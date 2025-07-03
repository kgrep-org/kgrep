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
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	testNamespace = "kgrep-integration-test"
	timeout       = 30 * time.Second
)

func TestMain(m *testing.M) {
	if os.Getenv("KGREP_INTEGRATION_TESTS") != "true" {
		fmt.Println("Skipping integration tests. Set KGREP_INTEGRATION_TESTS=true to run them.")
		os.Exit(0)
	}

	code := m.Run()
	os.Exit(code)
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

func createTestNamespace(t *testing.T, clientset kubernetes.Interface) {
	namespace := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: testNamespace,
		},
	}

	_, err := clientset.CoreV1().Namespaces().Create(context.Background(), namespace, metav1.CreateOptions{})
	if err != nil && !strings.Contains(err.Error(), "already exists") {
		require.NoError(t, err)
	}
}

func cleanupTestNamespace(t *testing.T, clientset kubernetes.Interface) {
	err := clientset.CoreV1().Namespaces().Delete(context.Background(), testNamespace, metav1.DeleteOptions{})
	if err != nil && !strings.Contains(err.Error(), "not found") {
		t.Logf("Warning: Failed to cleanup test namespace: %v", err)
	}
}

func createTestConfigMap(t *testing.T, clientset kubernetes.Interface) {
	configMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-config",
			Namespace: testNamespace,
		},
		Data: map[string]string{
			"app.properties": "database.url=postgresql://localhost:5432/mydb\napp.name=my-test-app\napp.version=1.0.0",
			"config.yaml":    "server:\n  host: localhost\n  port: 8080\nfeatures:\n  - authentication\n  - logging",
		},
	}

	_, err := clientset.CoreV1().ConfigMaps(testNamespace).Create(context.Background(), configMap, metav1.CreateOptions{})
	require.NoError(t, err)
}

func createTestSecret(t *testing.T, clientset kubernetes.Interface) {
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-secret",
			Namespace: testNamespace,
		},
		Data: map[string][]byte{
			"username": []byte("testuser"),
			"password": []byte("secretpassword123"),
			"api-key":  []byte("api-key-12345"),
		},
	}

	_, err := clientset.CoreV1().Secrets(testNamespace).Create(context.Background(), secret, metav1.CreateOptions{})
	require.NoError(t, err)
}

func createTestPod(t *testing.T, clientset kubernetes.Interface) {
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-pod",
			Namespace: testNamespace,
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

	_, err := clientset.CoreV1().Pods(testNamespace).Create(context.Background(), pod, metav1.CreateOptions{})
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			t.Fatal("Timeout waiting for pod to be ready")
		default:
			pod, err := clientset.CoreV1().Pods(testNamespace).Get(context.Background(), "test-pod", metav1.GetOptions{})
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

func createTestServiceAccount(t *testing.T, clientset kubernetes.Interface) {
	sa := &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-service-account",
			Namespace: testNamespace,
			Annotations: map[string]string{
				"description": "Test service account for kgrep integration tests",
				"owner":       "kgrep-team",
			},
		},
	}

	_, err := clientset.CoreV1().ServiceAccounts(testNamespace).Create(context.Background(), sa, metav1.CreateOptions{})
	require.NoError(t, err)
}

func runKgrepCommand(t *testing.T, args ...string) (string, error) {
	if _, err := os.Stat("../../kgrep"); os.IsNotExist(err) {
		buildCmd := exec.Command("go", "build", "-o", "../../kgrep", "../../main.go")
		buildCmd.Dir = "../../"
		err := buildCmd.Run()
		require.NoError(t, err, "Failed to build kgrep binary")
	}

	cmd := exec.Command("../../kgrep", args...)
	cmd.Dir = "."
	output, err := cmd.CombinedOutput()
	return string(output), err
}

func TestIntegration_ConfigMapSearch(t *testing.T) {
	clientset := setupKubernetesClient(t)
	createTestNamespace(t, clientset)
	defer cleanupTestNamespace(t, clientset)

	createTestConfigMap(t, clientset)

	output, err := runKgrepCommand(t, "configmaps", "-n", testNamespace, "-p", "database")
	require.NoError(t, err, "kgrep command failed: %s", output)

	assert.Contains(t, output, "test-config")
	assert.Contains(t, output, "database.url=postgresql://localhost:5432/mydb")
	t.Logf("ConfigMap search output: %s", output)
}

