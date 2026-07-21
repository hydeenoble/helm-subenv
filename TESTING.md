# Testing & Coverage Guide

This document outlines the testing and coverage setup for the helm-subenv project.

## Coverage Threshold

- **Minimum Coverage**: 70%
- This is enforced in CI/CD and locally via the `make coverage-check` target

## Local Testing

### Run Tests
```bash
make test
```
Runs all tests with race detection.

### Generate Coverage Report
```bash
make coverage
```
Runs tests with coverage, checks threshold, and generates an HTML report at `coverage.html`.

### Check Coverage Only
```bash
make coverage-check
```
Runs tests and verifies the coverage meets the 70% threshold. Fails if threshold not met.

### View Coverage Report
```bash
make coverage-report
```
Generates HTML coverage report from existing `coverage.out` file.

## CI/CD Coverage

The GitHub Actions workflow (`.github/workflows/ci.yml`) includes:

1. **Coverage Analysis** - Runs on Ubuntu + Go 1.22
2. **Coverage Threshold Check** - Fails if coverage < 70%
3. **Codecov Integration** - Uploads coverage to Codecov
4. **Coverage Report** - Artifacts uploaded for review

## Configuration Files

### `.codecov.yml`
Codecov configuration specifying:
- Precision: 2 decimal places
- Range: 70-100%
- Ignored files: test files, vendor, mocks
- Requires CI to pass before merge

### `.golangci.yml`
Linting configuration with coverage support and test checking enabled.

### `.gitignore`
Excludes coverage files from version control:
- `coverage.out`
- `coverage.html`
- `*.coverprofile`

## Understanding Coverage Output

When running `make coverage-check`, you'll see:

```
Total coverage: 85.3%
✓ Coverage 85.3% meets threshold
```

The HTML report (`coverage.html`) shows:
- Green: fully covered code
- Red: uncovered code
- Yellow: partially covered code

## Tips for Improving Coverage

1. **Identify uncovered lines**: Open `coverage.html` and look for red/yellow sections
2. **Add unit tests**: Write tests for error cases and edge conditions
3. **Test error paths**: Don't just test happy paths
4. **Use table-driven tests**: Efficient for testing multiple inputs
5. **Check benchmark tests**: Consider adding benchmarks for performance-critical code

## CI/CD Failure Scenarios

Coverage checks fail if:
- Total coverage < 70%
- Codecov upload fails and `fail_ci_if_error: true`
- Tests fail with race detector
- Lint checks fail

## Viewing Coverage in PRs

When you open a pull request:
1. Codecov bot comments with coverage changes
2. Coverage report is available as a GitHub Actions artifact
3. Coverage badge can be added to README (see example below)

## Coverage Badge

To add a coverage badge to your README:

```markdown
[![codecov](https://codecov.io/gh/hydeenoble/helm-env/branch/main/graph/badge.svg)](https://codecov.io/gh/hydeenoble/helm-env)
```

Replace `hydeenoble/helm-env` with your actual GitHub repository path.
