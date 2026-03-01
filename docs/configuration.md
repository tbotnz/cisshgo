# Configuration

cisshgo uses YAML files for configuration. There are two types of configuration files:

1. **Transcript Map** - Defines platforms, commands, and CLI contexts
2. **Inventory** - Defines multi-device topologies

## Transcript Map

The transcript map (`transcripts/transcript_map.yaml`) defines available platforms and their command responses.

### Basic Structure

```yaml
platforms:
  csr1000v:
    vendor: "cisco"
    hostname: "cisshgo1000v"
    password: "admin"
    command_transcripts:
      "show version": "transcripts/cisco/csr1000v/show_version.txt"
      "show ip interface brief": "transcripts/cisco/csr1000v/show_ip_interface_brief.txt"
    context_hierarchy:
      "(config)#": "#"
      "#": ">"
      ">": "exit"
    context_search:
      "configure terminal": "(config)#"
      "enable": "#"
      "base": ">"
```

### Fields

#### Platform Key

The top-level key (e.g., `csr1000v`) is the platform identifier used with `--platform` flag or in inventory files.

#### vendor

The device vendor (e.g., `cisco`, `arista`, `juniper`).

#### hostname

The hostname displayed in the CLI prompt. Supports Go template variables (see [Transcripts](transcripts.md)).

#### username

Optional. The SSH username required to authenticate. When set, cisshgo enforces both username and password — connections with a different username are rejected. When omitted, any username is accepted (only the password is checked).

```yaml
junos:
  username: "admin"   # only "admin" can connect
  password: "admin"
```

Also available as `{{.Username}}` in transcript templates.

#### password

SSH password for authentication.

#### command_transcripts

