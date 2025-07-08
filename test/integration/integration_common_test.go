package integration

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

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
