# Migration Guide: v0.2.0 to v1.0.0

This guide covers all breaking changes and new features introduced in v1.0.0.

---

## Breaking Changes

### 1. CLI Flag Names

All flags now use kebab-case (double-dash) instead of camelCase (single-dash).

| v0.2.0 | v1.0.0 | Short flag |
|--------|--------|------------|
| `-startingPort` | `--starting-port` | `-p` |
| `-transcriptMap` | `--transcript-map` | `-t` |
| `-listeners` | `--listeners` | `-l` |

**Before:**
```bash
./cisshgo -startingPort 10000 -transcriptMap custom.yaml -listeners 10
```

**After:**
```bash
./cisshgo --starting-port 10000 --transcript-map custom.yaml --listeners 10
```

All flags can also be set via environment variables:

```bash
CISSHGO_STARTING_PORT=10000 CISSHGO_TRANSCRIPT_MAP=custom.yaml ./cisshgo
```

---

### 2. Transcript Map YAML Schema

The `platforms` section changed from a list of single-key maps to a plain map.

**Before (v0.2.0):**
```yaml
platforms:
  - csr1000v:
      vendor: "cisco"
      hostname: "router"
      password: "admin"
      command_transcripts:
        "show version": "transcripts/cisco/csr1000v/show_version.txt"
      context_search:
        "enable": "#"
        "base": ">"
      context_hierarchy:
        "#": ">"
        ">": "exit"
```

**After (v1.0.0):**
```yaml
platforms:
  csr1000v:
    vendor: "cisco"
    hostname: "router"
    password: "admin"
    command_transcripts:
      "show version": "transcripts/cisco/csr1000v/show_version.txt"
    context_search:
      "enable": "#"
      "base": ">"
    context_hierarchy:
      "#": ">"
      ">": "exit"
```

See the [migration script](#migration-script) to automate this conversion.

---

### 3. Transcript Paths Resolved Relative to Transcript Map File

In v0.2.0, transcript paths were resolved relative to the **working directory** of the process. In v1.0.0, they are resolved relative to the **directory containing the transcript map file**.

If you run cisshgo from outside the repo root, update your transcript paths accordingly, or use absolute paths.

---

## New Features

### Inventory System (`--inventory`)

Define multi-device topologies in a YAML file instead of spawning identical listeners:

```yaml
# inventory.yaml
devices:
  - platform: csr1000v
    count: 10
  - platform: nxos
    count: 5
```

```bash
./cisshgo --inventory inventory.yaml --transcript-map transcript_map.yaml
```

### `--platform` Flag

When not using an inventory file, specify which platform to use (default: `csr1000v`):

```bash
./cisshgo --platform nxos --listeners 5
```

### Scenario-Based Stateful Responses

Define ordered command/response sequences for testing stateful workflows:

```yaml
scenarios:
  add-interface:
    platform: csr1000v
    sequence:
      - command: "show running-config"
        transcript: "transcripts/before.txt"
      - command: "interface GigabitEthernet0/0/2"
        transcript: "transcripts/generic_empty_return.txt"
      - command: "show running-config"
        transcript: "transcripts/after.txt"
```

Reference scenarios in your inventory:

```yaml
devices:
  - scenario: add-interface
    count: 3
```

### Username Enforcement

Optionally enforce a specific SSH username per platform:

```yaml
platforms:
  junos:
    username: "admin"   # only "admin" can connect
    password: "admin"
```

### Flexible Prompt Formatting

Use `prompt_format` for non-Cisco prompt styles:

```yaml
platforms:
  junos:
    username: "admin"
    prompt_format: "{username}@{hostname}{context}"
    # renders as: admin@hostname>
```

### Multi-line Prompts (`context_prefix_lines`)

For platforms like Junos that show a line above the prompt in config mode:

```yaml
platforms:
  junos:
    context_prefix_lines:
      "#": "[edit]"
    # config mode renders as:
    # [edit]
    # admin@hostname#
```

### Additional Platform Transcripts

v1.0.0 ships with transcripts for 7 platforms out of the box:

| Platform | Vendor |
|----------|--------|
| `csr1000v` | Cisco IOS XE |
| `ios` | Cisco IOS |
| `iosxr` | Cisco IOS XR |
| `asa` | Cisco ASA |
| `nxos` | Cisco NX-OS |
| `eos` | Arista EOS |
| `junos` | Juniper Junos |

### Graceful Shutdown

cisshgo now handles `SIGINT`/`SIGTERM` and shuts down all listeners cleanly.

### `--version` Flag

```bash
./cisshgo --version
```

---

## Migration Script

A Python script is provided to automate the transcript map schema migration:

```bash
# Requires Python 3 and PyYAML (pip install pyyaml)

# Preview output:
python3 scripts/migrate_transcript_map.py transcript_map.yaml

# Migrate in-place:
python3 scripts/migrate_transcript_map.py transcript_map.yaml > /tmp/map.yaml && mv /tmp/map.yaml transcript_map.yaml
```

The script handles:
- v0.2.0 list-of-maps → v1.0.0 map (migrates)
- Already v1.0.0 format (no-op, reports "no changes needed")
- Empty platforms list → empty map
- Missing `platforms` key → no change

---

## Step-by-Step Upgrade Checklist

1. **Update CLI flags** in any scripts or tooling (see [CLI flag table](#1-cli-flag-names))
2. **Migrate transcript map** from list-of-maps to map format (see [migration script](migration-script.md) or [schema change](#2-transcript-map-yaml-schema))
3. **Verify transcript paths** resolve correctly from the transcript map file's directory
4. **Optionally adopt** inventory files, scenarios, or new platform transcripts

---

## Troubleshooting

**`unknown flag: -startingPort`**
→ Rename to `--starting-port` (see [CLI flag names](#1-cli-flag-names))

**`parsing transcript map: yaml: unmarshal errors`**
→ Your transcript map uses the old list-of-maps format. Run the [migration script](migration-script.md) or update manually.

**`transcript map validation failed: platform "X" command "Y": <path>`**
→ Transcript paths are now resolved relative to the transcript map file. Check that paths are correct relative to where your `transcript_map.yaml` lives.

**`platform "csr1000v" not found in transcript map`**
→ The `--platform` flag defaults to `csr1000v`. If your transcript map uses a different platform name, pass `--platform <name>`.
