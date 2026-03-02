# Contributing

Thanks for your interest in contributing to cisshgo! This guide covers development setup, workflow, and how to add new features.

## Development Setup

**Prerequisites**: Go 1.26+

```bash
git clone https://github.com/tbotnz/cisshgo
cd cisshgo
go build ./...
go test ./...
```

## Workflow

cisshgo uses [GitHub Flow](https://docs.github.com/en/get-started/using-github/github-flow):

1. Fork the repo and create a branch from `main`
2. Make your changes
3. Ensure tests pass and coverage stays above 90%
4. Ensure code is formatted with `gofmt`
5. Open a pull request against `main`

### Branch Naming

- `feat/short-description` — new features
- `fix/short-description` — bug fixes
- `docs/short-description` — documentation
- `refactor/short-description` — refactoring
- `chore/short-description` — maintenance

### Commit Messages

Follow [Conventional Commits](https://www.conventionalcommits.org/):

```
feat: add show interfaces command for csr1000v
fix: handle empty input in exec mode
docs: update CLI reference
chore: update dependencies
```

Conventional commits are **required** — they feed into automated `CHANGELOG.md` generation via [git-cliff](https://git-cliff.org/). Non-conventional commits are excluded from the changelog.

## Running Tests

```bash
# All tests with race detection
go test -race ./...

# With coverage report
go test -coverprofile=coverage.out -covermode=atomic ./...
go tool cover -func=coverage.out

# View coverage in browser
go tool cover -html=coverage.out
```

**Coverage requirement**: 90% minimum. The CI workflow enforces this on every PR.

## Code Quality

```bash
# Format code
gofmt -w ./...

# Run linter (if golangci-lint is installed)
golangci-lint run
```

## Adding Command Transcripts

The simplest contribution is adding output for a new command on an existing platform.

### Steps

1. Create a plain text file in `transcripts/vendor/platform/`:

```bash
# Example: add "show interfaces" for IOS
vim transcripts/cisco/ios/show_interfaces.txt
```

2. Add the command to `transcripts/transcript_map.yaml`:

```yaml
platforms:
  ios:
    command_transcripts:
      "show interfaces": "transcripts/cisco/ios/show_interfaces.txt"
```

3. Test it:

```bash
./cisshgo --platform ios --listeners 1
ssh -p 10000 admin@localhost
# Try: show interfaces
```

### Using Templates

Transcripts support Go templates for dynamic content:

```text
{{.Hostname}} uptime is 4 hours, 55 minutes
Processor board ID {{.Platform}}-12345
```

Available variables: `Hostname`, `Vendor`, `Platform`, `Password`

See [Transcripts](transcripts.md#go-templates) for details.

## Adding a New Platform

### For Cisco-style Devices

1. Create transcript directory:

```bash
mkdir -p transcripts/vendor/platform
```

2. Add transcript files for common commands:
   - `show_version.txt`
   - `show_running-config.txt`
   - `show_ip_interface_brief.txt` (or equivalent)

3. Add platform entry to `transcript_map.yaml`:

```yaml
platforms:
  new_platform:
    vendor: "vendor_name"
    hostname: "device-hostname"
    password: "admin"
    command_transcripts:
      "show version": "transcripts/vendor/platform/show_version.txt"
      "terminal length 0": "transcripts/generic_empty_return.txt"
    context_hierarchy:
      "(config)#": "#"
      "#": ">"
      ">": "exit"
    context_search:
      "configure terminal": "(config)#"
      "enable": "#"
      "base": ">"
```

4. Add test in `fakedevices/genericFakeDevice_test.go`:

```go
func TestNewPlatformInitialization(t *testing.T) {
    // Test platform loads correctly
}
```

5. Test the platform:

```bash
./cisshgo --platform new_platform --listeners 1
ssh -p 10000 admin@localhost
```

### For Non-Cisco Devices

Devices with different CLI patterns (Juniper, F5, etc.) require a custom handler in `ssh_server/handlers/`. This is an advanced topic — open an issue to discuss the approach before implementing.

## Adding Tests

When adding new features, include tests:

```go
func TestNewFeature(t *testing.T) {
    // Arrange
    input := "test input"
    
    // Act
    result := NewFeature(input)
    
    // Assert
    if result != expected {
        t.Errorf("got %v, want %v", result, expected)
    }
}
```

Run tests frequently during development:

```bash
go test -v ./...
```

## Documentation

Update documentation when adding features:

- **README.md** - Quick start and overview
- **docs/** - Detailed documentation (this site)
- **Code comments** - Exported functions and types
- **CHANGELOG.md** - Automatically generated from commits

## Reporting Issues

Open a [GitHub issue](https://github.com/tbotnz/cisshgo/issues). For bugs, include:

- cisshgo version (e.g., `v0.2.0`)
- SSH client and command used
- Expected vs actual behavior
- Relevant logs

## Getting Help

- **Issues**: [github.com/tbotnz/cisshgo/issues](https://github.com/tbotnz/cisshgo/issues)
- **Discussions**: [github.com/tbotnz/cisshgo/discussions](https://github.com/tbotnz/cisshgo/discussions)

## Code of Conduct

Be respectful and constructive. We're all here to build something useful.
