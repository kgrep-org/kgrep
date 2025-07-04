package integration

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	testNamespacePrefix = "kgrep-integration-test"
	timeout             = 30 * time.Second
)

func TestMain(m *testing.M) {
	if os.Getenv("KGREP_INTEGRATION_TESTS") != "true" {
		fmt.Println("Skipping integration tests. Set KGREP_INTEGRATION_TESTS=true to run them.")
		os.Exit(0)
	}

	code := m.Run()
	os.Exit(code)
}

func getTestNamespace(t *testing.T) string {
	return fmt.Sprintf("%s-%d", testNamespacePrefix, time.Now().UnixNano())
}

func setupKubernetesClient(t *testing.T) kubernetes.Interface {
	kubeconfig := os.Getenv("KUBECONFIG")
	if kubeconfig == "" {
		homeDir, err := os.UserHomeDir()
		require.NoError(t, err)
		kubeconfig = fmt.Sprintf("%s/.kube/config", homeDir)
	}

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	require.NoError(t, err)

	clientset, err := kubernetes.NewForConfig(config)
	require.NoError(t, err)

	return clientset
}

func createTestNamespace(t *testing.T, clientset kubernetes.Interface, namespaceName string) {
	namespace := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: namespaceName,
		},
	}

	_, err := clientset.CoreV1().Namespaces().Create(context.Background(), namespace, metav1.CreateOptions{})
	if err != nil && !strings.Contains(err.Error(), "already exists") {
		require.NoError(t, err)
	}
}

func cleanupTestNamespace(t *testing.T, clientset kubernetes.Interface, namespaceName string) {
	err := clientset.CoreV1().Namespaces().Delete(context.Background(), namespaceName, metav1.DeleteOptions{})
	if err != nil && !strings.Contains(err.Error(), "not found") {
		t.Logf("Warning: Failed to cleanup test namespace: %v", err)
	}
}

func createTestConfigMap(t *testing.T, clientset kubernetes.Interface, namespaceName string) {
	configMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-config",
			Namespace: namespaceName,
		},
		Data: map[string]string{
			"app.properties": "database.url=postgresql://localhost:5432/mydb\napp.name=my-test-app\napp.version=1.0.0",
			"config.yaml":    "server:\n  host: localhost\n  port: 8080\nfeatures:\n  - authentication\n  - logging",
		},
	}

	_, err := clientset.CoreV1().ConfigMaps(namespaceName).Create(context.Background(), configMap, metav1.CreateOptions{})
	require.NoError(t, err)
}

func createTestSecret(t *testing.T, clientset kubernetes.Interface, namespaceName string) {
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-secret",
			Namespace: namespaceName,
		},
		Data: map[string][]byte{
			"username": []byte("testuser"),
			"password": []byte("secretpassword123"),
			"api-key":  []byte("api-key-12345"),
		},
	}

	_, err := clientset.CoreV1().Secrets(namespaceName).Create(context.Background(), secret, metav1.CreateOptions{})
	require.NoError(t, err)
}

func createTestPod(t *testing.T, clientset kubernetes.Interface, namespaceName string) {
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-pod",
			Namespace: namespaceName,
			Labels: map[string]string{
				"app":        "test-app",
				"component":  "backend",
				"version":    "v1.0.0",
				"managed-by": "kgrep-integration-test",
			},
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  "test-container",
					Image: "busybox:latest",
					Command: []string{
						"sh",
						"-c",
						"echo 'Starting test application'; echo 'Application initialized successfully'; while true; do echo 'Application running - $(date)'; sleep 30; done",
					},
				},
			},
			RestartPolicy: corev1.RestartPolicyAlways,
		},
	}

	_, err := clientset.CoreV1().Pods(namespaceName).Create(context.Background(), pod, metav1.CreateOptions{})
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			t.Fatal("Timeout waiting for pod to be ready")
		default:
			pod, err := clientset.CoreV1().Pods(namespaceName).Get(context.Background(), "test-pod", metav1.GetOptions{})
			if err != nil {
				time.Sleep(2 * time.Second)
				continue
			}
			if pod.Status.Phase == corev1.PodRunning {
				return
			}
			time.Sleep(2 * time.Second)
		}
	}
}

