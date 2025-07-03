.PHONY: build build-dev build-release clean test

# Default build (development version)
build:
	go build -o kgrep

# Development build with "dev" version
build-dev:
	go build -ldflags "-X github.com/hbelmiro/kgrep/cmd.Version=dev -X github.com/hbelmiro/kgrep/cmd.BuildTime=$(shell date -u '+%Y-%m-%d_%H:%M:%S') -X github.com/hbelmiro/kgrep/cmd.CommitHash=$(shell git rev-parse --short HEAD 2>/dev/null || echo 'unknown')" -o kgrep

# Release build with git tag version
build-release:
	@if [ -z "$(VERSION)" ]; then \
		echo "Error: VERSION is required. Usage: make build-release VERSION=1.0.0"; \
		exit 1; \
	fi
	go build -ldflags "-X github.com/hbelmiro/kgrep/cmd.Version=$(VERSION) -X github.com/hbelmiro/kgrep/cmd.BuildTime=$(shell date -u '+%Y-%m-%d_%H:%M:%S') -X github.com/hbelmiro/kgrep/cmd.CommitHash=$(shell git rev-parse --short HEAD 2>/dev/null || echo 'unknown')" -o kgrep

# Build with git describe (includes dirty state)
build-git:
	go build -ldflags "-X github.com/hbelmiro/kgrep/cmd.Version=$(shell git describe --tags --always --dirty 2>/dev/null || echo 'dev') -X github.com/hbelmiro/kgrep/cmd.BuildTime=$(shell date -u '+%Y-%m-%d_%H:%M:%S') -X github.com/hbelmiro/kgrep/cmd.CommitHash=$(shell git rev-parse --short HEAD 2>/dev/null || echo 'unknown')" -o kgrep

# Release build based on latest Git tag
build-tag:
	@LATEST_TAG=$$(git describe --tags --abbrev=0 2>/dev/null || echo "no-tags"); \
	if [ "$$LATEST_TAG" = "no-tags" ]; then \
		echo "Error: No Git tags found. Create a tag first: git tag v1.0.0"; \
		exit 1; \
	fi; \
	echo "Building release version: $$LATEST_TAG"; \
	go build -ldflags "-X github.com/hbelmiro/kgrep/cmd.Version=$$LATEST_TAG -X github.com/hbelmiro/kgrep/cmd.BuildTime=$(shell date -u '+%Y-%m-%d_%H:%M:%S') -X github.com/hbelmiro/kgrep/cmd.CommitHash=$(shell git rev-parse --short HEAD 2>/dev/null || echo 'unknown')" -o kgrep

# Build for specific tag
build-tag-version:
	@if [ -z "$(TAG)" ]; then \
		echo "Error: TAG is required. Usage: make build-tag-version TAG=v1.0.0"; \
		exit 1; \
	fi; \
	echo "Building version: $(TAG)"; \
	go build -ldflags "-X github.com/hbelmiro/kgrep/cmd.Version=$(TAG) -X github.com/hbelmiro/kgrep/cmd.BuildTime=$(shell date -u '+%Y-%m-%d_%H:%M:%S') -X github.com/hbelmiro/kgrep/cmd.CommitHash=$(shell git rev-parse --short HEAD 2>/dev/null || echo 'unknown')" -o kgrep

# Cross-compilation targets for specific tag
build-tag-version-linux-amd64:
	@if [ -z "$(TAG)" ]; then \
		echo "Error: TAG is required. Usage: make build-tag-version-linux-amd64 TAG=v1.0.0"; \
		exit 1; \
	fi; \
	echo "Building Linux amd64 version: $(TAG)"; \
	GOOS=linux GOARCH=amd64 go build -ldflags "-X github.com/hbelmiro/kgrep/cmd.Version=$(TAG) -X github.com/hbelmiro/kgrep/cmd.BuildTime=$(shell date -u '+%Y-%m-%d_%H:%M:%S') -X github.com/hbelmiro/kgrep/cmd.CommitHash=$(shell git rev-parse --short HEAD 2>/dev/null || echo 'unknown')" -o kgrep-amd64

build-tag-version-linux-arm64:
	@if [ -z "$(TAG)" ]; then \
		echo "Error: TAG is required. Usage: make build-tag-version-linux-arm64 TAG=v1.0.0"; \
		exit 1; \
	fi; \
	echo "Building Linux arm64 version: $(TAG)"; \
	GOOS=linux GOARCH=arm64 go build -ldflags "-X github.com/hbelmiro/kgrep/cmd.Version=$(TAG) -X github.com/hbelmiro/kgrep/cmd.BuildTime=$(shell date -u '+%Y-%m-%d_%H:%M:%S') -X github.com/hbelmiro/kgrep/cmd.CommitHash=$(shell git rev-parse --short HEAD 2>/dev/null || echo 'unknown')" -o kgrep-arm64

build-tag-version-windows-amd64:
	@if [ -z "$(TAG)" ]; then \
		echo "Error: TAG is required. Usage: make build-tag-version-windows-amd64 TAG=v1.0.0"; \
		exit 1; \
	fi; \
	echo "Building Windows amd64 version: $(TAG)"; \
	GOOS=windows GOARCH=amd64 go build -ldflags "-X github.com/hbelmiro/kgrep/cmd.Version=$(TAG) -X github.com/hbelmiro/kgrep/cmd.BuildTime=$(shell date -u '+%Y-%m-%d_%H:%M:%S') -X github.com/hbelmiro/kgrep/cmd.CommitHash=$(shell git rev-parse --short HEAD 2>/dev/null || echo 'unknown')" -o kgrep-amd64.exe

