# Getting Started

## Installation

### Pre-built Binaries

Download the latest release from [GitHub Releases](https://github.com/tbotnz/cisshgo/releases):

```bash
# Download for your platform (linux/darwin/windows, amd64/arm64)
wget https://github.com/tbotnz/cisshgo/releases/latest/download/cisshgo_linux_amd64.tar.gz
tar -xzf cisshgo_linux_amd64.tar.gz
./cisshgo
```

### Docker

Pull and run the latest release:

```bash
docker run -d -p 10000-10049:10000-10049 ghcr.io/tbotnz/cisshgo:latest
```

Run with custom options:

```bash
docker run -d -p 10000:10000 ghcr.io/tbotnz/cisshgo:latest --listeners 1 --starting-port 10000
```

### Building from Source

**Requirements**: Go 1.26+

```bash
git clone https://github.com/tbotnz/cisshgo
cd cisshgo
go build -o cisshgo cissh.go
./cisshgo
```

Or run directly:

```bash
go run cissh.go
```

## Basic Usage

### Start the Server

Default configuration (50 listeners on ports 10000-10049):

```bash
./cisshgo
```

Single listener on port 10000:

```bash
./cisshgo --listeners 1 --starting-port 10000
```

### Connect to a Device

```bash
ssh -p 10000 admin@localhost
```

Default password: `admin`

> **Note**: The hostname shown in the prompt (e.g., `cisshgo1000v#`) is determined by the platform's configuration in the transcript map file. The default platform is `csr1000v` with hostname `cisshgo1000v`.

### Example Session

```text
$ ssh -p 10000 admin@localhost
admin@localhost's password: 
test_device>enable
test_device#show version
Cisco IOS XE Software, Version 16.04.01
Cisco IOS Software [Everest], CSR1000V Software (X86_64_LINUX_IOSD-UNIVERSALK9-M), Version 16.4.1, RELEASE SOFTWARE (fc2)
Technical Support: http://www.cisco.com/techsupport
Copyright (c) 1986-2016 by Cisco Systems, Inc.
Compiled Sun 27-Nov-16 13:02 by mcpre

test_device#show ip interface brief
Interface              IP-Address      OK? Method Status                Protocol
GigabitEthernet1       10.0.0.1        YES NVRAM  up                    up      
GigabitEthernet2       unassigned      YES NVRAM  administratively down down    

test_device#configure terminal
test_device(config)#exit
test_device#
```

## Environment Variables

All CLI flags can be set via environment variables:

```bash
export CISSHGO_LISTENERS=10
export CISSHGO_STARTING_PORT=20000
export CISSHGO_PLATFORM=ios
./cisshgo
```

See [CLI Reference](cli-reference.md) for all available options.

## Next Steps

- [Configuration](configuration.md) - Learn about transcript map file and inventory files
- [Transcripts](transcripts.md) - Customize command outputs and add new commands
- [CLI Reference](cli-reference.md) - Complete CLI flag documentation