func createTestServiceAccount(t *testing.T, clientset kubernetes.Interface, namespaceName string) {
	sa := &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-service-account",
			Namespace: namespaceName,
			Annotations: map[string]string{
				"description": "Test service account for kgrep integration tests",
				"owner":       "kgrep-team",
			},
		},
	}

	_, err := clientset.CoreV1().ServiceAccounts(namespaceName).Create(context.Background(), sa, metav1.CreateOptions{})
	require.NoError(t, err)
}

func waitForLogs(clientset kubernetes.Interface, namespace, podName, expectedContent string, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("timeout waiting for logs containing '%s'", expectedContent)
		default:
			req := clientset.CoreV1().Pods(namespace).GetLogs(podName, &corev1.PodLogOptions{})
			logs, err := req.Do(context.Background()).Raw()
			if err == nil && strings.Contains(string(logs), expectedContent) {
				return nil
			}
			time.Sleep(2 * time.Second)
		}
	}
}

func runKgrepCommand(t *testing.T, args ...string) (string, error) {
	if _, err := os.Stat("../../kgrep"); os.IsNotExist(err) {
		buildCmd := exec.Command("go", "build", "-o", "kgrep", "main.go")
		buildCmd.Dir = "../../"
		err := buildCmd.Run()
		require.NoError(t, err, "Failed to build kgrep binary")
	}

	cmd := exec.Command("../../kgrep", args...)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

func TestIntegration_ConfigMapSearch(t *testing.T) {
	clientset := setupKubernetesClient(t)
	testNamespace := getTestNamespace(t)
	createTestNamespace(t, clientset, testNamespace)
	defer cleanupTestNamespace(t, clientset, testNamespace)

	createTestConfigMap(t, clientset, testNamespace)

	output, err := runKgrepCommand(t, "configmaps", "-n", testNamespace, "-p", "database")
	require.NoError(t, err, "kgrep command failed: %s", output)

	assert.Contains(t, output, "test-config")
	assert.Contains(t, output, "database.url=postgresql://localhost:5432/mydb")
	t.Logf("ConfigMap search output: %s", output)
}

func TestIntegration_SecretSearch(t *testing.T) {
	clientset := setupKubernetesClient(t)
	testNamespace := getTestNamespace(t)
	createTestNamespace(t, clientset, testNamespace)
	defer cleanupTestNamespace(t, clientset, testNamespace)

	createTestSecret(t, clientset, testNamespace)

	output, err := runKgrepCommand(t, "secrets", "-n", testNamespace, "-p", "username")
	require.NoError(t, err, "kgrep command failed: %s", output)

	assert.Contains(t, output, "test-secret")
	assert.Contains(t, output, "username")
	t.Logf("Secret search output: %s", output)
}

func TestIntegration_PodSearch(t *testing.T) {
	clientset := setupKubernetesClient(t)
	testNamespace := getTestNamespace(t)
	createTestNamespace(t, clientset, testNamespace)
	defer cleanupTestNamespace(t, clientset, testNamespace)

	createTestPod(t, clientset, testNamespace)

	output, err := runKgrepCommand(t, "pods", "-n", testNamespace, "-p", "test-app")
	require.NoError(t, err, "kgrep command failed: %s", output)

	assert.Contains(t, output, "test-pod")
	assert.Contains(t, output, "test-app")
	t.Logf("Pod search output: %s", output)
}

func TestIntegration_ServiceAccountSearch(t *testing.T) {
	clientset := setupKubernetesClient(t)
	testNamespace := getTestNamespace(t)
	createTestNamespace(t, clientset, testNamespace)
	defer cleanupTestNamespace(t, clientset, testNamespace)

	createTestServiceAccount(t, clientset, testNamespace)

	output, err := runKgrepCommand(t, "serviceaccounts", "-n", testNamespace, "-p", "kgrep-team")
	require.NoError(t, err, "kgrep command failed: %s", output)

	assert.Contains(t, output, "test-service-account")
	assert.Contains(t, output, "kgrep-team")
	t.Logf("ServiceAccount search output: %s", output)
}

func TestIntegration_LogSearch(t *testing.T) {
	clientset := setupKubernetesClient(t)
	testNamespace := getTestNamespace(t)
	createTestNamespace(t, clientset, testNamespace)
	defer cleanupTestNamespace(t, clientset, testNamespace)

	createTestPod(t, clientset, testNamespace)

	err := waitForLogs(clientset, testNamespace, "test-pod", "Application", timeout)
	require.NoError(t, err, "Failed to retrieve logs within timeout")

	output, err := runKgrepCommand(t, "logs", "-n", testNamespace, "-p", "Application")
	require.NoError(t, err, "kgrep command failed: %s", output)

	assert.Contains(t, output, "Application")
	t.Logf("Log search output: %s", output)
}

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

func setupDynamicClient(t *testing.T) dynamic.Interface {
	kubeconfig := os.Getenv("KUBECONFIG")
	if kubeconfig == "" {
		homeDir, err := os.UserHomeDir()
		require.NoError(t, err)
		kubeconfig = fmt.Sprintf("%s/.kube/config", homeDir)
	}

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	require.NoError(t, err)

	dynamicClient, err := dynamic.NewForConfig(config)
	require.NoError(t, err)

	return dynamicClient
}

func applyTestCRD(t *testing.T) {
	cmd := exec.Command("kubectl", "apply", "-f", "test-crd.yaml")
	cmd.Dir = "."
	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "Failed to apply test CRD: %s", output)
	t.Logf("Applied test CRD: %s", output)

	time.Sleep(5 * time.Second)
}

