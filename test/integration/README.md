# Integration Tests

This directory contains integration tests for kgrep that require a real Kubernetes cluster to run. These tests complement the unit tests by testing scenarios that cannot be properly tested with mocked Kubernetes clients.

## What These Tests Cover

The integration tests cover scenarios that require interaction with a real Kubernetes cluster:

1. **End-to-End Resource Searching**: Tests actual resource creation and searching in a real cluster
2. **Cross-Namespace Operations**: Tests searching across multiple namespaces
3. **Real Log Aggregation**: Tests log searching from actual running pods
4. **Resource Type Interactions**: Tests with different Kubernetes resource types (ConfigMaps, Secrets, Pods, ServiceAccounts)
5. **CLI Command Integration**: Tests the actual kgrep CLI commands against real resources
6. **Multi-Resource Scenarios**: Tests complex scenarios with multiple resource types

## Test Structure

The integration tests create a dedicated test namespace (`kgrep-integration-test`) and:

1. Set up test resources (ConfigMaps, Secrets, Pods, ServiceAccounts)
2. Run kgrep commands against these resources
3. Verify the expected output
4. Clean up resources after each test

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
3. **Non-Destructive**: Tests only create test resources and don't modify existing ones
4. **Opt-In**: Tests only run when explicitly enabled via environment variable

## CI/CD Integration

The integration tests run automatically in CI/CD via GitHub Actions:

- **Trigger**: On pushes to `main` branch and pull requests
- **Environment**: Ubuntu with kind-created Kubernetes cluster
- **Isolation**: Each CI run gets a fresh kind cluster

## Test Development

When adding new integration tests:

1. Follow the existing pattern of setup/test/cleanup
2. Use the dedicated test namespace
3. Include proper cleanup in `defer` statements
4. Add meaningful assertions and logging
5. Test both success and failure scenarios where applicable

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
```

### Kind Cluster Issues

If kind cluster creation fails:
```bash
kind delete cluster --name kgrep-integration-test
docker system prune -f
```
