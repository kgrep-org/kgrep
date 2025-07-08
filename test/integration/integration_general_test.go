package integration

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
)

func TestIntegration_CurrentNamespaceSearch(t *testing.T) {
	clientset := setupKubernetesClient(t)

	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	configOverrides := &clientcmd.ConfigOverrides{}
	kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)

	currentNamespace, _, err := kubeConfig.Namespace()
	if err != nil || currentNamespace == "" {
		currentNamespace = "default"
	}

	createTestConfigMap(t, clientset, currentNamespace)
	defer func() {
		err := clientset.CoreV1().ConfigMaps(currentNamespace).Delete(context.Background(), "test-config", metav1.DeleteOptions{})
		if err != nil {
			t.Logf("Warning: Failed to cleanup ConfigMap in namespace %s: %v", currentNamespace, err)
		}
	}()

	output, err := runKgrepCommand(t, "configmaps", "-p", "my-test-app")
	require.NoError(t, err, "kgrep command failed: %s", output)

	assert.Contains(t, output, "test-config")
	assert.Contains(t, output, "my-test-app")
	t.Logf("Current namespace (%s) search output: %s", currentNamespace, output)
}

func TestIntegration_MultipleResourceTypesSetup(t *testing.T) {
	clientset := setupKubernetesClient(t)
	testNamespace := getTestNamespace(t)
	createTestNamespace(t, clientset, testNamespace)
	defer cleanupTestNamespace(t, clientset, testNamespace)

	createTestConfigMap(t, clientset, testNamespace)
	createTestSecret(t, clientset, testNamespace)
	createTestPod(t, clientset, testNamespace)
	createTestServiceAccount(t, clientset, testNamespace)

	tests := []struct {
		name    string
		cmd     []string
		expects string
	}{
		{
			name:    "ConfigMap search",
			cmd:     []string{"configmaps", "-n", testNamespace, "-p", "app.name"},
			expects: "test-config",
		},
		{
			name:    "Secret search",
			cmd:     []string{"secrets", "-n", testNamespace, "-p", "api-key"},
			expects: "test-secret",
		},
		{
			name:    "Pod search",
			cmd:     []string{"pods", "-n", testNamespace, "-p", "managed-by"},
			expects: "test-pod",
		},
		{
			name:    "ServiceAccount search",
			cmd:     []string{"serviceaccounts", "-n", testNamespace, "-p", "description"},
			expects: "test-service-account",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := runKgrepCommand(t, tt.cmd...)
			require.NoError(t, err, "kgrep command failed: %s", output)
			assert.Contains(t, output, tt.expects)
			t.Logf("%s output: %s", tt.name, output)
		})
	}
}

func TestIntegration_AllNamespaces_FlagValidation(t *testing.T) {
	tests := []struct {
		name     string
		cmd      []string
		expected string
	}{
		{
			name:     "Pods mutual exclusion",
			cmd:      []string{"pods", "--namespace", "test-ns", "--all-namespaces", "-p", "test"},
			expected: "--all-namespaces and --namespace cannot be used together",
		},
		{
			name:     "ConfigMaps mutual exclusion",
			cmd:      []string{"configmaps", "--namespace", "test-ns", "--all-namespaces", "-p", "test"},
			expected: "--all-namespaces and --namespace cannot be used together",
		},
		{
			name:     "Secrets mutual exclusion",
			cmd:      []string{"secrets", "--namespace", "test-ns", "--all-namespaces", "-p", "test"},
			expected: "--all-namespaces and --namespace cannot be used together",
		},
		{
			name:     "ServiceAccounts mutual exclusion",
			cmd:      []string{"serviceaccounts", "--namespace", "test-ns", "--all-namespaces", "-p", "test"},
			expected: "--all-namespaces and --namespace cannot be used together",
		},
		{
			name:     "Resources mutual exclusion",
			cmd:      []string{"resources", "--kind", "Pod", "--namespace", "test-ns", "--all-namespaces", "-p", "test"},
			expected: "--all-namespaces and --namespace cannot be used together",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := runKgrepCommand(t, tt.cmd...)
			assert.Error(t, err, "Expected error for mutual exclusion")
			assert.Contains(t, output, tt.expected)
			t.Logf("%s validation output: %s", tt.name, output)
		})
	}
}

func TestIntegration_AllNamespaces_NoResults(t *testing.T) {
	clientset := setupKubernetesClient(t)
	namespace1 := getTestNamespace(t)
	namespace2 := getTestNamespace(t)

	createTestNamespace(t, clientset, namespace1)
	createTestNamespace(t, clientset, namespace2)
	defer func() {
		cleanupTestNamespace(t, clientset, namespace1)
		cleanupTestNamespace(t, clientset, namespace2)
	}()

	createTestConfigMap(t, clientset, namespace1)
	createTestSecret(t, clientset, namespace2)

	output, err := runKgrepCommand(t, "configmaps", "--all-namespaces", "-p", "non-existent-pattern")
	require.NoError(t, err, "kgrep command failed: %s", output)
	assert.Contains(t, output, "No occurrences of 'non-existent-pattern' found")
	t.Logf("All namespaces no results output: %s", output)
}

func TestIntegration_AllNamespaces_NamespaceDisplayFormat(t *testing.T) {
	clientset := setupKubernetesClient(t)
	testNamespace := getTestNamespace(t)

	createTestNamespace(t, clientset, testNamespace)
	defer cleanupTestNamespace(t, clientset, testNamespace)

	createTestConfigMap(t, clientset, testNamespace)

	output, err := runKgrepCommand(t, "configmaps", "--all-namespaces", "-p", "my-test-app")
	require.NoError(t, err, "kgrep command failed: %s", output)

	assert.Contains(t, output, fmt.Sprintf("%s/test-config", testNamespace))
	t.Logf("All namespaces namespace display format output: %s", output)
}
