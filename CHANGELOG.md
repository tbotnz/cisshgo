# Changelog

All notable changes to this project will be documented in this file.
## [Unreleased]

### Bug Fixes

- Change ParseArgs to return int instead of *int for startingPort ([31205b8](https://github.com/tbotnz/cisshgo/commit/31205b818f5d675d0c5e9c52e8c95dee1fec354a))
- Rename sshlistners package to sshlisteners ([2e7ca9e](https://github.com/tbotnz/cisshgo/commit/2e7ca9e51de2afa5b66930a5309d3886d04a09b9))
- Eliminate FakeDevice data race and ssh.Handle global state (#60) ([7423e32](https://github.com/tbotnz/cisshgo/commit/7423e322ed7f0ee7bbbcd621345b60283b80a678))

### CI/CD

- Raise coverage threshold from 60% to 90% ([e174c15](https://github.com/tbotnz/cisshgo/commit/e174c151c0b78b4198b8754443f7d75f1f6f9df8))

### Chores

- Update master references to main ([3cca7f4](https://github.com/tbotnz/cisshgo/commit/3cca7f41e220ef409fe6e1baae48a649b0072d80))
- Fix Go Report Card issues (#55) ([576dad4](https://github.com/tbotnz/cisshgo/commit/576dad4956a7215eb869ca5b04e0ca8b9f118c02))

### Documentation

- Add CONTRIBUTING.md with GitHub Flow workflow and transcript guide ([5434278](https://github.com/tbotnz/cisshgo/commit/54342788114d86b5b68e05a17a0448797fb251cb))
- Add package-level doc comments for pkg.go.dev ([a8313e5](https://github.com/tbotnz/cisshgo/commit/a8313e5613c07e3ac8b0ff24d4c9f9d8523c13a3))
- Add git-cliff for automated changelog generation (#66) ([f2386e1](https://github.com/tbotnz/cisshgo/commit/f2386e1d01c5a914833114a8b3c958c6dd0c41e8))

### Features

- Add graceful shutdown via context.Context (#61) ([eb120de](https://github.com/tbotnz/cisshgo/commit/eb120dede28b2686dc0ac4f4562e2005373956fc))

### Refactoring

- Change TranscriptMap.Platforms from list-of-maps to map (#62) ([8814135](https://github.com/tbotnz/cisshgo/commit/881413506fcd3ea4c635eec802df19162af8a92f))
- Replace flag package with Kong for CLI argument parsing (#64) ([0391eba](https://github.com/tbotnz/cisshgo/commit/0391ebab3aef81a13d2349fb2f2686ab59009f72))

## [0.2.0] - 2026-02-26
### Breaking Changes

- Replace log.Fatal with error returns, achieve 93% coverage (#39) ([585a431](https://github.com/tbotnz/cisshgo/commit/585a431d6f0060b31ce984e1c87c03bbac9a6e53))

### CI/CD

- Add GitHub Actions test workflow (#36) ([d5a3918](https://github.com/tbotnz/cisshgo/commit/d5a3918f8ced5398d7317d33fda506d49bd6eb2f))
- Add coverage reporting and threshold enforcement (#37) ([94e305f](https://github.com/tbotnz/cisshgo/commit/94e305f1bc4d439c1e1b4b9b7bd9c48d9a8ba204))
- Switch Docker registry from Docker Hub to GHCR ([04d10a8](https://github.com/tbotnz/cisshgo/commit/04d10a86cd5b1d999f9557febd36d7ef749e7340))
- Install syft for SBOM generation in release workflow ([1d14912](https://github.com/tbotnz/cisshgo/commit/1d1491202a8452fae5c89f120325f9f1445d3f54))

### Documentation

- Add MIT License ([37b4c79](https://github.com/tbotnz/cisshgo/commit/37b4c79519fc6a1f60c9fdc9bd6dc38c5afac7eb))
- Add GitHub Actions release workflow and update README ([096074d](https://github.com/tbotnz/cisshgo/commit/096074da38853c73f2c42a192a8cb07e222be614))
- Add standard Go ecosystem badges to README (#41) ([12ba87f](https://github.com/tbotnz/cisshgo/commit/12ba87f6613fafd03d119dded53fc3c749df10e5))

### Features

- Add Docker release support and migrate to goreleaser v2 ([0697e99](https://github.com/tbotnz/cisshgo/commit/0697e993ad027b65556853dcfe0357529774032f))
- Add SSH exec mode support (#42) ([b15019e](https://github.com/tbotnz/cisshgo/commit/b15019efce42ad2d3440238d95762b2dee472c8d))

### Refactoring

- Replace log.Fatal with error returns, achieve 93% coverage (#39) ([585a431](https://github.com/tbotnz/cisshgo/commit/585a431d6f0060b31ce984e1c87c03bbac9a6e53))

### Testing

- Add unit tests to achieve 75.5% coverage (#38) ([e5c64d4](https://github.com/tbotnz/cisshgo/commit/e5c64d4c8fb88ac891c690253abd3652faf3a24d))

## [0.0.1] - 2020-09-01


