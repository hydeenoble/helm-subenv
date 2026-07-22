# Contributing to helm-subenv

First off, thank you for considering contributing to helm-subenv! It's people like you that make helm-subenv such a great tool.

## Code of Conduct

This project and everyone participating in it is governed by our commitment to providing a welcoming and inspiring community for all. Please be respectful and constructive in your interactions.

## How Can I Contribute?

### Reporting Bugs

Before creating bug reports, please check the existing issues to avoid duplicates. When you create a bug report, include as many details as possible:

* **Use a clear and descriptive title**
* **Describe the exact steps to reproduce the problem**
* **Provide specific examples** to demonstrate the steps
* **Describe the behavior you observed** and what behavior you expected
* **Include your environment details** (OS, Helm version, Go version, etc.)

### Suggesting Enhancements

Enhancement suggestions are tracked as GitHub issues. When creating an enhancement suggestion, include:

* **Use a clear and descriptive title**
* **Provide a detailed description** of the suggested enhancement
* **Explain why this enhancement would be useful** to most users
* **List any similar features** in other tools if applicable

### Pull Requests

1. Fork the repository and create your branch from `main`
2. If you've added code that should be tested, add tests
3. If you've changed APIs, update the documentation
4. Ensure the test suite passes
5. Make sure your code follows the existing code style
6. Write a clear commit message

## Development Setup

### Prerequisites

* Go 1.21 or higher
* Make (optional, but recommended)
* Git

### Setting Up Your Development Environment

1. Fork and clone the repository:
   ```bash
   git clone https://github.com/YOUR_USERNAME/helm-subenv.git
   cd helm-subenv
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Build the project:
   ```bash
   make build
   # or
   go build -o bin/subenv .
   ```

4. Run tests:
   ```bash
   make test
   # or
   go test ./...
   ```

### Development Workflow

1. Create a new branch for your feature or bugfix:
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. Make your changes and commit them:
   ```bash
   git add .
   git commit -m "Add your descriptive commit message"
   ```

3. Run tests and ensure they pass:
   ```bash
   make test
   ```

4. Format your code:
   ```bash
   make fmt
   ```

5. Run the linter:
   ```bash
   make vet
   ```

6. Push to your fork:
   ```bash
   git push origin feature/your-feature-name
   ```

7. Create a Pull Request from your fork to the main repository

## Testing

We use Go's built-in testing framework. Tests are located in `*_test.go` files.

### Running Tests

```bash
# Run all tests
make test

# Run tests with coverage
make coverage

# Run tests for a specific package
go test ./cmd/...
```

### Writing Tests

* Write tests for all new functionality
* Ensure tests are deterministic and don't depend on external state
* Use table-driven tests where appropriate
* Mock external dependencies

## Code Style

* Follow standard Go conventions and idioms
* Use `gofmt` to format your code (run `make fmt`)
* Run `go vet` to catch common mistakes (run `make vet`)
* Write clear, descriptive variable and function names
* Add comments for exported functions and complex logic
* Keep functions small and focused

## Commit Messages

* Use the present tense ("Add feature" not "Added feature")
* Use the imperative mood ("Move cursor to..." not "Moves cursor to...")
* Limit the first line to 72 characters or less
* Reference issues and pull requests liberally after the first line

Example:
```
Add support for default values in environment variables

- Implement ${VAR:-default} syntax
- Add tests for default value functionality
- Update documentation

Fixes #123
```

## Documentation

* Update the README.md if you change functionality
* Add comments to exported functions
* Update examples if behavior changes
* Keep documentation clear and concise

## Release Process

Releases are handled by maintainers using GitHub releases and GoReleaser:

1. Update version in `plugin.yaml`
2. Create and push a new tag: `git tag -a vX.Y.Z -m "Release vX.Y.Z"`
3. Push the tag: `git push origin vX.Y.Z`
4. GitHub Actions will automatically build and create a release

## Questions?

Feel free to open an issue with your question or reach out to the maintainers.

## License

By contributing, you agree that your contributions will be licensed under the MIT License.
