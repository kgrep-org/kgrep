# kgrep

[![License](https://img.shields.io/badge/license-Apache%202.0-blue.svg)](LICENSE)
[![Build Status](https://github.com/hbelmiro/kgrep/actions/workflows/ci.yaml/badge.svg)](https://github.com/hbelmiro/kgrep/actions/workflows/ci.yaml)
[![Latest Release](https://img.shields.io/github/v/release/hbelmiro/kgrep)](https://github.com/hbelmiro/kgrep/releases)

`kgrep` is a command-line utility designed to simplify the process of searching and analyzing logs and resources in Kubernetes. Unlike traditional methods that involve printing resource definitions and grepping through them, `kgrep` allows you to search across multiple logs or resources simultaneously, making it easier to find what you need quickly.

## Key Features

* **Resource Searching**: Search the content of Kubernetes resources such as `ConfigMaps` for specific patterns within designated namespaces.

* **Log Searching**: Inspect logs from a group of pods or entire namespaces, filtering by custom patterns to locate relevant entries.

* **Namespace Specification**: Every search command supports namespace specification, allowing users to focus their queries on particular sections of their Kubernetes cluster.

* **Pattern-based Filtering**: Utilize pattern matching to refine search results, ensuring that only the most pertinent data is returned.

## Installation

### Prerequisites

- [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/) installed and configured to connect to your
  Kubernetes cluster.

### Download the binary and add it to your PATH

Download a release from https://github.com/hbelmiro/kgrep/releases, uncompress it, and add it to your PATH.

#### ⚠️ Unverified app warning on macOS

You can see this warning when trying to run `kgrep` for the first time on macOS.

```
"kgrep" Not Opened
Apple could not verify "kgrep" is free of malware that may harm your Mac or compromise your privacy.
```

![kgrep-not-opened.png](resources/kgrep-not-opened.png)

If you see that, click "Done" and allow `kgrep` to run in macOS settings, like the following screenshot.

![allow-kgrep.png](resources/allow-kgrep.png)

When you try to run it again, you'll see a final warning. Just click "Open Anyway" and it won't warn you anymore.

![open-anyway.png](resources/open-anyway.png)

## Example

To search for `example` in `ConfigMaps` definitions in the `my_namespace` namespace: 

```shell
$ kgrep configmaps -n my_namespace -p "example"
configmaps/example-config-4khgb5fg64[7]:     internal.config.kubernetes.io/previousNames: "example-config-4khgb5fg64"
configmaps/example-config-4khgb5fg64[48]:   name: "example-config-4khgb5fg64"
configmaps/example-config-5fmk4f7h8k[7]:     internal.config.kubernetes.io/previousNames: "example-config-5fmk4f7h8k"
configmaps/example-config-5fmk4f7h8k[57]:   name: "example-config-5fmk4f7h8k"
configmaps/acme-manager-config[104]:     \  frameworks:\n  - \"batch/job\"\n  - \"example.org/mpijob\"\n  - \"acme.io/acmejob\"\
configmaps/acme-manager-config[105]:     \n  - \"acme.io/acmecluster\"\n  - \"jobset.x-k8s.io/jobset\"\n  - \"example.org/mxjob\"\
configmaps/acme-manager-config[106]:     \n  - \"example.org/paddlejob\"\n  - \"example.org/acmejob\"\n  - \"example.org/tfjob\"\
configmaps/acme-manager-config[107]:     \n  - \"example.org/xgboostjob\"\n# - \"pod\"\n  externalFrameworks:\n
```

💡 Type `kgrep --help` to check all the commands.

## Building the project

This project is written in Go and uses the Cobra library for CLI commands.

## Running the application in dev mode

You can run the application in development mode using:
```shell script
go run main.go
```

## Building the application

You can build the application using:
```shell script
go build -o kgrep
```

This will produce a `kgrep` executable in the current directory.

## Installing the application

You can install the application to your GOPATH using:
```shell script
go install
```

## Cross-compiling for different platforms

Go makes it easy to cross-compile for different platforms:

For Linux:
```shell script
GOOS=linux GOARCH=amd64 go build -o kgrep-linux-amd64
```

For macOS:
```shell script
GOOS=darwin GOARCH=amd64 go build -o kgrep-darwin-amd64
```

For Windows:
```shell script
GOOS=windows GOARCH=amd64 go build -o kgrep-windows-amd64.exe
```

## Creating Releases

This project uses a tag-based release system that automatically builds and publishes binaries when you push a Git tag.

### Prerequisites

- Make sure you have the latest changes committed
- Ensure all tests pass: `go test ./...`

### Creating a New Release

1. **Create a release tag and build:**
   ```bash
   ./scripts/release.sh create 1.0.0
   ```
   This will:
   - Run all tests
   - Create a Git tag (e.g., `v1.0.0`)
   - Build the release binary with version information
   - Show you the next steps

2. **Push the tag to trigger the GitHub Actions release:**
   ```bash
   git push origin v1.0.0
   ```

3. **GitHub Actions will automatically:**
   - Build binaries for Linux and macOS (amd64/arm64)
   - Create a GitHub release
   - Upload the binaries to the release

### Building for Existing Tags

If you need to rebuild a release for an existing tag:

```bash
# Build for the latest tag
./scripts/release.sh latest

# Build for a specific tag
./scripts/release.sh build v1.0.0
```

### Development Builds

For development builds:

```bash
# Development build with "dev" version
make build-dev

# Or using the Makefile
make build-dev
```

### Available Release Commands

```bash
# Show all available commands
./scripts/release.sh help

# List all available tags
./scripts/release.sh list

# Clean build artifacts
./scripts/release.sh clean
```

### Version Information

The built binaries include version information that can be displayed with:

```bash
./kgrep version
```

This shows:
- Version number (from Git tag or "dev" for development builds)
- Build timestamp
- Git commit hash

### Makefile Targets

You can also use Makefile targets for building:

```bash
# Development build
make build-dev

# Build for latest tag
make build-tag

# Build for specific tag
make build-tag-version TAG=v1.0.0

# Create new tag and build
make tag-and-build VERSION=1.0.0

# List available tags
make list-tags
```
