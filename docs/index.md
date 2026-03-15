# cisshgo

Simple, small, fast, concurrent SSH server to emulate network equipment (i.e. Cisco IOS) for testing purposes.

## What is cisshgo?

cisshgo is a lightweight SSH server that emulates network devices like Cisco routers and switches. It's designed for testing network automation tools, scripts, and applications without requiring physical hardware or virtual machines.

## Key Features

- **Fast & Lightweight**: Written in Go, minimal resource usage
- **Concurrent**: Spawn multiple device instances simultaneously
- **Flexible**: Support for multiple vendors and platforms (Cisco IOS/IOS-XR/NX-OS/ASA, Arista EOS, Juniper JunOS)
- **Customizable**: Easy to add new commands and platforms via YAML configuration
- **Stateful**: Scenario support for commands that change device state
- **Multi-device**: Inventory system for complex network topologies
- **Template Support**: Go templates in command outputs for dynamic responses

## Use Cases

- Testing network automation scripts (Ansible, Netmiko, NAPALM)
- CI/CD pipelines for network automation
- Development and debugging without physical equipment
- Training and demonstrations
- Integration testing for network management tools

## Quick Example

```bash
# Start 50 SSH listeners on ports 10000-10049
./cisshgo

# Connect to any device
ssh -p 10000 admin@localhost
```

```text
test_device#show version
Cisco IOS XE Software, Version 16.04.01
...
```

## Project Status

[![CI](https://github.com/tbotnz/cisshgo/actions/workflows/test.yml/badge.svg)](https://github.com/tbotnz/cisshgo/actions/workflows/test.yml)
[![coverage](https://raw.githubusercontent.com/tbotnz/cisshgo/badges/.badges/master/coverage.svg)](https://github.com/tbotnz/cisshgo/actions/workflows/test.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/tbotnz/cisshgo)](https://goreportcard.com/report/github.com/tbotnz/cisshgo)

## License

MIT License - see [LICENSE](https://github.com/tbotnz/cisshgo/blob/master/LICENSE) file for details.

## Disclaimer

Cisco IOS is the property/trademark of Cisco Systems, Inc.
