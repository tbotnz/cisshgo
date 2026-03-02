# Transcripts

Transcripts are plain text files containing command output. They support Go templates for dynamic content.

## Creating Transcripts

### Basic Transcript

Create a plain text file with the command output:

```text
Cisco IOS XE Software, Version 16.04.01
Cisco IOS Software [Everest], CSR1000V Software (X86_64_LINUX_IOSD-UNIVERSALK9-M), Version 16.4.1, RELEASE SOFTWARE (fc2)
Technical Support: http://www.cisco.com/techsupport
Copyright (c) 1986-2016 by Cisco Systems, Inc.
Compiled Sun 27-Nov-16 13:02 by mcpre
```

Save it in the appropriate vendor/platform directory:

```
transcripts/cisco/csr1000v/show_version.txt
```

### Adding to Transcript Map

Add an entry in `transcripts/transcript_map.yaml`:

```yaml
platforms:
  csr1000v:
    command_transcripts:
      "show version": "transcripts/cisco/csr1000v/show_version.txt"
```

## Go Templates

Transcripts support Go template syntax for dynamic content.

### Available Variables

The `FakeDevice` struct provides these template variables:

```go
type FakeDevice struct {
    Vendor            string            // Device vendor (e.g., "cisco")
    Platform          string            // Platform identifier (e.g., "csr1000v")
    Hostname          string            // Device hostname
    Username          string            // SSH username (from platform config)
    Password          string            // SSH password
    SupportedCommands SupportedCommands // Available commands
    ContextSearch     map[string]string // CLI contexts
    ContextHierarchy  map[string]string // Context navigation
}
```

### Template Example

```text
ROM: IOS-XE ROMMON
{{.Hostname}} uptime is 4 hours, 55 minutes
Uptime for this control processor is 4 hours, 56 minutes
System returned to ROM by reload
System image file is "bootflash:packages.conf"
Last reload reason: reload
```

When rendered with hostname `cisshgo1000v`:

```text
ROM: IOS-XE ROMMON
cisshgo1000v uptime is 4 hours, 55 minutes
Uptime for this control processor is 4 hours, 56 minutes
System returned to ROM by reload
System image file is "bootflash:packages.conf"
Last reload reason: reload
```

#### Using Templates in Practice

1. Create a transcript file `transcripts/cisco/ios/show_version.txt`:

```text
Cisco IOS Software, {{.Platform}} Software
ROM: Bootstrap program is IOS
{{.Hostname}} uptime is 1 day, 2 hours, 30 minutes
```

2. Reference it in the transcript map:

```yaml
platforms:
  ios:
    hostname: "my-router"
    command_transcripts:
      "show version": "transcripts/cisco/ios/show_version.txt"
```

3. Run and test:

```bash
./cisshgo --platform ios --listeners 1
ssh -p 10000 admin@localhost
```

Output will show:
```text
Cisco IOS Software, ios Software
ROM: Bootstrap program is IOS
my-router uptime is 1 day, 2 hours, 30 minutes
```

### Template Functions

