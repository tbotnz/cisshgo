# Contributing to cisshgo

Thanks for your interest in contributing. This document covers how to get set up, the development workflow, and how to add new device transcripts.

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

1. Fork the repo and create a branch from `master`
2. Make your changes
3. Ensure tests pass and coverage stays above 90%: `go test -race -coverprofile=coverage.out ./...`
4. Ensure code is formatted: `gofmt -w ./...`
5. Open a pull request against `master`

Branch naming convention:
- `feat/short-description` — new features
- `fix/short-description` — bug fixes
- `chore/short-description` — maintenance, refactoring, docs

## Running Tests

```bash
# All tests with race detection
go test -race ./...

# With coverage report
go test -coverprofile=coverage.out -covermode=atomic ./...
go tool cover -func=coverage.out
```

Coverage must stay at or above 90%. The CI workflow enforces this on every PR.

## Adding Command Transcripts

The simplest contribution is adding output for a new command on an existing platform.

1. Create a plain text file in the appropriate `transcripts/vendor/platform/` directory
2. Add an entry to `transcripts/transcript_map.yaml`:

```yaml
platforms:
  - csr1000v:
      command_transcripts:
        "show new command": "transcripts/cisco/csr1000v/show_new_command.txt"
```

Transcripts support Go templates. The following variables from `FakeDevice` are available:

| Variable | Description |
|----------|-------------|
| `{{.Hostname}}` | Current hostname of the device |
| `{{.Vendor}}` | Vendor string (e.g. `cisco`) |
| `{{.Platform}}` | Platform string (e.g. `csr1000v`) |
| `{{.Password}}` | Device password |

## Adding a New Platform

1. Create a directory under `transcripts/vendor/platform/`
2. Add transcript files for each supported command
3. Add the platform entry to `transcript_map.yaml` with `hostname`, `password`, `command_transcripts`, `context_search`, and `context_hierarchy`
4. Add a `FakeDevice` initialization test in `fakedevices/genericFakeDevice_test.go`

See the existing `csr1000v` entry as a reference.

## Commit Messages

Follow [Conventional Commits](https://www.conventionalcommits.org/):

```
feat: add show interfaces command for csr1000v
fix: handle empty input in exec mode
chore: update dependencies
```

Conventional commits are required — they feed directly into the automated `CHANGELOG.md` generation via [git-cliff](https://git-cliff.org/) on each release. Non-conventional commits will be silently excluded from the changelog.

## Reporting Issues

Open a GitHub issue. For bugs, include:
- cisshgo version (the release tag, e.g. `v0.2.0`)
- The SSH client and command you were running
- Expected vs actual behavior
