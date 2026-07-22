# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Comprehensive test suite with multiple test cases
- Makefile for easier development workflow
- GitHub Actions CI/CD pipeline
- Contributing guidelines (CONTRIBUTING.md)
- Code linting configuration (.golangci.yml)
- Detailed documentation and examples in README
- Roadmap with potential future features

### Changed
- Refactored code for better error handling and maintainability
- Improved code structure with separate functions for different operations
- Updated go.mod to use modern Go version (1.21) and latest dependencies
- Enhanced README with open source standards
- Changed file permissions from 0777 to 0644 for security

### Fixed
- Improved error messages with context
- Better handling of edge cases in file processing
- Fixed deprecated ioutil usage (replaced with os package)

## [1.0.0] - 2021-XX-XX

### Added
- Initial release
- Environment variable substitution in Helm values files
- Support for single file processing
- Support for multiple files processing
- Support for directory (recursive) processing
- Support for mixed files and directories
- Support for both `$VAR` and `${VAR}` syntax
- Helm plugin installation support
- Cross-platform support (Linux, macOS)

[Unreleased]: https://github.com/hydeenoble/helm-subenv/compare/v1.0.0...HEAD
[1.0.0]: https://github.com/hydeenoble/helm-subenv/releases/tag/v1.0.0
