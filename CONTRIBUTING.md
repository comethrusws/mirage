# Contributing to Mirage

Thank you for your interest in contributing to Mirage! This document provides guidelines for contributing to the project.

## Getting Started

1. Fork the repository
2. Clone your fork: `git clone https://github.com/YOUR_USERNAME/mirage.git`
3. Create a feature branch: `git checkout -b feature/amazing-feature`
4. Make your changes
5. Run tests: `make test`
6. Commit your changes: `git commit -m 'feat: add amazing feature'`
7. Push to the branch: `git push origin feature/amazing-feature`
8. Open a Pull Request

## Development Setup

### Prerequisites
- Go 1.21 or higher
- Make

### Building from Source
```bash
make build
```

### Running Tests
```bash
make test
```

### Running Locally
```bash
make dev
```

## Code Style

- Follow standard Go conventions
- Run `go fmt` before committing
- Write meaningful commit messages following [Conventional Commits](https://www.conventionalcommits.org/)
- No comments in code - use clear naming instead

## Commit Message Format

```
<type>: <description>

[optional body]
[optional footer]
```

Types:
- `feat:` New feature
- `fix:` Bug fix
- `docs:` Documentation changes
- `refactor:` Code refactoring
- `test:` Adding tests
- `chore:` Maintenance tasks
- `build:` Build system changes

## Pull Request Process

1. Update the README.md with details of changes if applicable
2. Update the CHANGELOG.md following the existing format
3. Ensure all tests pass
4. Request review from maintainers
5. Squash commits before merging

## Reporting Bugs

Use GitHub Issues with the bug report template. Include:
- Description of the bug
- Steps to reproduce
- Expected behavior
- Actual behavior
- Environment details (OS, Go version)

## Feature Requests

Use GitHub Issues with the feature request template. Include:
- Clear description of the feature
- Use case / motivation
- Possible implementation approach

## Questions?

Join discussions in GitHub Discussions or open an issue.

## License

By contributing, you agree that your contributions will be licensed under the MIT License.
