package integration

import (
	"context"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestIntegration_GenericResourceSearch(t *testing.T) {
	clientset := setupKubernetesClient(t)
	testNamespace := getTestNamespace(t)
	createTestNamespace(t, clientset, testNamespace)
	defer cleanupTestNamespace(t, clientset, testNamespace)

	createTestPod(t, clientset, testNamespace)

	output, err := runKgrepCommand(t, "resources", "--kind", "Pod", "--pattern", "backend", "--namespace", testNamespace)
	require.NoError(t, err, "kgrep command failed: %s", output)

	assert.Contains(t, output, "test-pod")
	assert.Contains(t, output, "backend")
	t.Logf("Generic resource search output: %s", output)
}

func TestIntegration_CustomResourceAutoDiscovery(t *testing.T) {
	clientset := setupKubernetesClient(t)
	dynamicClient := setupDynamicClient(t)
	testNamespace := getTestNamespace(t)

	applyTestCRD(t)
	defer deleteTestCRD(t)

	createTestNamespace(t, clientset, testNamespace)
	defer cleanupTestNamespace(t, clientset, testNamespace)
	defer cleanupTestCustomResources(t, dynamicClient, testNamespace)

	createTestCustomResource(t, dynamicClient, testNamespace)

	output, err := runKgrepCommand(t, "resources", "--kind", "TestApplication", "--pattern", "fraud", "--namespace", testNamespace)
	require.NoError(t, err, "kgrep command failed: %s", output)

	assert.Contains(t, output, "test-app")
	assert.Contains(t, output, "fraud")
	t.Logf("Custom resource auto-discovery output: %s", output)
}

func TestIntegration_ResourceNameFormats(t *testing.T) {
	clientset := setupKubernetesClient(t)
	dynamicClient := setupDynamicClient(t)
	testNamespace := getTestNamespace(t)

	applyTestCRD(t)
	defer deleteTestCRD(t)

	createTestNamespace(t, clientset, testNamespace)
	defer cleanupTestNamespace(t, clientset, testNamespace)
	defer cleanupTestCustomResources(t, dynamicClient, testNamespace)

	createTestCustomResource(t, dynamicClient, testNamespace)

	testCases := []struct {
		name         string
		resourceName string
		description  string
	}{
		{
			name:         "plural name",
			resourceName: "testapplications",
			description:  "Test with plural resource name",
		},
		{
			name:         "short name",
			resourceName: "tapp",
			description:  "Test with short name (alias)",
		},
		{
			name:         "resource.group format",
			resourceName: "testapplications.test.kgrep.io",
			description:  "Test with resource.group format",
		},
		{
			name:         "Kind name",
			resourceName: "TestApplication",
			description:  "Test with Kind name",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			output, err := runKgrepCommand(t, "resources", "--kind", tc.resourceName, "--pattern", "fraud", "--namespace", testNamespace)
			require.NoError(t, err, "kgrep command failed for %s: %s", tc.resourceName, output)

			assert.Contains(t, output, "test-app")
			assert.Contains(t, output, "fraud")
			t.Logf("%s (%s) output: %s", tc.description, tc.resourceName, output)
		})
	}
}

func TestIntegration_ResourceGroupFallback(t *testing.T) {
	clientset := setupKubernetesClient(t)
	dynamicClient := setupDynamicClient(t)
	testNamespace := getTestNamespace(t)

	applyTestCRD(t)
	defer deleteTestCRD(t)

	createTestNamespace(t, clientset, testNamespace)
	defer cleanupTestNamespace(t, clientset, testNamespace)
	defer cleanupTestCustomResources(t, dynamicClient, testNamespace)

	createTestCustomResource(t, dynamicClient, testNamespace)

	output, err := runKgrepCommand(t, "resources", "--kind", "testapplications.test.kgrep.io", "--pattern", "fraud", "--namespace", testNamespace)
	require.NoError(t, err, "kgrep command failed for fallback test: %s", output)

	assert.Contains(t, output, "test-app")
	assert.Contains(t, output, "fraud")
	t.Logf("Resource.group fallback test output: %s", output)

	output2, err := runKgrepCommand(t, "resources", "--kind", "testapplications", "--pattern", "fraud", "--namespace", testNamespace)
	require.NoError(t, err, "kgrep command failed for simple plural test: %s", output2)

	assert.Contains(t, output2, "test-app")
	assert.Contains(t, output2, "fraud")

	assert.Contains(t, output, "test-app")
	assert.Contains(t, output2, "test-app")
	t.Logf("Simple plural name test output: %s", output2)
}

func TestIntegration_CustomResourceErrors(t *testing.T) {
	clientset := setupKubernetesClient(t)
	testNamespace := getTestNamespace(t)

	createTestNamespace(t, clientset, testNamespace)
	defer cleanupTestNamespace(t, clientset, testNamespace)

	output, err := runKgrepCommand(t, "resources", "--kind", "NonExistentResource", "--pattern", "anything", "--namespace", testNamespace)
	require.Error(t, err, "Expected error for non-existent resource type")

	assert.Contains(t, output, "error")
	t.Logf("Non-existent resource error output: %s", output)

	output2, err := runKgrepCommand(t, "resources", "--kind", "invalid..resource...name", "--pattern", "anything", "--namespace", testNamespace)
	require.Error(t, err, "Expected error for malformed resource name")

	assert.Contains(t, output2, "error")
	t.Logf("Malformed resource name error output: %s", output2)
}

