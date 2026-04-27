# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/).
This project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Fixed
- Now it does graceful shutdown when receives system signal

### Changed
- Environment variable names got prefix to ensure uniqueness

## [v2.3.1]

### Fixed
- Fix bug preventing door open/close actions in Home Assistant integration

## [v2.3.0]

### Added
- Add wire protocol documentation (#17)
- Add "ports-list" command that lists all ports with their types (#20)
- Add port enumeration to SDK (GET_PORTS, GET_TYPE) (#20)
- If `autologin` flag is false in the configuration, Home Assistant integration will not do *PERIODIC* logout-login sequency to renew the token use to communication with the gateway

### Changed
- Updated linter versions and fixed code smells
- 'groups group' subcommand has been renamed to 'groups list'
- 'users group' subcommand has been renamed to 'users list'
- README updated
- Increased gateway authorization token lifecycle period (automatic logout-login sequence happens less frequently: only in every 6 hours)

### Fixed
- Drain stale LOGOUT responses during login (#19)
- Fixed debug logging: byte array was logged as an ascii string instead of hex string
- Because of unchecked type cast, panic occurred when a command failed, for example, because of permission denied response
- Home Assistant integration now use token and its timestamp from the config file and write renewd token back into it

## [v2.2.3]

### Added
- Warning message when no device port specified 

### Fixed
- Addressed gosec linter warnings related to unsafe integer-to-byte conversions.

## [v2.2.2]

### Added
- Maintain an increased status report frequency while the door is in motion, regardless of the trigger source.
- Add `doorStatusSupported` option to Home Assistant integration.
- Perform Hormann `logout` when terminating Home Assistant integration process.
- Add multi door support.

### Changed

### Fixed
- Fixed local docker build (build argument was missing).
- Fixed potential resource leak

## [v2.2.1]

### Added
- Introduce CHANGELOG.md

### Changed
- Release workflow changed to use goreleaser
- updated tools: golang, golangci-lint, gosec

## [v2.2.0] - 2025-08-12

### Added
- add fast update + dependency updates

## [v2.1.0] - 2025-07-22

### Added
- Improve Home Assistant integration with auto discovery.
- Add docker compose support with example yaml file.

## [v2.0.0] - 2025-07-20

### Added
- Add basic Home Assistant integration.

## [v1.0.0] - 2024-03-29

### Added
- Door handling features are tested on real Hörmann BiSecur gateway.

## [v0.0.1-alpha3] - 2024-02-04

### Added
- ?? 

## [v0.0.1-alpha2] - 2024-01-16 

### Added
- refactor ping: separate functionality of cli client and sdk properly

## [v0.0.1-alpha1] - 2024-01-11

### Added
- Add version string