func TestIntegration_SecretSearch(t *testing.T) {
	clientset := setupKubernetesClient(t)
	createTestNamespace(t, clientset)
	defer cleanupTestNamespace(t, clientset)

	createTestSecret(t, clientset)

	output, err := runKgrepCommand(t, "secrets", "-n", testNamespace, "-p", "testuser")
	require.NoError(t, err, "kgrep command failed: %s", output)

	assert.Contains(t, output, "test-secret")
	assert.Contains(t, output, "testuser")
	t.Logf("Secret search output: %s", output)
}

func TestIntegration_PodSearch(t *testing.T) {
	clientset := setupKubernetesClient(t)
	createTestNamespace(t, clientset)
	defer cleanupTestNamespace(t, clientset)

	createTestPod(t, clientset)

	output, err := runKgrepCommand(t, "pods", "-n", testNamespace, "-p", "test-app")
	require.NoError(t, err, "kgrep command failed: %s", output)

	assert.Contains(t, output, "test-pod")
	assert.Contains(t, output, "test-app")
	t.Logf("Pod search output: %s", output)
}

func TestIntegration_ServiceAccountSearch(t *testing.T) {
	clientset := setupKubernetesClient(t)
	createTestNamespace(t, clientset)
	defer cleanupTestNamespace(t, clientset)

	createTestServiceAccount(t, clientset)

	output, err := runKgrepCommand(t, "serviceaccounts", "-n", testNamespace, "-p", "kgrep-team")
	require.NoError(t, err, "kgrep command failed: %s", output)

	assert.Contains(t, output, "test-service-account")
	assert.Contains(t, output, "kgrep-team")
	t.Logf("ServiceAccount search output: %s", output)
}

func TestIntegration_LogSearch(t *testing.T) {
	clientset := setupKubernetesClient(t)
	createTestNamespace(t, clientset)
	defer cleanupTestNamespace(t, clientset)

	createTestPod(t, clientset)

	time.Sleep(10 * time.Second)

	output, err := runKgrepCommand(t, "logs", "-n", testNamespace, "-p", "Application")
	require.NoError(t, err, "kgrep command failed: %s", output)

	assert.Contains(t, output, "Application")
	t.Logf("Log search output: %s", output)
}

func TestIntegration_CrossNamespaceSearch(t *testing.T) {
	clientset := setupKubernetesClient(t)
	createTestNamespace(t, clientset)
	defer cleanupTestNamespace(t, clientset)

	createTestConfigMap(t, clientset)

	output, err := runKgrepCommand(t, "configmaps", "-p", "my-test-app")
	require.NoError(t, err, "kgrep command failed: %s", output)

	assert.Contains(t, output, "test-config")
	assert.Contains(t, output, "my-test-app")
	t.Logf("Cross-namespace search output: %s", output)
}

func TestIntegration_GenericResourceSearch(t *testing.T) {
	clientset := setupKubernetesClient(t)
	createTestNamespace(t, clientset)
	defer cleanupTestNamespace(t, clientset)

	createTestPod(t, clientset)

	output, err := runKgrepCommand(t, "resources", "--kind", "Pod", "--pattern", "backend", "--namespace", testNamespace)
	require.NoError(t, err, "kgrep command failed: %s", output)

	assert.Contains(t, output, "test-pod")
	assert.Contains(t, output, "backend")
	t.Logf("Generic resource search output: %s", output)
}

func TestIntegration_MultipleResourceTypesSetup(t *testing.T) {
	clientset := setupKubernetesClient(t)
	createTestNamespace(t, clientset)
	defer cleanupTestNamespace(t, clientset)

	createTestConfigMap(t, clientset)
	createTestSecret(t, clientset)
	createTestPod(t, clientset)
	createTestServiceAccount(t, clientset)
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