func TestIntegration_KubectlCompatibility(t *testing.T) {
	clientset := setupKubernetesClient(t)
	dynamicClient := setupDynamicClient(t)
	testNamespace := getTestNamespace(t)

	applyTestCRD(t)
	defer deleteTestCRD(t)

	createTestNamespace(t, clientset, testNamespace)
	defer cleanupTestNamespace(t, clientset, testNamespace)
	defer cleanupTestCustomResources(t, dynamicClient, testNamespace)

	createTestCustomResource(t, dynamicClient, testNamespace)

	kubectlCmd := exec.Command("kubectl", "get", "testapplications", "-n", testNamespace, "-o", "name")
	kubectlOutput, err := kubectlCmd.CombinedOutput()
	require.NoError(t, err, "kubectl command failed: %s", kubectlOutput)

	kubectlResources := strings.TrimSpace(string(kubectlOutput))
	assert.Contains(t, kubectlResources, "test-app")
	t.Logf("kubectl found resources: %s", kubectlResources)

	kgrepOutput, err := runKgrepCommand(t, "resources", "--kind", "testapplications", "--pattern", "fraud", "--namespace", testNamespace)
	require.NoError(t, err, "kgrep command failed: %s", kgrepOutput)

	assert.Contains(t, kgrepOutput, "test-app")
	t.Logf("kgrep found resources: %s", kgrepOutput)

	assert.Contains(t, kubectlResources, "test-app")
	assert.Contains(t, kgrepOutput, "test-app")
}

func TestIntegration_CustomResourceSearchContent(t *testing.T) {
	clientset := setupKubernetesClient(t)
	dynamicClient := setupDynamicClient(t)
	testNamespace := getTestNamespace(t)

	applyTestCRD(t)
	defer deleteTestCRD(t)

	createTestNamespace(t, clientset, testNamespace)
	defer cleanupTestNamespace(t, clientset, testNamespace)
	defer cleanupTestCustomResources(t, dynamicClient, testNamespace)

	createTestCustomResource(t, dynamicClient, testNamespace)

	testCases := []struct {
		name     string
		pattern  string
		expected string
	}{
		{
			name:     "search in labels",
			pattern:  "fraud",
			expected: "test-app",
		},
		{
			name:     "search in spec",
			pattern:  "Test Application",
			expected: "test-app",
		},
		{
			name:     "search in tags",
			pattern:  "fraud-detection",
			expected: "test-app",
		},
		{
			name:     "search in status",
			pattern:  "ApplicationReady",
			expected: "test-app",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			output, err := runKgrepCommand(t, "resources", "--kind", "TestApplication", "--pattern", tc.pattern, "--namespace", testNamespace)
			require.NoError(t, err, "kgrep command failed for %s: %s", tc.name, output)

			assert.Contains(t, output, tc.expected)
			assert.Contains(t, output, tc.pattern)
			t.Logf("%s output: %s", tc.name, output)
		})
	}
}

func TestIntegration_ClusterScopedResourceSearch(t *testing.T) {
	clientset := setupKubernetesClient(t)
	namespaceName := getTestNamespace(t)

	createTestNamespace(t, clientset, namespaceName)
	defer cleanupTestNamespace(t, clientset, namespaceName)

	output, err := runKgrepCommand(t, "resources", "-k", "namespace", "-p", namespaceName)
	assert.NoError(t, err)
	assert.Contains(t, output, namespaceName, "Should find the test namespace")

	output, err = runKgrepCommand(t, "resources", "-k", "namespaces", "-p", namespaceName)
	assert.NoError(t, err)
	assert.Contains(t, output, namespaceName, "Should find the test namespace using plural form")
}

func TestIntegration_NamespacedResourceIsolation(t *testing.T) {
	clientset := setupKubernetesClient(t)
	namespace1 := getTestNamespace(t)
	namespace2 := getTestNamespace(t)

	createTestNamespace(t, clientset, namespace1)
	createTestNamespace(t, clientset, namespace2)
	defer func() {
		cleanupTestNamespace(t, clientset, namespace1)
		cleanupTestNamespace(t, clientset, namespace2)
	}()

	configMap1 := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-config",
			Namespace: namespace1,
		},
		Data: map[string]string{
			"config": "unique-content-namespace1",
		},
	}

	configMap2 := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-config",
			Namespace: namespace2,
		},
		Data: map[string]string{
			"config": "unique-content-namespace2",
		},
	}

	_, err := clientset.CoreV1().ConfigMaps(namespace1).Create(context.Background(), configMap1, metav1.CreateOptions{})
	require.NoError(t, err)

	_, err = clientset.CoreV1().ConfigMaps(namespace2).Create(context.Background(), configMap2, metav1.CreateOptions{})
	require.NoError(t, err)

	output, err := runKgrepCommand(t, "resources", "-k", "configmap", "-n", namespace1, "-p", "unique-content-namespace1")
	assert.NoError(t, err)
	assert.Contains(t, output, "unique-content-namespace1", "Should find content from namespace1")

	output, err = runKgrepCommand(t, "resources", "-k", "configmap", "-n", namespace1, "-p", "unique-content-namespace2")
	assert.NoError(t, err)
	assert.Contains(t, output, "No occurrences", "Should not find content from namespace2 when searching in namespace1")

	output, err = runKgrepCommand(t, "resources", "-k", "configmap", "-n", namespace2, "-p", "unique-content-namespace2")
	assert.NoError(t, err)
	assert.Contains(t, output, "unique-content-namespace2", "Should find content from namespace2")

	output, err = runKgrepCommand(t, "resources", "-k", "configmap", "-n", namespace2, "-p", "unique-content-namespace1")
	assert.NoError(t, err)
	assert.Contains(t, output, "No occurrences", "Should not find content from namespace1 when searching in namespace2")
}

