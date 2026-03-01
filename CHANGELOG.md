# Changelog

All notable changes to this project will be documented in this file.

This project adheres to [Keep a Changelog](https://keepachangelog.com/en/1.0.0/)
and [Semantic Versioning](https://semver.org/).

**Stability commitment (v1.0.0+):** No breaking changes in minor or patch releases.
Deprecated features will be supported for at least one major version.
Migration guides will be provided for all major version upgrades.

Changes are generated from [Conventional Commits](https://www.conventionalcommits.org/).
## [Unreleased]

### ⚠️ Breaking Changes

- Document breaking CLI flag renames from Kong migration ([4d7da81](https://github.com/tbotnz/cisshgo/commit/4d7da81880a6ca65f30f96e8c4c9b5fb0061e438))

CLI flags were renamed as part of the Kong migration in #64.
Users with scripts or tooling using the old flag names must update them.

- V1.0.0 — first stable release with breaking changes from v0.2.0 ([474af1f](https://github.com/tbotnz/cisshgo/commit/474af1f0c9a1a15dc97ef5855ac76968fa3a7f84))
### 🐛 Bug Fixes

- Change ParseArgs to return int instead of *int for startingPort ([31205b8](https://github.com/tbotnz/cisshgo/commit/31205b818f5d675d0c5e9c52e8c95dee1fec354a))
- Rename sshlistners package to sshlisteners ([2e7ca9e](https://github.com/tbotnz/cisshgo/commit/2e7ca9e51de2afa5b66930a5309d3886d04a09b9))
- Eliminate FakeDevice data race and ssh.Handle global state (#60) ([7423e32](https://github.com/tbotnz/cisshgo/commit/7423e322ed7f0ee7bbbcd621345b60283b80a678))
- Resolve transcript paths relative to transcript map file location (#76) ([1c34c23](https://github.com/tbotnz/cisshgo/commit/1c34c23c0febdf3f186475114c907e9a86ec8958))

### 👷 CI/CD

- Raise coverage threshold from 60% to 90% ([e174c15](https://github.com/tbotnz/cisshgo/commit/e174c151c0b78b4198b8754443f7d75f1f6f9df8))

### 📚 Documentation

- Add CONTRIBUTING.md with GitHub Flow workflow and transcript guide ([5434278](https://github.com/tbotnz/cisshgo/commit/54342788114d86b5b68e05a17a0448797fb251cb))
- Add package-level doc comments for pkg.go.dev ([a8313e5](https://github.com/tbotnz/cisshgo/commit/a8313e5613c07e3ac8b0ff24d4c9f9d8523c13a3))
- Add git-cliff for automated changelog generation (#66) ([f2386e1](https://github.com/tbotnz/cisshgo/commit/f2386e1d01c5a914833114a8b3c958c6dd0c41e8))
- Promote breaking changes to top section in changelog ([45f40ba](https://github.com/tbotnz/cisshgo/commit/45f40bab62fea63f0eae66109d16ba08fc604485))
- Add emoji section headers and fix breaking change deduplication in changelog ([6e0a8f1](https://github.com/tbotnz/cisshgo/commit/6e0a8f1c7ae203221b933e3aae2af34f305669c8))
- Regenerate changelog with Kong breaking change entry ([b81b88f](https://github.com/tbotnz/cisshgo/commit/b81b88fe63124a6f39187bab91bb5da9c1a1f26f))
- Show commit body for features and breaking changes in changelog ([2f25762](https://github.com/tbotnz/cisshgo/commit/2f25762069443a95d11c8a28d91374cec754f686))
- Add MkDocs Material documentation site ([92505a5](https://github.com/tbotnz/cisshgo/commit/92505a54dff89c24c322c544684a7b93c0676837))
- Fix CLI flag formats and scenarios documentation ([9c6908c](https://github.com/tbotnz/cisshgo/commit/9c6908ca3131df6c5e57e490e002cd79cfbbe234))
- Polish documentation with improvements from technical review ([b436db3](https://github.com/tbotnz/cisshgo/commit/b436db3d1ed8cb840eaa3a359de1cbc5f775ddf6))
- Fix remaining flag format inconsistencies ([eb82291](https://github.com/tbotnz/cisshgo/commit/eb82291c57bc8cc02a423421e295bbba382cfc83))
- Add comprehensive usage examples for all features ([206e530](https://github.com/tbotnz/cisshgo/commit/206e5305de90e9f3f99c7cefa42dc9e8931a5861))
- Add migration guide for v0.2.0 to v1.0.0 (#95) ([bf76cb4](https://github.com/tbotnz/cisshgo/commit/bf76cb4d6baa4c9f63ce04a1f0b8e7420040f463))

### 🔧 Refactoring

- Change TranscriptMap.Platforms from list-of-maps to map (#62) ([8814135](https://github.com/tbotnz/cisshgo/commit/881413506fcd3ea4c635eec802df19162af8a92f))
- Replace flag package with Kong for CLI argument parsing (#64) ([0391eba](https://github.com/tbotnz/cisshgo/commit/0391ebab3aef81a13d2349fb2f2686ab59009f72))
- Split utils package into config, transcript, and cmdmatch ([cf55700](https://github.com/tbotnz/cisshgo/commit/cf55700c1a14c95a8d9b8ef8e6c7458d1b70808e))

### 🚀 Features

- Add graceful shutdown via context.Context (#61) ([eb120de](https://github.com/tbotnz/cisshgo/commit/eb120dede28b2686dc0ac4f4562e2005373956fc))

- GenericListener now accepts context.Context and shuts down cleanly
  on cancellation via ssh.Server.Shutdown
- run() uses sync.WaitGroup to wait for all listeners to exit
- main() wires up signal.NotifyContext for SIGINT/SIGTERM
- Update tests for new signatures; add shutdown verification test

- Implement inventory system for multi-device topology management (#67) ([5a2b784](https://github.com/tbotnz/cisshgo/commit/5a2b7848d57850af4218132dfc85f4a2e9e89947))

- Add Inventory/InventoryEntry types and LoadInventory() to utils
- Add --inventory and --platform flags to CLI struct
- InitGeneric now reads vendor from transcript map (drops vendor param)
- InitGeneric returns error for unknown platform
- run() uses inventory when provided, falls back to --platform/--listeners
- Add iosxr platform to transcript_map.yaml with show version transcript
- Add transcripts/inventory_example.yaml demonstrating multi-device usage

- Add device transcript library for 5 additional platforms (#69) ([dd8b6ce](https://github.com/tbotnz/cisshgo/commit/dd8b6ce3fd5db9efa42291cbf3ad2b054e79f402))

Add show version, show ip interface brief, and show running-config
transcripts sourced from NTC Templates test fixtures for:
- Cisco IOS (ISR4321, 15.6)
- Cisco ASA (ASAv, 9.12)
- Cisco NX-OS (N9K, 9.3)
- Arista EOS (DCS-7050CX3, 4.27)
- Juniper Junos (MX240, 21.4)

Also update inventory_example.yaml to demonstrate all 7 platforms.

- Add ENV variable support for all CLI flags (#75) ([61dfe19](https://github.com/tbotnz/cisshgo/commit/61dfe190224e3923f2c315126229f3cc5debab01))

Each flag can now be set via environment variable (CLI takes precedence):
  CISSHGO_LISTENERS, CISSHGO_STARTING_PORT, CISSHGO_TRANSCRIPT_MAP,
  CISSHGO_PLATFORM, CISSHGO_INVENTORY

- Validate transcript map paths at startup before spawning listeners (#77) ([eaeaf48](https://github.com/tbotnz/cisshgo/commit/eaeaf48098a92f05a4f88595646f4859a18ef79f))

Add ValidateTranscriptMap() that checks all command_transcripts paths
exist on disk before any listeners are spawned. Reports all missing
files in a single error rather than failing on the first.

- Implement scenario-based stateful command responses (#79) ([9419d11](https://github.com/tbotnz/cisshgo/commit/9419d111eefd02a3fe39a81eebb86e5f920e3348))

* feat: Implement scenario-based stateful command responses

Add scenario support to transcript_map.yaml — a scenario defines an
ordered sequence of (command, transcript) pairs layered on top of a
platform. Each SSH session gets its own sequence pointer that advances
as commands match the next expected step; non-matching commands fall
through to normal command_transcripts behavior.

- Add --version flag to display version information (#88) ([2f35a9b](https://github.com/tbotnz/cisshgo/commit/2f35a9b2fe32ad4b0e56456cce5f0457ca3ae800))

- Add Version field to CLI struct with kong.VersionFlag
- Add version, commit, date variables set via ldflags
- Update Kong parser to include version vars
- Add version flag to CLI reference documentation
- GoReleaser already configured with ldflags

- Add username field to platform config for SSH auth enforcement ([5d68f9b](https://github.com/tbotnz/cisshgo/commit/5d68f9b250d8933d0fafa79bdb48744e582db5bd))

- Add Username field to transcript.Platform and fakedevices.FakeDevice
- InitGeneric populates Username from platform config
- PasswordHandler enforces username when set (any username accepted if empty)
- Update Junos platform entry with username: admin
- Add TestGenericListener_UsernameEnforcement
- Document username field in configuration.md and transcripts.md

- Add prompt_format field for flexible prompt construction (#93) ([f0bcdd8](https://github.com/tbotnz/cisshgo/commit/f0bcdd819990948124311ac63c6582e421c11110))

- Add PromptFormat to transcript.Platform and fakedevices.FakeDevice
- Add buildPrompt() helper using strings.NewReplacer for {hostname},
  {username}, {context} variables
- Replace all t.SetPrompt(fd.Hostname+...) calls with buildPrompt()
- Update Junos platform with prompt_format: '{username}@{hostname}{context}'
- Add TestBuildPrompt_Default, TestBuildPrompt_Format, TestHandler_JunosStylePrompt
- Document prompt_format in configuration.md

- Support multi-line prompts via context_prefix_lines (#94) ([a815dfa](https://github.com/tbotnz/cisshgo/commit/a815dfa8a85735b70d67a5e31f6d8556fd1b5ade))

Tested that golang.org/x/term handles \n in prompt strings correctly.
Uses the additive context_prefix_lines map (keyed by context value) to
prepend a line above the prompt for specific contexts.

- Add ContextPrefixLines to transcript.Platform and fakedevices.FakeDevice
- Add devicePrompt() helper wrapping buildPrompt() with prefix line lookup
- Update buildPrompt() to accept and prepend prefixLine when non-empty
- Update Junos platform with context_prefix_lines: '#': '[edit]'
- Add TestBuildPrompt_PrefixLine
- Document context_prefix_lines in configuration.md

- Add transcript map migration script for v0.2.0 to v1.0.0 (#96) ([88b0cd7](https://github.com/tbotnz/cisshgo/commit/88b0cd70238706896534456ec583c85bcd6f24d9))

scripts/migrate_transcript_map.py converts the platforms list-of-maps
schema to the v1.0.0 map format. Handles already-migrated files,
empty platforms, and missing platforms key gracefully.

Also updates migration guide to reference the script.


### 🧪 Testing

- Add transcript_map.yaml integrity test (#74) ([5f5c55f](https://github.com/tbotnz/cisshgo/commit/5f5c55f36d0fddb707b7038ea400486e97e298eb))
- Add coverage for inventory branch, ambiguous context, and TranscriptReader error (#78) ([945624d](https://github.com/tbotnz/cisshgo/commit/945624d653f9fa63aa028a01c0e7ca9c6c42e62f))
- Add integration tests for run() listener spawning (#81) ([02d481b](https://github.com/tbotnz/cisshgo/commit/02d481ba910c6a9a48b3329fb4ece78d91fdc8b4))

### 🧹 Chores

- Update master references to main ([3cca7f4](https://github.com/tbotnz/cisshgo/commit/3cca7f41e220ef409fe6e1baae48a649b0072d80))
- Fix Go Report Card issues (#55) ([576dad4](https://github.com/tbotnz/cisshgo/commit/576dad4956a7215eb869ca5b04e0ca8b9f118c02))

## [0.2.0] - 2026-02-26

### ⚠️ Breaking Changes

- Replace log.Fatal with error returns, achieve 93% coverage (#39) ([585a431](https://github.com/tbotnz/cisshgo/commit/585a431d6f0060b31ce984e1c87c03bbac9a6e53))

* refactor!: Replace log.Fatal with error returns, achieve 93% coverage

### 👷 CI/CD

- Add GitHub Actions test workflow (#36) ([d5a3918](https://github.com/tbotnz/cisshgo/commit/d5a3918f8ced5398d7317d33fda506d49bd6eb2f))
- Add coverage reporting and threshold enforcement (#37) ([94e305f](https://github.com/tbotnz/cisshgo/commit/94e305f1bc4d439c1e1b4b9b7bd9c48d9a8ba204))
- Switch Docker registry from Docker Hub to GHCR ([04d10a8](https://github.com/tbotnz/cisshgo/commit/04d10a86cd5b1d999f9557febd36d7ef749e7340))
- Install syft for SBOM generation in release workflow ([1d14912](https://github.com/tbotnz/cisshgo/commit/1d1491202a8452fae5c89f120325f9f1445d3f54))

### 📚 Documentation

- Add MIT License ([37b4c79](https://github.com/tbotnz/cisshgo/commit/37b4c79519fc6a1f60c9fdc9bd6dc38c5afac7eb))
- Add GitHub Actions release workflow and update README ([096074d](https://github.com/tbotnz/cisshgo/commit/096074da38853c73f2c42a192a8cb07e222be614))
- Add standard Go ecosystem badges to README (#41) ([12ba87f](https://github.com/tbotnz/cisshgo/commit/12ba87f6613fafd03d119dded53fc3c749df10e5))

### 🚀 Features

- Add Docker release support and migrate to goreleaser v2 ([0697e99](https://github.com/tbotnz/cisshgo/commit/0697e993ad027b65556853dcfe0357529774032f))

- Migrate .goreleaser.yml from v0 to v2 format
- Add Docker multi-arch builds (amd64, arm64)
- Add SBOM generation for security compliance
- Add build metadata (version, commit, date) to binaries
- Add Dockerfile.goreleaser for automated releases
- Modernize standard Dockerfile to Go 1.26 with multi-stage build
- Add MIT License file for archive inclusion
- Comment out UPX compression (not recommended for Go 1.26)
- Add dist/ to .gitignore

- Add SSH exec mode support (#42) ([b15019e](https://github.com/tbotnz/cisshgo/commit/b15019efce42ad2d3440238d95762b2dee472c8d))

Handle exec requests (e.g., ssh host "show version") in addition
to interactive shell sessions. Check s.RawCommand() at the top of
the handler — if non-empty, process the command, write output, and
exit. Abbreviated command matching works in exec mode.


### 🧪 Testing

- Add unit tests to achieve 75.5% coverage (#38) ([e5c64d4](https://github.com/tbotnz/cisshgo/commit/e5c64d4c8fb88ac891c690253abd3652faf3a24d))

## [0.0.1] - 2020-09-01


