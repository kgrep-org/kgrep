package integration

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

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

func TestIntegration_AllNamespaces_ConfigMaps(t *testing.T) {
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
			Name:      "test-config-1",
			Namespace: namespace1,
		},
		Data: map[string]string{
			"config": "shared-pattern-in-config",
		},
	}

	configMap2 := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-config-2",
			Namespace: namespace2,
		},
		Data: map[string]string{
			"config": "shared-pattern-in-config",
		},
	}

	_, err := clientset.CoreV1().ConfigMaps(namespace1).Create(context.Background(), configMap1, metav1.CreateOptions{})
	require.NoError(t, err)

	_, err = clientset.CoreV1().ConfigMaps(namespace2).Create(context.Background(), configMap2, metav1.CreateOptions{})
	require.NoError(t, err)

	output, err := runKgrepCommand(t, "configmaps", "--all-namespaces", "-p", "shared-pattern-in-config")
	require.NoError(t, err, "kgrep command failed: %s", output)

	assert.Contains(t, output, "test-config-1")
	assert.Contains(t, output, "test-config-2")
	assert.Contains(t, output, namespace1)
	assert.Contains(t, output, namespace2)
	assert.Contains(t, output, "shared-pattern-in-config")
	t.Logf("All namespaces ConfigMap search output: %s", output)
}
