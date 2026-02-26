# cisshgo

[![CI](https://github.com/tbotnz/cisshgo/actions/workflows/test.yml/badge.svg)](https://github.com/tbotnz/cisshgo/actions/workflows/test.yml)
[![coverage](https://raw.githubusercontent.com/tbotnz/cisshgo/badges/.badges/main/coverage.svg)](https://github.com/tbotnz/cisshgo/actions/workflows/test.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/tbotnz/cisshgo)](https://goreportcard.com/report/github.com/tbotnz/cisshgo)
[![Go Reference](https://pkg.go.dev/badge/github.com/tbotnz/cisshgo.svg)](https://pkg.go.dev/github.com/tbotnz/cisshgo)
[![Release](https://img.shields.io/github/v/release/tbotnz/cisshgo)](https://github.com/tbotnz/cisshgo/releases)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Simple, small, fast, concurrent SSH server to emulate network equipment (i.e. Cisco IOS) for testing purposes.

## Quick Start

### Using Pre-built Binaries

Download the latest release from [GitHub Releases](https://github.com/tbotnz/cisshgo/releases) and run:

```bash
./cisshgo
```

### Building from Source

#### Quick Run

```bash
go run cissh.go
```

#### Build and Run

```bash
go build -o cisshgo cissh.go
./cisshgo
```

### Using Docker

Pull and run the latest release:

```bash
docker run -d -p 10000-10049:10000-10049 ghcr.io/tbotnz/cisshgo:latest
```

Or run with custom options:

```bash
docker run -d -p 10000:10000 ghcr.io/tbotnz/cisshgo:latest -listeners 1 -startingPort 10000
```

Or build locally:

```bash
docker build -t cisshgo .
docker run -d -p 10000-10049:10000-10049 cisshgo
```

### Using GoReleaser (for maintainers)

Build a local snapshot release:

```bash
goreleaser release --snapshot --clean --skip=publish
```

## Releasing

Releases are automated via GitHub Actions. To create a new release:

1. Create and push a tag:

   ```bash
   git tag v0.1.2
   git push origin v0.1.2
   ```

2. GitHub Actions will automatically:
   - Build binaries for all platforms (linux/darwin/windows, amd64/arm64)
   - Create multi-arch Docker images and push to Docker Hub
   - Generate SBOMs for security compliance
   - Create GitHub Release with binaries, archives, and checksums
   - Build deb/rpm packages

### Required Secrets

The following secrets must be configured in the GitHub repository:

- `DOCKER_USERNAME` - Docker Hub username
- `DOCKER_PASSWORD` - Docker Hub token/password

## Connecting

SSH into any of the open ports with `admin` as the password:

```bash
ssh -p 10000 admin@localhost
```

Default password: `admin`

## Example Session

```text
test_device#show version
Cisco IOS XE Software, Version 16.04.01
Cisco IOS Software [Everest], CSR1000V Software (X86_64_LINUX_IOSD-UNIVERSALK9-M), Version 16.4.1, RELEASE SOFTWARE (fc2)
Technical Support: http://www.cisco.com/techsupport
Copyright (c) 1986-2016 by Cisco Systems, Inc.
Compiled Sun 27-Nov-16 13:02 by mcpre
...
ROM: IOS-XE ROMMON
```

Available commands:

- `show version`
- `show ip interface brief`
- `show running-config`

Additional commands can be added by modifying `transcripts/transcript_map.yaml`.

## Advanced Usage

### Command Line Options

```text
  -listeners int
        How many listeners do you wish to spawn? (default 50)
  -startingPort int
        What port do you want to start at? (default 10000)
  -transcriptMap string
        What file contains the map of commands to transcribed output? (default "transcripts/transcript_map.yaml")
```

### Example: Single Listener

```bash
./cisshgo -listeners 1 -startingPort 10000
```

## Expanding Platform Support

cisshgo is built with modularity in mind to support easy expansion or customization.

### Customized Output in Command Transcripts

Transcripts support Go templating. For example, in `show_version.txt`:

```text
ROM: IOS-XE ROMMON
{{.Hostname}} uptime is 4 hours, 55 minutes
Uptime for this control processor is 4 hours, 56 minutes
```

Available template variables from `fakedevices.FakeDevice`:

```go
type FakeDevice struct {
    Vendor            string            // Vendor of this fake device
    Platform          string            // Platform of this fake device
    Hostname          string            // Hostname of the fake device
    Password          string            // Password of the fake device
    SupportedCommands SupportedCommands // What commands this fake device supports
    ContextSearch     map[string]string // The available CLI prompt/contexts on this fake device
    ContextHierarchy  map[string]string // The hierarchy of the available contexts
}
```

### Adding Additional Command Transcripts

1. Create a plain text file in the appropriate `vendor/platform` folder
2. Add an entry in `transcripts/transcript_map.yaml`:

```yaml
---
platforms:
  - csr1000v:
      command_transcripts:
        "my new fancy command": "transcripts/cisco/csr1000v/my_new_fancy_command.txt"
```

### Adding Additional "Cisco-style" Platforms

Supply additional device types and transcripts in `transcript_map.yaml`.
This works for devices with similar interaction patterns (e.g., `configure terminal` leading to `(config)#` mode).

### Adding Additional Non-"Cisco-style" Platforms

**NOTE:** This feature is not fully implemented yet!

For platforms with different interaction patterns (e.g., Juniper, F5):

1. Implement a new handler module under `ssh_server/handlers`
2. Add it to the device mapping in `cissh.go`

The handler controls SSH session emulation and provides conditional logic to simulate the device experience.

## License

MIT License - see [LICENSE](LICENSE) file for details.

## Disclaimer

Cisco IOS is the property/trademark of Cisco.
