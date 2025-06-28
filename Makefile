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
	rm -f kgrep

# Run tests
test:
	go test ./...

# Install development version
install-dev:
	go install -ldflags "-X github.com/hbelmiro/kgrep/cmd.Version=dev -X github.com/hbelmiro/kgrep/cmd.BuildTime=$(shell date -u '+%Y-%m-%d_%H:%M:%S') -X github.com/hbelmiro/kgrep/cmd.CommitHash=$(shell git rev-parse --short HEAD 2>/dev/null || echo 'unknown')"

# Show current version
version:
	@echo "Current version: $(shell ./kgrep version 2>/dev/null || echo 'not built')" 