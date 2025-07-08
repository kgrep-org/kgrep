# Integration Tests

This directory contains integration tests for kgrep that require a real Kubernetes cluster to run. These tests complement the unit tests by testing scenarios that cannot be properly tested with mocked Kubernetes clients.

## What These Tests Cover

The integration tests cover scenarios that require interaction with a real Kubernetes cluster.

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
