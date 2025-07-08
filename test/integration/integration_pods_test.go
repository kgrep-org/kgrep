package integration

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

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

func TestIntegration_AllNamespaces_Pods(t *testing.T) {
	clientset := setupKubernetesClient(t)
	namespace1 := getTestNamespace(t)
	namespace2 := getTestNamespace(t)

	createTestNamespace(t, clientset, namespace1)
	createTestNamespace(t, clientset, namespace2)
	defer func() {
		cleanupTestNamespace(t, clientset, namespace1)
		cleanupTestNamespace(t, clientset, namespace2)
	}()

	pod1 := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-pod-1",
			Namespace: namespace1,
			Labels: map[string]string{
				"app":     "shared-app-label",
				"version": "v1.0.0",
			},
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:    "test-container",
					Image:   "busybox:latest",
					Command: []string{"sh", "-c", "sleep 3600"},
				},
			},
			RestartPolicy: corev1.RestartPolicyAlways,
		},
	}

	pod2 := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-pod-2",
			Namespace: namespace2,
			Labels: map[string]string{
				"app":     "shared-app-label",
				"version": "v2.0.0",
			},
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:    "test-container",
					Image:   "busybox:latest",
					Command: []string{"sh", "-c", "sleep 3600"},
				},
			},
			RestartPolicy: corev1.RestartPolicyAlways,
		},
	}

	_, err := clientset.CoreV1().Pods(namespace1).Create(context.Background(), pod1, metav1.CreateOptions{})
	require.NoError(t, err)

	_, err = clientset.CoreV1().Pods(namespace2).Create(context.Background(), pod2, metav1.CreateOptions{})
	require.NoError(t, err)

	output, err := runKgrepCommand(t, "pods", "--all-namespaces", "-p", "shared-app-label")
	require.NoError(t, err, "kgrep command failed: %s", output)

	assert.Contains(t, output, "test-pod-1")
	assert.Contains(t, output, "test-pod-2")
	assert.Contains(t, output, namespace1)
	assert.Contains(t, output, namespace2)
	assert.Contains(t, output, "shared-app-label")
	t.Logf("All namespaces Pod search output: %s", output)
}
