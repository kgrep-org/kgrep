#!/bin/bash

# Release script for kgrep
# Handles creating Git tags and building release versions

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Function to show usage
show_usage() {
    echo "Usage: $0 [COMMAND] [OPTIONS]"
    echo ""
    echo "Commands:"
    echo "  create VERSION    Create a new release tag and build"
    echo "  build TAG         Build release for existing tag"
    echo "  latest            Build release for latest tag"
    echo "  list              List all available tags"
    echo "  clean             Clean build artifacts"
    echo ""
    echo "Examples:"
    echo "  $0 create 1.0.0"
    echo "  $0 build v1.0.0"
    echo "  $0 latest"
    echo ""
}

# Function to validate version format
validate_version() {
    local version=$1
    if [[ ! $version =~ ^[0-9]+\.[0-9]+\.[0-9]+(-[a-zA-Z0-9.-]+)?(\+[a-zA-Z0-9.-]+)?$ ]]; then
        print_error "Invalid version format: $version"
        print_error "Use semantic versioning: MAJOR.MINOR.PATCH (e.g., 1.0.0)"
        exit 1
    fi
}

# Function to check if tag exists
tag_exists() {
    local tag=$1
    git rev-parse "$tag" >/dev/null 2>&1
}

# Function to check if working directory is clean
check_clean_worktree() {
    if ! git diff-index --quiet HEAD --; then
        print_warning "Working directory has uncommitted changes"
        read -p "Continue anyway? (y/N): " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            print_error "Release cancelled"
            exit 1
        fi
    fi
}

# Function to run tests
run_tests() {
    print_status "Running tests..."
    if ! go test ./...; then
        print_error "Tests failed. Aborting release."
        exit 1
    fi
    print_success "All tests passed"
}

# Function to create release
create_release() {
    local version=$1
    local tag="v$version"
    
    print_status "Creating release for version $version"
    
    # Validate version
    validate_version "$version"
    
    # Check if tag already exists
    if tag_exists "$tag"; then
        print_error "Tag $tag already exists"
        exit 1
    fi
    
    # Check clean worktree
    check_clean_worktree
    
    # Run tests
    run_tests
    
    # Create tag
    print_status "Creating Git tag: $tag"
    git tag "$tag"
    
    # Build release
    print_status "Building release..."
    go build -ldflags "-X github.com/hbelmiro/kgrep/cmd.Version=$tag -X github.com/hbelmiro/kgrep/cmd.BuildTime=$(date -u '+%Y-%m-%d_%H:%M:%S') -X github.com/hbelmiro/kgrep/cmd.CommitHash=$(git rev-parse --short HEAD)" -o kgrep
    
    # Test the built version
    print_status "Testing built version..."
    ./kgrep version
    
    print_success "Release $tag created successfully!"
    print_status "Next steps:"
    echo "  1. Push the tag: git push origin $tag"
    echo "  2. Create a GitHub release for $tag"
    echo "  3. Upload the kgrep binary to the release"
}

# Function to build for existing tag
build_release() {
    local tag=$1
    
    print_status "Building release for tag: $tag"
    
    # Check if tag exists
    if ! tag_exists "$tag"; then
        print_error "Tag $tag does not exist"
        print_status "Available tags:"
        git tag --sort=-version:refname | head -10
        exit 1
    fi
    
    # Checkout the tag
    print_status "Checking out tag: $tag"
    git checkout "$tag"
    
    # Build release
    print_status "Building release..."
    go build -ldflags "-X github.com/hbelmiro/kgrep/cmd.Version=$tag -X github.com/hbelmiro/kgrep/cmd.BuildTime=$(date -u '+%Y-%m-%d_%H:%M:%S') -X github.com/hbelmiro/kgrep/cmd.CommitHash=$(git rev-parse --short HEAD)" -o kgrep
    
    # Test the built version
    print_status "Testing built version..."
    ./kgrep version
    
    print_success "Release built successfully for $tag"
}

# Function to build latest tag
build_latest() {
    local latest_tag=$(git tag --sort=-version:refname | head -1 2>/dev/null || echo "")
    
    if [ -z "$latest_tag" ]; then
        print_error "No tags found in repository"
        print_status "Create a tag first: $0 create 1.0.0"
        exit 1
    fi
    
    print_status "Latest tag: $latest_tag"
    build_release "$latest_tag"
}

# Function to list tags
list_tags() {
    print_status "Available tags:"
    if git tag | wc -l | grep -q "0"; then
        print_warning "No tags found"
        print_status "Create your first tag: $0 create 1.0.0"
    else
        git tag --sort=-version:refname
    fi
}

# Function to clean build artifacts
clean_builds() {
    print_status "Cleaning build artifacts..."
    rm -f kgrep
    print_success "Build artifacts cleaned"
}

# Main script logic
case "${1:-}" in
    "create")
        if [ -z "${2:-}" ]; then
            print_error "Version is required"
            show_usage
            exit 1
        fi
        create_release "$2"
        ;;
    "build")
        if [ -z "${2:-}" ]; then
            print_error "Tag is required"
            show_usage
            exit 1
        fi
        build_release "$2"
        ;;
    "latest")
        build_latest
        ;;
    "list")
        list_tags
        ;;
    "clean")
        clean_builds
        ;;
    "help"|"-h"|"--help")
        show_usage
        ;;
    "")
        print_error "Command is required"
        show_usage
        exit 1
        ;;
    *)
        print_error "Unknown command: $1"
        show_usage
        exit 1
        ;;
esac 