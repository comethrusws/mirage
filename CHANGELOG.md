# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- GoReleaser configuration for multi-platform releases
- Homebrew tap support
- GitHub Actions CI/CD workflows
- Installation script for one-liner setup
- Comprehensive project documentation

### Changed
- Removed all code comments for cleaner codebase
- Improved code organization and naming

## [0.1.0] - 2026-01-01

### Added
- HTTP proxy server with request forwarding
- YAML-based scenario configuration
- Request/response mocking with path, method, and header matching
- Traffic recording to JSON files
- CLI commands (`start`, `record`, `scenarios`, `replay`)
- Web dashboard for request monitoring
- Scenario toggle functionality
- Request delay simulation
- Glob pattern matching for paths

### Technical
- Built with Go 1.21+
- Uses Cobra for CLI
- Gorilla Mux for routing
- YAML v3 for configuration

[Unreleased]: https://github.com/comethrusws/mirage/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/comethrusws/mirage/releases/tag/v0.1.0
