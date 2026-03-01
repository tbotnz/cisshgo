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

The top-level key (e.g., `csr1000v`) is the platform identifier used with `-platform` flag or in inventory files.

#### vendor

The device vendor (e.g., `cisco`, `arista`, `juniper`).

#### hostname

The hostname displayed in the CLI prompt. Supports Go template variables (see [Transcripts](transcripts.md)).

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

Scenarios enable stateful command responses. A scenario is a sequence of transcript files that change based on previous commands.

```yaml
platforms:
  csr1000v-add-interface:
    vendor: "cisco"
    hostname: "cisshgo1000v"
    password: "admin"
    scenario:
      - trigger: "interface GigabitEthernet3"
        transcript_updates:
          "show running-config": "transcripts/scenarios/csr1000v-add-interface/running_config_after.txt"
    command_transcripts:
      "show running-config": "transcripts/scenarios/csr1000v-add-interface/running_config_before.txt"
    context_hierarchy:
      "(config)#": "#"
      "#": ">"
      ">": "exit"
    context_search:
      "configure terminal": "(config)#"
      "enable": "#"
      "base": ">"
```

When the trigger command is executed, the specified transcripts are updated for subsequent commands.

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

A scenario identifier (also defined in the transcript map). Mutually exclusive with `platform`.

#### count

Number of listeners to spawn for this platform/scenario.

### Usage

```bash
./cisshgo -inventory transcripts/inventory_example.yaml
```

When using an inventory file:
- The `-listeners` flag is ignored
- Listeners are spawned sequentially starting from `-startingPort`
- Each device gets a unique port

Example with 6 total devices starting at port 10000:
- Ports 10000-10001: csr1000v (2 devices)
- Ports 10002-10004: ios (3 devices)
- Port 10005: csr1000v-add-interface scenario (1 device)

## Path Resolution

All paths in the transcript map are resolved relative to the transcript map file location, not the current working directory. This allows the transcript map to be placed anywhere without breaking relative paths.

## Validation

cisshgo validates the transcript map at startup:

- All referenced transcript files must exist
- Platform/scenario references in inventory must exist in transcript map
- Inventory entries must specify exactly one of `platform` or `scenario`

If validation fails, cisshgo exits with an error before spawning any listeners.
