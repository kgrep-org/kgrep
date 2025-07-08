package integration

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

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

func TestIntegration_AllNamespaces_ServiceAccounts(t *testing.T) {
	clientset := setupKubernetesClient(t)
	namespace1 := getTestNamespace(t)
	namespace2 := getTestNamespace(t)

	createTestNamespace(t, clientset, namespace1)
	createTestNamespace(t, clientset, namespace2)
	defer func() {
		cleanupTestNamespace(t, clientset, namespace1)
		cleanupTestNamespace(t, clientset, namespace2)
	}()

	sa1 := &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-service-account-1",
			Namespace: namespace1,
			Annotations: map[string]string{
				"description": "shared-annotation-pattern",
				"team":        "integration-test",
			},
		},
	}

	sa2 := &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-service-account-2",
			Namespace: namespace2,
			Annotations: map[string]string{
				"description": "shared-annotation-pattern",
				"team":        "integration-test",
			},
		},
	}

	_, err := clientset.CoreV1().ServiceAccounts(namespace1).Create(context.Background(), sa1, metav1.CreateOptions{})
	require.NoError(t, err)

	_, err = clientset.CoreV1().ServiceAccounts(namespace2).Create(context.Background(), sa2, metav1.CreateOptions{})
	require.NoError(t, err)

	output, err := runKgrepCommand(t, "serviceaccounts", "--all-namespaces", "-p", "shared-annotation-pattern")
	require.NoError(t, err, "kgrep command failed: %s", output)

	assert.Contains(t, output, "test-service-account-1")
	assert.Contains(t, output, "test-service-account-2")
	assert.Contains(t, output, namespace1)
	assert.Contains(t, output, namespace2)
	assert.Contains(t, output, "shared-annotation-pattern")
	t.Logf("All namespaces ServiceAccount search output: %s", output)
}
