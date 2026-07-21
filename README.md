# helm-subenv

[![License MIT](https://img.shields.io/badge/license-MIT-blue.svg?style=flat)](LICENSE)
[![GitHub release](https://img.shields.io/github/v/release/hydeenoble/helm-subenv.svg)](https://github.com/hydeenoble/helm-subenv/releases)
[![Artifact Hub](https://img.shields.io/endpoint?url=https://artifacthub.io/badge/repository/helm-subenv)](https://artifacthub.io/packages/search?repo=helm-subenv)
[![CI](https://github.com/hydeenoble/helm-subenv/workflows/CI/badge.svg)](https://github.com/hydeenoble/helm-subenv/actions)
[![codecov](https://codecov.io/gh/hydeenoble/helm-env/branch/main/graph/badge.svg)](https://codecov.io/gh/hydeenoble/helm-env)
[![Go Report Card](https://goreportcard.com/badge/github.com/hydeenoble/helm-env)](https://goreportcard.com/report/github.com/hydeenoble/helm-env)

A powerful Helm plugin that enables environment variable substitution in your Helm values files, perfect for CI/CD pipelines and dynamic configuration management.

## 📋 Table of Contents

- [Features](#-features)
- [Installation](#-installation)
- [Usage](#-usage)
- [Examples](#-examples)
- [Configuration](#-configuration)
- [Development](#-development)
- [Contributing](#-contributing)
- [Roadmap](#-roadmap--potential-features)
- [License](#-license)
- [Support](#-support)

## ✨ Features

- 🔄 **Environment Variable Substitution**: Replace placeholders with actual environment variable values
- 📁 **Multiple File Support**: Process multiple files in a single command
- 📂 **Directory Processing**: Recursively process all files in a directory
- 🔀 **Mixed Mode**: Combine files and directories in one operation
- 🎯 **Flexible Syntax**: Supports both `$VAR` and `${VAR}` syntax
- 🚀 **CI/CD Ready**: Designed for seamless integration with CI/CD pipelines
- ⚡ **Fast & Lightweight**: Minimal dependencies, maximum performance
- 🔒 **Safe Operations**: Validates paths and handles errors gracefully

## 📦 Installation

### Using Helm Plugin Manager

The simplest way to install:

```bash
helm plugin install https://github.com/hydeenoble/helm-subenv.git
```

### Install Specific Version

```bash
helm plugin install https://github.com/hydeenoble/helm-subenv.git --version v1.0.0
```

### Verify Installation

```bash
helm plugin list
```

You should see `subenv` in the list of installed plugins.

## 🚀 Usage

### Basic Syntax

```bash
helm subenv -f <path>
```

### Single File

Process a single values file:

```bash
helm subenv -f values.yaml
```

### Multiple Files

Process multiple files at once:

```bash
helm subenv -f values.yaml -f secrets.yaml -f config.yaml
```

### Directory Processing

Recursively process all files in a directory:

```bash
helm subenv -f ./config-dir
```

### Mixed Mode

Combine files and directories:

```bash
helm subenv -f values.yaml -f ./config-dir -f secrets.yaml
```

### Help

Display help information:

```bash
helm subenv --help
```

## 📚 Examples

### Example 1: Basic Substitution

**Input** (`values.yaml`):
```yaml
image:
  repository: $REGISTRY/$IMAGE_NAME
  tag: $IMAGE_TAG
  pullPolicy: Always
```

**Environment Variables**:
```bash
export REGISTRY=docker.io
export IMAGE_NAME=myapp
export IMAGE_TAG=v1.2.3
```

**Command**:
```bash
helm subenv -f values.yaml
```

**Output** (`values.yaml`):
```yaml
image:
  repository: docker.io/myapp
  tag: v1.2.3
  pullPolicy: Always
```

### Example 2: Braces Syntax

**Input**:
```yaml
database:
  host: ${DB_HOST}
  port: ${DB_PORT}
  name: ${DB_NAME}
  url: postgresql://${DB_HOST}:${DB_PORT}/${DB_NAME}
```

**Environment Variables**:
```bash
export DB_HOST=localhost
export DB_PORT=5432
export DB_NAME=mydb
```

**Result**:
```yaml
database:
  host: localhost
  port: 5432
  name: mydb
  url: postgresql://localhost:5432/mydb
```

### Example 3: CI/CD Pipeline Integration

#### GitHub Actions

```yaml
name: Deploy
on:
  push:
    branches: [main]

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - name: Install Helm
        uses: azure/setup-helm@v3
        
      - name: Install helm-subenv
        run: helm plugin install https://github.com/hydeenoble/helm-subenv.git
        
      - name: Substitute environment variables
        env:
          IMAGE_TAG: ${{ github.sha }}
          REGISTRY: ghcr.io/${{ github.repository }}
        run: helm subenv -f values.yaml
        
      - name: Deploy with Helm
        run: helm upgrade --install myapp ./chart -f values.yaml
```

#### GitLab CI

```yaml
deploy:
  stage: deploy
  image: alpine/helm:latest
  before_script:
    - helm plugin install https://github.com/hydeenoble/helm-subenv.git
  script:
    - export IMAGE_TAG=$CI_COMMIT_SHA
    - export REGISTRY=$CI_REGISTRY_IMAGE
    - helm subenv -f values.yaml
    - helm upgrade --install myapp ./chart -f values.yaml
```

#### Jenkins

```groovy
pipeline {
    agent any
    
    environment {
        IMAGE_TAG = "${env.BUILD_NUMBER}"
        REGISTRY = "docker.io/myorg"
    }
    
    stages {
        stage('Prepare') {
            steps {
                sh 'helm plugin install https://github.com/hydeenoble/helm-subenv.git'
            }
        }
        
        stage('Deploy') {
            steps {
                sh 'helm subenv -f values.yaml'
                sh 'helm upgrade --install myapp ./chart -f values.yaml'
            }
        }
    }
}
```

### Example 4: Multiple Environments

**Directory Structure**:
```
config/
├── base-values.yaml
├── dev/
│   └── values.yaml
├── staging/
│   └── values.yaml
└── prod/
    └── values.yaml
```

**Command**:
```bash
# Process all environment configs
helm subenv -f config/base-values.yaml -f config/$ENVIRONMENT/
```

## ⚙️ Configuration

### Environment Variable Syntax

The plugin supports two syntax styles:

1. **Simple**: `$VARIABLE_NAME`
2. **Braces**: `${VARIABLE_NAME}`

### Behavior Notes

- **Missing Variables**: If an environment variable doesn't exist, it will be replaced with an empty string
- **File Permissions**: Processed files maintain their original permissions
- **Recursive Processing**: When processing directories, all files are processed recursively
- **Error Handling**: The plugin will exit with an error if any file cannot be processed

## 🛠️ Development

### Prerequisites

- Go 1.21 or higher
- Make (optional, but recommended)
- Git

### Building from Source

```bash
# Clone the repository
git clone https://github.com/hydeenoble/helm-subenv.git
cd helm-subenv

# Build the binary
make build

# Or without Make
go build -o bin/subenv .
```

### Running Tests

```bash
# Run all tests
make test

# Run tests with coverage
make coverage

# Run specific tests
go test ./cmd/... -v
```

### Code Quality

```bash
# Format code
make fmt

# Run linter
make lint

# Run vet
make vet

# Run all checks
make all
```

### Local Testing

```bash
# Build and test locally
go build -o subenv .
./subenv -f path/to/test/file.yaml
```

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details.

### Quick Start for Contributors

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Run tests (`make test`)
5. Commit your changes (`git commit -m 'Add amazing feature'`)
6. Push to the branch (`git push origin feature/amazing-feature`)
7. Open a Pull Request

### Code of Conduct

Please be respectful and constructive in all interactions. We're committed to providing a welcoming and inspiring community for all.

## 🗺️ Roadmap & Potential Features

We're constantly working to improve helm-subenv. Here are some features we're considering for future releases:

### Planned Features

- [ ] **Default Values Support**: `${VAR:-default}` syntax for fallback values
- [ ] **Required Variables**: `${VAR:?error message}` to fail if variable is not set
- [ ] **Variable Validation**: Validate that all required variables are set before processing
- [ ] **Dry Run Mode**: Preview changes without modifying files (`--dry-run` flag)
- [ ] **Backup Creation**: Automatically create backups before modifying files (`--backup` flag)
- [ ] **Output to Different File**: Write results to a new file instead of overwriting (`--output` flag)
- [ ] **Variable Listing**: List all variables found in files (`--list-vars` flag)
- [ ] **JSON Support**: Support for JSON configuration files
- [ ] **TOML Support**: Support for TOML configuration files
- [ ] **Environment File Loading**: Load variables from `.env` files (`--env-file` flag)
- [ ] **Variable Prefix Filtering**: Only substitute variables with specific prefixes
- [ ] **Encryption Support**: Decrypt encrypted values before substitution
- [ ] **Template Functions**: Support for basic template functions (uppercase, lowercase, etc.)
- [ ] **Conditional Substitution**: Only substitute if certain conditions are met
- [ ] **Verbose Mode**: Detailed logging of substitutions (`--verbose` flag)
- [ ] **Quiet Mode**: Suppress all output except errors (`--quiet` flag)
- [ ] **Parallel Processing**: Process multiple files concurrently for better performance
- [ ] **Watch Mode**: Continuously watch files for changes and auto-substitute
- [ ] **Integration Tests**: Comprehensive integration test suite
- [ ] **Performance Benchmarks**: Benchmark suite for performance tracking
- [ ] **Docker Image**: Pre-built Docker image for containerized environments
- [ ] **Kubernetes Operator**: Kubernetes operator for automatic substitution
- [ ] **Web UI**: Simple web interface for testing substitutions
- [ ] **VS Code Extension**: Extension for IDE integration
- [ ] **Shell Completion**: Auto-completion for bash, zsh, and fish

### Community Requested Features

Have an idea? [Open an issue](https://github.com/hydeenoble/helm-subenv/issues/new) and let us know!

## 🔧 Troubleshooting

### Common Issues

**Issue**: Plugin installation fails
```bash
# Solution: Ensure you have proper permissions
sudo helm plugin install https://github.com/hydeenoble/helm-subenv.git
```

**Issue**: Variables not being substituted
```bash
# Solution: Verify environment variables are exported
export MY_VAR=value
echo $MY_VAR  # Should print: value
helm subenv -f values.yaml
```

**Issue**: Permission denied when writing files
```bash
# Solution: Check file permissions
chmod 644 values.yaml
helm subenv -f values.yaml
```

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 💬 Support

- 📖 [Documentation](https://github.com/hydeenoble/helm-subenv)
- 🐛 [Issue Tracker](https://github.com/hydeenoble/helm-subenv/issues)
- 💡 [Feature Requests](https://github.com/hydeenoble/helm-subenv/issues/new)
- 📧 [Email](mailto:hydeenoble39@gmail.com)

## 🙏 Acknowledgments

- Built with [Cobra](https://github.com/spf13/cobra) - A Commander for modern Go CLI interactions
- Uses [envsubst](https://github.com/a8m/envsubst) - Environment variable substitution for Go
- Inspired by the need for better CI/CD integration with Helm

## 📊 Project Stats

![GitHub stars](https://img.shields.io/github/stars/hydeenoble/helm-subenv?style=social)
![GitHub forks](https://img.shields.io/github/forks/hydeenoble/helm-subenv?style=social)
![GitHub watchers](https://img.shields.io/github/watchers/hydeenoble/helm-subenv?style=social)

---

**Made with ❤️ by [hydeenoble](https://github.com/hydeenoble)**

If you find this project useful, please consider giving it a ⭐️!