func deleteTestCRD(t *testing.T) {
	cmd := exec.Command("kubectl", "delete", "-f", "test-crd.yaml", "--ignore-not-found=true")
	cmd.Dir = "."
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Logf("Warning: Failed to delete test CRD: %s", output)
	} else {
		t.Logf("Deleted test CRD: %s", output)
	}
}

func createTestCustomResource(t *testing.T, dynamicClient dynamic.Interface, namespaceName string) {
	gvr := schema.GroupVersionResource{
		Group:    "test.kgrep.io",
		Version:  "v1",
		Resource: "testapplications",
	}

	testResource := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "test.kgrep.io/v1",
			"kind":       "TestApplication",
			"metadata": map[string]interface{}{
				"name":      "test-app",
				"namespace": namespaceName,
				"labels": map[string]interface{}{
					"app":        "test-application",
					"component":  "backend",
					"fraud":      "detection",
					"version":    "v1.0.0",
					"managed-by": "kgrep-integration-test",
				},
			},
			"spec": map[string]interface{}{
				"name":        "Test Application",
				"version":     "1.0.0",
				"description": "A test application for kgrep integration tests with fraud detection capabilities",
				"tags": []interface{}{
					"testing",
					"fraud-detection",
					"integration",
					"backend-service",
				},
			},
			"status": map[string]interface{}{
				"phase": "Ready",
				"conditions": []interface{}{
					map[string]interface{}{
						"type":    "Ready",
						"status":  "True",
						"reason":  "ApplicationReady",
						"message": "Application is ready for fraud detection",
					},
				},
			},
		},
	}

	_, err := dynamicClient.Resource(gvr).Namespace(namespaceName).Create(context.Background(), testResource, metav1.CreateOptions{})
	require.NoError(t, err)
	t.Logf("Created test custom resource in namespace %s", namespaceName)
}

func cleanupTestCustomResources(t *testing.T, dynamicClient dynamic.Interface, namespaceName string) {
	gvr := schema.GroupVersionResource{
		Group:    "test.kgrep.io",
		Version:  "v1",
		Resource: "testapplications",
	}

	err := dynamicClient.Resource(gvr).Namespace(namespaceName).DeleteCollection(context.Background(), metav1.DeleteOptions{}, metav1.ListOptions{})
	if err != nil {
		t.Logf("Warning: Failed to cleanup test custom resources: %v", err)
	}
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

	// This mimics the behavior from issue #125 where datasciencepipelinesapplications.opendatahub.io
	// should fallback to datasciencepipelinesapplications
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
