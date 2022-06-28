# helm-subenv
[![License MIT](https://img.shields.io/badge/license-MIT-blue.svg?style=flat)](LICENSE)
[![GitHub release](https://img.shields.io/github/v/release/hydeenoble/helm-subenv.svg)](https://github.com/hydeenoble/helm-subenv/releases)

This Helm plugin allows you to substitute the environment variables specified in your helm values file with their respective values in the environment from within a CICD pipeline.

## Install

The installation itself is simple as:

```bash
helm plugin install https://github.com/hydeenoble/helm-subenv.git
```
You can install a specific release version:
```bash
helm plugin install https://github.com/hydeenoble/helm-subenv.git --version <release version>
```

To use the plugin, you do not need any special dependencies. The installer will download the latest release with prebuilt binary from [GitHub releases](https://github.com/hydeenoble/helm-subenv/releases).

## Usage

### Single file usage
```bash
helm subenv -f <path to values file>
```

### Multiple files usage
```bash
helm subenv -f <path to values file> -f <path to values file> -f <path to values file>
```

### Directory usage
The plugin can also be used to recursively substitute environment variables in all the files in a specified directory.
```bash
helm subenv -f <path to directory>
```

### Mix files and directories
You can also decide to mix files and directories:
```bash
helm subenv -f <path to values file> -f <path to directory>
```

## Example
Sample helm values file:
```yaml
# values.yaml

image:
  repository: $REGISTRY/$IMAGE_NAME
  tag: $IMAGE_TAG
```
Environment variables configured in your environment (this should most likely be configured with your CI environment): 
```txt
REGISTRY => docker.com
IMAGE_NAME => helm-subenv
IMAGE_TAG => test
```
Substitute Env:
```
helm subenv -f values.yaml
```
Result: 
```yaml
image:
  repository: docker.com/helm-subenv
  tag: test
```
**Note:** If the value of the environment variable does not exist, it will be replaced with an empty string. For instance, from the above example, if `IMAGE_TAG` does not exist as an environment variable in the environment the result would have been: 

```yaml
image:
  repository: docker.com/helm-subenv
  tag:
```

## Uninstall
```
helm plugin remove subenv
```
## License

[MIT](LICENSE)