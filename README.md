# helm-subenv
[![License MIT](https://img.shields.io/badge/license-MIT-blue.svg?style=flat)](LICENSE)
[![GitHub release](https://img.shields.io/github/v/release/hydeenoble/helm-subenv.svg)](https://github.com/hydeenoble/helm-subenv/releases)

This Helm plugin allows you to substitue the environment variables specified in your helm values file with their respective values in the environment from within a CICD pipeline.

## Install

The installation itself is simple as:

```bash
helm plugin install https://github.com/hydeenoble/helm-subenv.git
```
To use the plugin, you do not need any special dependencies. The installer will
download versioned release with prebuilt binary from [github releases](https://github.com/hydeenoble/helm-subenv/releases).

## Example
Sample helm values file:
```yaml
# values.yaml

image:
  repository: $REGISTRY/$IMAGE_NAME
  tag: $IMAGE_TAG
```
Environment variables configured in your environment (this should most likely be configure with you CI environment): 
```txt
REGISTRY => docker.com
IMAGE_NAME => helm-subenv
IMAGE_TAG => test
```

Result: 
```yaml
image:
  repository: docker.com/helm-subenv
  tag: test
```
**Note:** If the value of the environment variable does not exit, it will be replaced with an empty string. For instance, from the above example, if `IMAGE_TAG` does not exist as an environment variable in the environment the result would have been: 

```yaml
image:
  repository: docker.com/helm-subenv
  tag:
```
## License

[MIT](LICENSE)