func TestIntegration_NonExistentNamespaceIsolation(t *testing.T) {
	clientset := setupKubernetesClient(t)
	realNamespace := getTestNamespace(t)
	nonExistentNamespace := "non-existent-namespace-12345"

	createTestNamespace(t, clientset, realNamespace)
	defer cleanupTestNamespace(t, clientset, realNamespace)

	configMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-config",
			Namespace: realNamespace,
		},
		Data: map[string]string{
			"config": "content-in-real-namespace",
		},
	}

	_, err := clientset.CoreV1().ConfigMaps(realNamespace).Create(context.Background(), configMap, metav1.CreateOptions{})
	require.NoError(t, err)

	output, err := runKgrepCommand(t, "resources", "-k", "configmap", "-n", nonExistentNamespace, "-p", "content-in-real-namespace")
	assert.NoError(t, err)
	assert.Contains(t, output, "No occurrences", "Should not find content from real namespace when searching in non-existent namespace")

	output, err = runKgrepCommand(t, "resources", "-k", "configmap", "-n", realNamespace, "-p", "content-in-real-namespace")
	assert.NoError(t, err)
	assert.Contains(t, output, "content-in-real-namespace", "Should find content when searching in correct namespace")
}

func TestIntegration_ClusterScopedVsNamespacedBehavior(t *testing.T) {
	clientset := setupKubernetesClient(t)
	testNamespace := getTestNamespace(t)

	createTestNamespace(t, clientset, testNamespace)
	defer cleanupTestNamespace(t, clientset, testNamespace)

	output, err := runKgrepCommand(t, "resources", "-k", "namespace", "-n", testNamespace, "-p", "default")
	assert.NoError(t, err)
	assert.Contains(t, output, "default", "Should find default namespace even when specifying a different namespace for cluster-scoped resources")

	createTestPod(t, clientset, testNamespace)

	output, err = runKgrepCommand(t, "resources", "-k", "pod", "-n", testNamespace, "-p", "test-pod")
	assert.NoError(t, err)
	assert.Contains(t, output, "test-pod", "Should find pod in specified namespace")

	output, err = runKgrepCommand(t, "resources", "-k", "pod", "-n", "default", "-p", "test-pod")
	assert.NoError(t, err)
	assert.Contains(t, output, "No occurrences", "Should not find pod from test namespace when searching in default namespace")
}

func TestIntegration_AllNamespaces_Resources(t *testing.T) {
	clientset := setupKubernetesClient(t)
	namespace1 := getTestNamespace(t)
	namespace2 := getTestNamespace(t)

	createTestNamespace(t, clientset, namespace1)
	createTestNamespace(t, clientset, namespace2)
	defer func() {
		cleanupTestNamespace(t, clientset, namespace1)
		cleanupTestNamespace(t, clientset, namespace2)
	}()

	configMap1 := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-config-generic-1",
			Namespace: namespace1,
		},
		Data: map[string]string{
			"config": "generic-resource-pattern",
		},
	}

	configMap2 := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-config-generic-2",
			Namespace: namespace2,
		},
		Data: map[string]string{
			"config": "generic-resource-pattern",
		},
	}

	_, err := clientset.CoreV1().ConfigMaps(namespace1).Create(context.Background(), configMap1, metav1.CreateOptions{})
	require.NoError(t, err)

	_, err = clientset.CoreV1().ConfigMaps(namespace2).Create(context.Background(), configMap2, metav1.CreateOptions{})
	require.NoError(t, err)

	output, err := runKgrepCommand(t, "resources", "--kind", "ConfigMap", "--all-namespaces", "-p", "generic-resource-pattern")
	require.NoError(t, err, "kgrep command failed: %s", output)

	assert.Contains(t, output, "test-config-generic-1")
	assert.Contains(t, output, "test-config-generic-2")
	assert.Contains(t, output, namespace1)
	assert.Contains(t, output, namespace2)
	assert.Contains(t, output, "generic-resource-pattern")
	t.Logf("All namespaces generic resource search output: %s", output)
}
