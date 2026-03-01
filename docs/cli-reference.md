# CLI Reference

Complete reference for all cisshgo command-line flags and environment variables.

## Quick Reference

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--listeners` | `-l` | `50` | Number of SSH listeners to spawn |
| `--starting-port` | `-p` | `10000` | Starting port number |
| `--transcript-map` | `-t` | `transcripts/transcript_map.yaml` | Path to transcript map file |
| `--platform` | `-P` | `csr1000v` | Platform to emulate |
| `--inventory` | `-i` | - | Path to inventory file (optional) |

## Usage

```bash
cisshgo [flags]
```

## Flags

### -l, --listeners

**Type**: `int`  
**Default**: `50`  
**Environment**: `CISSHGO_LISTENERS`

Number of SSH listeners to spawn.

```bash
# Spawn 10 listeners
./cisshgo --listeners 10
./cisshgo -l 10

# Using environment variable
export CISSHGO_LISTENERS=10
./cisshgo
```

Ignored when using `--inventory`.

### -p, --starting-port

**Type**: `int`  
**Default**: `10000`  
**Environment**: `CISSHGO_STARTING_PORT`

Starting port number for SSH listeners. Listeners are spawned on sequential ports.

```bash
# Start at port 20000
./cisshgo --starting-port 20000
./cisshgo -p 20000

# Using environment variable
export CISSHGO_STARTING_PORT=20000
./cisshgo
```

With default settings (50 listeners, starting port 10000), listeners run on ports 10000-10049.

### -t, --transcript-map

**Type**: `path`  
**Default**: `transcripts/transcript_map.yaml`  
**Environment**: `CISSHGO_TRANSCRIPT_MAP`

Path to the transcript map YAML file.

```bash
# Use custom transcript map
./cisshgo --transcript-map /path/to/custom_map.yaml
./cisshgo -t /path/to/custom_map.yaml

# Using environment variable
export CISSHGO_TRANSCRIPT_MAP=/path/to/custom_map.yaml
./cisshgo
```

See [Configuration](configuration.md#transcript-map) for transcript map format.

### -P, --platform

**Type**: `string`  
**Default**: `csr1000v`  
**Environment**: `CISSHGO_PLATFORM`

Platform to use when no inventory is provided. Must match a platform key in the transcript map.

```bash
# Use IOS platform
./cisshgo --platform ios
./cisshgo -P ios

# Using environment variable
export CISSHGO_PLATFORM=ios
./cisshgo
```

Ignored when using `--inventory`.

### -i, --inventory

**Type**: `path`  
**Optional**: Yes  
**Environment**: `CISSHGO_INVENTORY`

Path to inventory YAML file for multi-device topologies.

```bash
# Use inventory file
./cisshgo --inventory transcripts/inventory_example.yaml
./cisshgo -i transcripts/inventory_example.yaml

# Using environment variable
export CISSHGO_INVENTORY=transcripts/inventory_example.yaml
./cisshgo
```

When specified:
- `--listeners` and `--platform` flags are ignored
- Devices are spawned according to inventory configuration
- Each device gets a sequential port starting from `--starting-port`

See [Configuration](configuration.md#inventory) for inventory format.

## Examples

### Single Device

```bash
# One CSR1000v on port 10000
./cisshgo --listeners 1 --starting-port 10000 --platform csr1000v
```

### Multiple Platforms

```bash
# Use inventory for mixed topology
./cisshgo --inventory my_topology.yaml --starting-port 20000
```

### Custom Transcripts

```bash
# Use custom transcript map
./cisshgo --transcript-map /opt/transcripts/custom.yaml --platform ios
```

### Environment Variables

```bash
# Configure entirely via environment
export CISSHGO_LISTENERS=5
export CISSHGO_STARTING_PORT=30000
export CISSHGO_PLATFORM=nxos
export CISSHGO_TRANSCRIPT_MAP=/opt/transcripts/map.yaml
./cisshgo
```

### Docker

```bash
# Single listener
docker run -d -p 10000:10000 \
  ghcr.io/tbotnz/cisshgo:latest \
  --listeners 1 --starting-port 10000

# Custom platform
docker run -d -p 10000-10009:10000-10009 \
  ghcr.io/tbotnz/cisshgo:latest \
  --listeners 10 --platform ios

# With custom transcripts (mount volume)
docker run -d -p 10000-10049:10000-10049 \
  -v /path/to/transcripts:/transcripts \
  ghcr.io/tbotnz/cisshgo:latest \
  --transcript-map /transcripts/custom.yaml
```

## Exit Codes

- `0` - Success
- `1` - Configuration error (invalid flags, missing files, validation failure)
- `2` - Runtime error (port binding failure, etc.)

## Logging

cisshgo uses structured logging (slog) and logs to stderr:

```text
2026/02/28 19:00:00 INFO Starting cisshgo listeners=50 startingPort=10000 platform=csr1000v
2026/02/28 19:00:00 INFO Listener started port=10000 platform=csr1000v
2026/02/28 19:00:00 INFO Listener started port=10001 platform=csr1000v
...
```

## Help

```bash
# Show help
./cisshgo -h
./cisshgo --help
```