Standard Go template functions are available (see [text/template documentation](https://pkg.go.dev/text/template) for the full reference):

- String manipulation: `printf`, `print`, `println`
- Conditionals: `if`, `else`, `end`
- Loops: `range`, `end`
- Comparisons: `eq`, `ne`, `lt`, `le`, `gt`, `ge`

Example with conditionals:

```text
{{if eq .Platform "csr1000v"}}
Cisco CSR1000V (VXE) processor with 2070375K/3075K bytes of memory.
{{else}}
Cisco IOS processor with 1048576K bytes of memory.
{{end}}
```

## Command Matching

Commands are matched using fuzzy matching, which allows for:

- **Partial commands**: `sh ver` matches `show version`
- **Abbreviations**: `conf t` matches `configure terminal`
- **Extra whitespace**: `show  version` matches `show version`

### Matching Algorithm

The fuzzy matcher:

1. Normalizes whitespace
2. Splits command into tokens
3. Matches each token as a prefix
4. Returns the best match if unambiguous

### Ambiguous Commands

If multiple commands match, cisshgo returns an error:

```text
test_device#sh
% Ambiguous command: "sh"
```

To avoid ambiguity, ensure command prefixes are unique within each context.

## Supported Platforms

cisshgo includes transcripts for these platforms:

### Cisco

- **IOS** (`ios`) - Classic Cisco IOS
- **IOS-XE** (`csr1000v`) - Cisco CSR1000v
- **IOS-XR** (`iosxr`) - Cisco IOS-XR
- **NX-OS** (`nxos`) - Cisco Nexus
- **ASA** (`asa`) - Cisco ASA Firewall

### Arista

- **EOS** (`eos`) - Arista EOS

### Juniper

- **JunOS** (`junos`) - Juniper JunOS

## Adding New Platforms

### Cisco-style Platforms

For devices with Cisco-like CLI (enable mode, config mode, etc.):

1. Create transcript directory: `transcripts/vendor/platform/`
2. Add transcript files for common commands
3. Add platform entry to `transcript_map.yaml`

Example:

```yaml
platforms:
  new_platform:
    vendor: "cisco"
    hostname: "new-device"
    password: "admin"
    command_transcripts:
      "show version": "transcripts/cisco/new_platform/show_version.txt"
    context_hierarchy:
      "(config)#": "#"
      "#": ">"
      ">": "exit"
    context_search:
      "configure terminal": "(config)#"
      "enable": "#"
      "base": ">"
```

### Non-Cisco Platforms

For devices with different CLI patterns (e.g., Juniper, F5), you'll need to implement a custom handler in the `ssh_server/handlers` package. This is an advanced topic - see the [Contributing](contributing.md) guide for details.

## Common Commands

Most platforms should support these basic commands:

- `show version` - Device version information
- `show running-config` - Current configuration
- `show ip interface brief` - Interface status (Cisco)
- `show interfaces` - Interface details (Juniper)
- `terminal length 0` - Disable pagination
- `enable` - Enter privileged mode
- `configure terminal` - Enter configuration mode
- `exit` - Exit current mode

## Empty Returns

For commands that produce no output (like `terminal length 0`), use the generic empty return:

```yaml
command_transcripts:
  "terminal length 0": "transcripts/generic_empty_return.txt"
```

The file should be empty or contain only whitespace.

## Troubleshooting

### Template Errors

If you see errors like `template: ...: executing ... at <.InvalidField>: can't evaluate field`:

- Check that the field name matches exactly (case-sensitive)
- Available fields: `Vendor`, `Platform`, `Hostname`, `Username`, `Password`
- Use `{{.Hostname}}` not `{{.hostname}}`

### Command Not Found

If cisshgo returns "Command not found" or no output:

- Verify the command is spelled correctly in `transcript_map.yaml`
- Check that the command is defined for the current CLI context
- Try the full command - fuzzy matching may not work for very short inputs (e.g., `s` alone)
- Use `show version` instead of just `sh` if ambiguous

### Ambiguous Command Errors

If you see `% Ambiguous command`:

- Multiple commands match your input
- Use more characters to make the command unique
- Example: `sh` matches both `show` and `shutdown` - use `sho` or `show`

### File Not Found Errors

If cisshgo fails to start with "file not found" or "no such file":

- Check that all transcript paths in `transcript_map.yaml` are correct
- Paths are relative to the transcript map file location, not your current directory
- Verify files exist: `ls -la transcripts/cisco/ios/show_version.txt`
- Check for typos in filenames

### Connection Refused

If SSH connection fails:

- Verify cisshgo is running: check for "Listener started" log messages
- Check the port number matches: `ssh -p 10000` for default starting port
- Ensure no firewall is blocking the port
- Try `netstat -tuln | grep 10000` to verify the port is listening
