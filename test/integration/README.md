# Integration Tests

This directory contains integration tests for kgrep that require a real Kubernetes cluster to run. These tests complement the unit tests by testing scenarios that cannot be properly tested with mocked Kubernetes clients.

## What These Tests Cover

The integration tests cover scenarios that require interaction with a real Kubernetes cluster:

### Core Resource Tests
1. **End-to-End Resource Searching**: Tests actual resource creation and searching in a real cluster
2. **Cross-Namespace Operations**: Tests searching across multiple namespaces
3. **Real Log Aggregation**: Tests log searching from actual running pods
4. **Resource Type Interactions**: Tests with different Kubernetes resource types (ConfigMaps, Secrets, Pods, ServiceAccounts)
5. **CLI Command Integration**: Tests the actual kgrep CLI commands against real resources
6. **Multi-Resource Scenarios**: Tests complex scenarios with multiple resource types

### All-Namespaces Tests
7. **All-Namespaces ConfigMaps**: Tests --all-namespaces flag with ConfigMap resources across multiple namespaces
8. **All-Namespaces Secrets**: Tests --all-namespaces flag with Secret resources across multiple namespaces
9. **All-Namespaces Pods**: Tests --all-namespaces flag with Pod resources across multiple namespaces
10. **All-Namespaces ServiceAccounts**: Tests --all-namespaces flag with ServiceAccount resources across multiple namespaces
11. **All-Namespaces Resources**: Tests --all-namespaces flag with generic resource command across multiple namespaces
12. **All-Namespaces Flag Validation**: Tests mutual exclusion of --all-namespaces and --namespace flags
13. **All-Namespaces No Results**: Tests --all-namespaces behavior when no results are found
14. **All-Namespaces Display Format**: Tests that namespace information is properly displayed in output format

### Custom Resource Tests
15. **Custom Resource Auto-Discovery**: Tests that kgrep can discover and search custom resources without needing `--api-version` flag
16. **Resource Name Format Compatibility**: Tests the 4 different resource name formats (plural, short name, resource.group, Kind)
17. **Resource.group Fallback Logic**: Tests fallback when `resource.group` format fails and falls back to just resource name
18. **kubectl Compatibility**: Verifies that kgrep behaves identically to kubectl for resource discovery
19. **Custom Resource Error Handling**: Tests meaningful error messages with real API server errors
20. **Custom Resource Content Search**: Tests searching within different parts of custom resources (spec, status, labels, etc.)

## Test Structure

The integration tests create a dedicated test namespace (`kgrep-integration-test`) and:

1. Set up test resources (ConfigMaps, Secrets, Pods, ServiceAccounts)
2. **For custom resource tests**: Apply a test CRD and create custom resources
3. Run kgrep commands against these resources
4. Verify the expected output
5. Clean up resources after each test

## Custom Resource Test Setup

The custom resource tests use a simple test CRD (`TestApplication`) that mimics real-world custom resources:

- **CRD**: `testapplications.test.kgrep.io`
- **Resource Names**: 
  - Plural: `testapplications`
  - Short Name: `tapp`
  - Kind: `TestApplication`
  - Full Name: `testapplications.test.kgrep.io`
- **Test Content**: Contains searchable text in spec, status, labels, and metadata

This setup allows testing of the exact scenarios from GitHub issue #125:
- `testapplications` (plural name)
- `tapp` (short name)
- `testapplications.test.kgrep.io` (resource.group format)
- `TestApplication` (Kind name)

## Running the Tests

### Prerequisites

- Go 1.24+
- A running Kubernetes cluster accessible via `kubectl`
- `kubectl` configured to connect to your cluster

### Option 1: Using Your Existing Cluster

If you have a running Kubernetes cluster (local or remote):

```bash
# Run integration tests using the current kubectl context
make test-integration

# Or directly with go test
cd test/integration
KGREP_INTEGRATION_TESTS=true go test -v -timeout=10m ./...
```

### Option 2: Using kind (Recommended for Local Development)

To run tests with a temporary kind cluster:

```bash
# This will create a kind cluster, run tests, and clean up
make test-integration-kind
```

Prerequisites for kind:
- [kind](https://kind.sigs.k8s.io/docs/user/quick-start/#installation) installed
- [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/) installed
- Docker running

### Option 3: Manual kind Setup

If you want more control over the kind cluster:

```bash
# Create kind cluster
kind create cluster --name kgrep-integration-test --image kindest/node:v1.31.0

# Wait for cluster to be ready
kubectl wait --for=condition=Ready nodes --all --timeout=300s

# Run tests
cd test/integration
KGREP_INTEGRATION_TESTS=true go test -v -timeout=10m ./...

# Clean up
kind delete cluster --name kgrep-integration-test
```

## Environment Variables

- `KGREP_INTEGRATION_TESTS`: Set to `"true"` to enable integration tests (required)
- `KUBECONFIG`: Path to kubeconfig file (optional, defaults to `~/.kube/config`)

## Test Safety

The integration tests are designed to be safe:

1. **Isolated Namespace**: All test resources are created in a dedicated namespace
2. **Resource Cleanup**: Resources are cleaned up after each test
3. **CRD Cleanup**: Test CRDs are automatically removed after custom resource tests
4. **Non-Destructive**: Tests only create test resources and don't modify existing ones
5. **Opt-In**: Tests only run when explicitly enabled via environment variable

## CI/CD Integration

The integration tests run automatically in CI/CD via GitHub Actions:

- **Trigger**: On pushes to `main` branch and pull requests
- **Environment**: Ubuntu with kind-created Kubernetes cluster
- **Isolation**: Each CI run gets a fresh kind cluster
- **Coverage**: Includes both core resource tests and custom resource tests

## Test Development

When adding new integration tests:

1. Follow the existing pattern of setup/test/cleanup
2. Use the dedicated test namespace
3. Include proper cleanup in `defer` statements
4. For custom resource tests, use the provided test CRD
5. Add meaningful assertions and logging
6. Test both success and failure scenarios where applicable

## Troubleshooting

### Tests Skip with "Integration tests disabled"

Make sure to set the environment variable:
```bash
export KGREP_INTEGRATION_TESTS=true
```

### Connection Issues

Verify your kubectl configuration:
```bash
kubectl cluster-info
kubectl get nodes
```

### Resource Cleanup Issues

If tests fail to clean up resources:
```bash
kubectl delete namespace kgrep-integration-test
kubectl delete crd testapplications.test.kgrep.io
```

### Kind Cluster Issues

If kind cluster creation fails:
```bash
kind delete cluster --name kgrep-integration-test
docker system prune -f
```

### Custom Resource Test Issues

If custom resource tests fail:
```bash
# Check if CRD was applied
kubectl get crd testapplications.test.kgrep.io

# Check custom resources
kubectl get testapplications --all-namespaces

# Clean up manually if needed
kubectl delete crd testapplications.test.kgrep.io
```

## Test Coverage Summary

| Test Category | Tests | Description |
|---------------|-------|-------------|
| Core Resources | 7 tests | ConfigMaps, Secrets, Pods, ServiceAccounts, Logs |
| All-Namespaces | 8 tests | --all-namespaces flag functionality across all resource types |
| Custom Resources | 6 tests | Auto-discovery, format compatibility, kubectl compatibility |
| Error Handling | 2 tests | Core resource errors, custom resource errors |
| Multi-Resource | 2 tests | Multiple resource types, cross-namespace |

**Total**: 25 integration tests covering all major functionality that requires real Kubernetes API interaction.
