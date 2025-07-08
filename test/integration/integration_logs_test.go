package integration

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