Maps CLI commands to transcript files. Commands are matched using fuzzy matching (see [Transcripts](transcripts.md#command-matching)).

```yaml
command_transcripts:
  "show version": "transcripts/cisco/ios/show_version.txt"
  "show running-config": "transcripts/cisco/ios/show_running-config.txt"
  "terminal length 0": "transcripts/generic_empty_return.txt"
```

Paths are relative to the transcript map file location.

#### context_hierarchy

Defines the CLI context hierarchy and how to navigate between contexts.

```yaml
context_hierarchy:
  "(config)#": "#"      # From config mode, go to enable mode
  "#": ">"              # From enable mode, go to user mode
  ">": "exit"           # From user mode, exit
```

#### context_search

Maps commands to the context they enter.

```yaml
context_search:
  "configure terminal": "(config)#"
  "enable": "#"
  "base": ">"
```

The `base` key defines the initial context when a user connects.

### Scenarios

Scenarios enable stateful command responses that simulate device configuration changes. Unlike static platforms, scenarios define a sequence of commands that, when executed in order, change the device's responses to subsequent commands.

#### Structure

Scenarios are defined in a separate `scenarios:` section at the same level as `platforms:`:

```yaml
scenarios:
  csr1000v-add-interface:
    platform: csr1000v
    sequence:
      - command: "show running-config"
        transcript: "transcripts/scenarios/csr1000v-add-interface/running_config_before.txt"
      - command: "interface GigabitEthernet0/0/2"
        transcript: "transcripts/generic_empty_return.txt"
      - command: "ip address 172.16.0.1 255.255.255.0"
        transcript: "transcripts/generic_empty_return.txt"
      - command: "end"
        transcript: "transcripts/generic_empty_return.txt"
      - command: "show running-config"
        transcript: "transcripts/scenarios/csr1000v-add-interface/running_config_after.txt"
```

#### Fields

**platform**

The base platform identifier from the `platforms:` section. The scenario inherits all settings (vendor, hostname, password, contexts) from this platform.

**sequence**

An ordered list of command/transcript pairs. Each step defines:
- `command`: The exact command to match
- `transcript`: The transcript file to return for that command

#### How Scenarios Work

1. **Initial State**: When a scenario starts, the first command in the sequence is active
2. **Command Execution**: When a user executes a command that matches the current sequence step, that transcript is returned
3. **State Progression**: After a sequence command is executed, the scenario advances to the next step
4. **State Changes**: Subsequent commands (like the second `show running-config`) return different transcripts, simulating configuration changes
5. **Non-sequence Commands**: Commands not in the sequence fall back to the base platform's command_transcripts

#### Example Flow

Using the `csr1000v-add-interface` scenario above:

```text
1. User: show running-config
   → Returns: running_config_before.txt (no GigabitEthernet0/0/2)
   
2. User: interface GigabitEthernet0/0/2
   → Returns: empty (enters config mode)
   
3. User: ip address 172.16.0.1 255.255.255.0
   → Returns: empty (configures IP)
   
4. User: end
   → Returns: empty (exits config mode)
   
5. User: show running-config
   → Returns: running_config_after.txt (now includes GigabitEthernet0/0/2 with IP)
```

The device appears to have been configured, even though it's just playing back different transcripts.

#### Scenario Behavior

**Out-of-sequence commands**: Commands must be executed in the exact order defined in the sequence. Executing commands out of order or skipping steps will not advance the sequence state.

**Non-sequence commands**: Commands not in the sequence (like `show version`) work normally using the base platform's command transcripts.

**Sequence completion**: Once all sequence steps are executed, the scenario remains in its final state. The sequence does not reset or loop.

#### Running Scenarios

Scenarios are used via inventory files. Create an inventory file:

```yaml
devices:
  - scenario: csr1000v-add-interface
    count: 1
```

Then run cisshgo:

```bash
./cisshgo --inventory my_inventory.yaml --starting-port 10000
```

Connect and execute the sequence:

```bash
ssh -p 10000 admin@localhost
# Execute commands in order to progress through the scenario
```

## Inventory

The inventory file defines a multi-device topology with different platforms and counts.

### Structure

```yaml
devices:
  - platform: csr1000v
    count: 2
  - platform: ios
    count: 3
  - scenario: csr1000v-add-interface
    count: 1
```

### Fields

#### platform

The platform identifier from the transcript map.

#### scenario

A scenario identifier (also defined in the transcript map under `scenarios:`). Mutually exclusive with `platform`.

Scenarios are stateful command sequences - see [Scenarios](#scenarios) for details.

#### count

Number of listeners to spawn for this platform/scenario.

### Usage

```bash
./cisshgo --inventory transcripts/inventory_example.yaml
```

When using an inventory file:
- The `--listeners` flag is ignored
- Listeners are spawned sequentially starting from `--starting-port`
- Each device gets a unique port

Example with 6 total devices starting at port 10000:
- Ports 10000-10001: csr1000v (2 devices)
- Ports 10002-10004: ios (3 devices)
- Port 10005: csr1000v-add-interface scenario (1 device)

#### Complete Example

Create `my_lab.yaml`:

```yaml
devices:
  - platform: csr1000v
    count: 2
  - platform: ios
    count: 1
  - scenario: csr1000v-add-interface
    count: 1
```

Run cisshgo:

```bash
./cisshgo --inventory my_lab.yaml --starting-port 10000
```

Connect to devices:

```bash
# CSR1000v devices
ssh -p 10000 admin@localhost
ssh -p 10001 admin@localhost

# IOS device
ssh -p 10002 admin@localhost

# Scenario device (stateful)
ssh -p 10003 admin@localhost
```

## Path Resolution

All paths in the transcript map are resolved relative to the transcript map file location, not the current working directory. This allows the transcript map to be placed anywhere without breaking relative paths.

## Validation

cisshgo validates the transcript map at startup:

- All referenced transcript files must exist
- Platform/scenario references in inventory must exist in transcript map
- Inventory entries must specify exactly one of `platform` or `scenario`

If validation fails, cisshgo exits with an error before spawning any listeners.