build-tag-version-windows-arm64:
	@if [ -z "$(TAG)" ]; then \
		echo "Error: TAG is required. Usage: make build-tag-version-windows-arm64 TAG=v1.0.0"; \
		exit 1; \
	fi; \
	echo "Building Windows arm64 version: $(TAG)"; \
	GOOS=windows GOARCH=arm64 go build -ldflags "-X github.com/hbelmiro/kgrep/cmd.Version=$(TAG) -X github.com/hbelmiro/kgrep/cmd.BuildTime=$(shell date -u '+%Y-%m-%d_%H:%M:%S') -X github.com/hbelmiro/kgrep/cmd.CommitHash=$(shell git rev-parse --short HEAD 2>/dev/null || echo 'unknown')" -o kgrep-arm64.exe

build-tag-version-darwin-amd64:
	@if [ -z "$(TAG)" ]; then \
		echo "Error: TAG is required. Usage: make build-tag-version-darwin-amd64 TAG=v1.0.0"; \
		exit 1; \
	fi; \
	echo "Building macOS amd64 version: $(TAG)"; \
	GOOS=darwin GOARCH=amd64 go build -ldflags "-X github.com/hbelmiro/kgrep/cmd.Version=$(TAG) -X github.com/hbelmiro/kgrep/cmd.BuildTime=$(shell date -u '+%Y-%m-%d_%H:%M:%S') -X github.com/hbelmiro/kgrep/cmd.CommitHash=$(shell git rev-parse --short HEAD 2>/dev/null || echo 'unknown')" -o kgrep-amd64

build-tag-version-darwin-arm64:
	@if [ -z "$(TAG)" ]; then \
		echo "Error: TAG is required. Usage: make build-tag-version-darwin-arm64 TAG=v1.0.0"; \
		exit 1; \
	fi; \
	echo "Building macOS arm64 version: $(TAG)"; \
	GOOS=darwin GOARCH=arm64 go build -ldflags "-X github.com/hbelmiro/kgrep/cmd.Version=$(TAG) -X github.com/hbelmiro/kgrep/cmd.BuildTime=$(shell date -u '+%Y-%m-%d_%H:%M:%S') -X github.com/hbelmiro/kgrep/cmd.CommitHash=$(shell git rev-parse --short HEAD 2>/dev/null || echo 'unknown')" -o kgrep-arm64

# Create a new tag and build release
tag-and-build:
	@if [ -z "$(VERSION)" ]; then \
		echo "Error: VERSION is required. Usage: make tag-and-build VERSION=1.0.0"; \
		exit 1; \
	fi; \
	echo "Creating tag v$(VERSION)..."; \
	git tag v$(VERSION); \
	echo "Building release version v$(VERSION)..."; \
	go build -ldflags "-X github.com/hbelmiro/kgrep/cmd.Version=v$(VERSION) -X github.com/hbelmiro/kgrep/cmd.BuildTime=$(shell date -u '+%Y-%m-%d_%H:%M:%S') -X github.com/hbelmiro/kgrep/cmd.CommitHash=$(shell git rev-parse --short HEAD 2>/dev/null || echo 'unknown')" -o kgrep; \
	echo "Release built successfully! Don't forget to push the tag: git push origin v$(VERSION)"

# List all available tags
list-tags:
	@echo "Available tags:"; \
	git tag --sort=-version:refname | head -10

# Clean build artifacts
clean:
	rm -f kgrep kgrep-amd64 kgrep-arm64 kgrep.exe kgrep-amd64.exe kgrep-arm64.exe

# Run tests
test:
	go test ./...

# Run integration tests (requires a running Kubernetes cluster)
test-integration:
	@echo "Running integration tests..."
	@echo "Make sure you have a running Kubernetes cluster accessible via kubectl"
	cd test/integration && KGREP_INTEGRATION_TESTS=true go test -v -timeout=10m ./...

# Run integration tests with kind (creates a local cluster)
test-integration-kind:
	@echo "Setting up kind cluster for integration tests..."
	@if ! command -v kind > /dev/null 2>&1; then \
		echo "Error: kind is not installed. Please install it first: https://kind.sigs.k8s.io/docs/user/quick-start/#installation"; \
		exit 1; \
	fi
	@if ! command -v kubectl > /dev/null 2>&1; then \
		echo "Error: kubectl is not installed. Please install it first: https://kubernetes.io/docs/tasks/tools/install-kubectl/"; \
		exit 1; \
	fi
	@echo "Creating kind cluster..."
	kind create cluster --name kgrep-integration-test --image kindest/node:v1.31.0 || true
	@echo "Waiting for cluster to be ready..."
	kubectl wait --for=condition=Ready nodes --all --timeout=300s
	@echo "Running integration tests..."
	cd test/integration && KGREP_INTEGRATION_TESTS=true go test -v -timeout=10m ./...
	@echo "Cleaning up kind cluster..."
	kind delete cluster --name kgrep-integration-test || true

# Install development version
install-dev:
	go install -ldflags "-X github.com/hbelmiro/kgrep/cmd.Version=dev -X github.com/hbelmiro/kgrep/cmd.BuildTime=$(shell date -u '+%Y-%m-%d_%H:%M:%S') -X github.com/hbelmiro/kgrep/cmd.CommitHash=$(shell git rev-parse --short HEAD 2>/dev/null || echo 'unknown')"

# Show current version
version:
	@echo "Current version: $(shell ./kgrep version 2>/dev/null || echo 'not built')" 