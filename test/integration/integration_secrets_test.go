package integration

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

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

func TestIntegration_AllNamespaces_Secrets(t *testing.T) {
	clientset := setupKubernetesClient(t)
	namespace1 := getTestNamespace(t)
	namespace2 := getTestNamespace(t)

	createTestNamespace(t, clientset, namespace1)
	createTestNamespace(t, clientset, namespace2)
	defer func() {
		cleanupTestNamespace(t, clientset, namespace1)
		cleanupTestNamespace(t, clientset, namespace2)
	}()

	secret1 := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-secret-1",
			Namespace: namespace1,
			Annotations: map[string]string{
				"description": "secret-pattern-test",
				"purpose":     "integration-testing",
			},
		},
		Data: map[string][]byte{
			"password": []byte("mypassword123"),
		},
	}

	secret2 := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-secret-2",
			Namespace: namespace2,
			Annotations: map[string]string{
				"description": "secret-pattern-test",
				"purpose":     "integration-testing",
			},
		},
		Data: map[string][]byte{
			"password": []byte("mypassword456"),
		},
	}

	_, err := clientset.CoreV1().Secrets(namespace1).Create(context.Background(), secret1, metav1.CreateOptions{})
	require.NoError(t, err)

	_, err = clientset.CoreV1().Secrets(namespace2).Create(context.Background(), secret2, metav1.CreateOptions{})
	require.NoError(t, err)

	output, err := runKgrepCommand(t, "secrets", "--all-namespaces", "-p", "secret-pattern-test")
	require.NoError(t, err, "kgrep command failed: %s", output)

	assert.Contains(t, output, "test-secret-1")
	assert.Contains(t, output, "test-secret-2")
	assert.Contains(t, output, namespace1)
	assert.Contains(t, output, namespace2)
	assert.Contains(t, output, "secret-pattern-test")
	t.Logf("All namespaces Secret search output: %s", output)
}
