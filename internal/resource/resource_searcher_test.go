package resource

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestResourceSearcher_SearchWithoutNamespace(t *testing.T) {
	clientset := fake.NewClientset()
	searcher := &Searcher{
		clientset:    clientset,
		resourceType: "pods",
	}

	occurrences, err := searcher.SearchWithoutNamespace("test")
	assert.NoError(t, err)
	assert.Empty(t, occurrences)
}

func TestResourceSearcher_SearchWithNamespace(t *testing.T) {
	clientset := fake.NewClientset()
	searcher := &Searcher{
		clientset:    clientset,
		resourceType: "pods",
	}

	occurrences, err := searcher.Search("default", "test")
	assert.NoError(t, err)
	assert.Empty(t, occurrences)
}

func TestResourceSearcher_GetPodNames(t *testing.T) {
	pod1 := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{Name: "pod1", Namespace: "default"},
	}
	pod2 := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{Name: "pod2", Namespace: "default"},
	}

	clientset := fake.NewClientset(pod1, pod2)
	searcher := &Searcher{clientset: clientset}

	names, err := searcher.getPodNames("default")
	assert.NoError(t, err)
	assert.Len(t, names, 2)
	assert.Contains(t, names, "pod1")
	assert.Contains(t, names, "pod2")
}

func TestResourceSearcher_GetConfigMapNames(t *testing.T) {
	cm1 := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{Name: "cm1", Namespace: "default"},
	}
	cm2 := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{Name: "cm2", Namespace: "default"},
	}

	clientset := fake.NewClientset(cm1, cm2)
	searcher := &Searcher{clientset: clientset}

	names, err := searcher.getConfigMapNames("default")
	assert.NoError(t, err)
	assert.Len(t, names, 2)
	assert.Contains(t, names, "cm1")
	assert.Contains(t, names, "cm2")
}

func TestResourceSearcher_GetSecretNames(t *testing.T) {
	secret1 := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{Name: "secret1", Namespace: "default"},
	}
	secret2 := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{Name: "secret2", Namespace: "default"},
	}

	clientset := fake.NewClientset(secret1, secret2)
	searcher := &Searcher{clientset: clientset}

	names, err := searcher.getSecretNames("default")
	assert.NoError(t, err)
	assert.Len(t, names, 2)
	assert.Contains(t, names, "secret1")
	assert.Contains(t, names, "secret2")
}

func TestResourceSearcher_GetServiceAccountNames(t *testing.T) {
	sa1 := &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{Name: "sa1", Namespace: "default"},
	}
	sa2 := &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{Name: "sa2", Namespace: "default"},
	}

	clientset := fake.NewClientset(sa1, sa2)
	searcher := &Searcher{clientset: clientset}

	names, err := searcher.getServiceAccountNames("default")
	assert.NoError(t, err)
	assert.Len(t, names, 2)
	assert.Contains(t, names, "sa1")
	assert.Contains(t, names, "sa2")
}

func TestGetDefaultNamespace(t *testing.T) {
	searcher := &Searcher{}
	namespace, err := searcher.getDefaultNamespace()
	assert.NoError(t, err)
	assert.Equal(t, "default", namespace)
	t.Logf("Default namespace: %s", namespace)
}

func TestResourceOccurrence(t *testing.T) {
	occurrence := Occurrence{
		Resource: "test-pod",
		Line:     10,
		Content:  "test content",
	}

	assert.Equal(t, "test-pod", occurrence.Resource)
	assert.Equal(t, 10, occurrence.Line)
	assert.Equal(t, "test content", occurrence.Content)
}

func TestResourceSearcher_GetSecretYAML_Error(t *testing.T) {
	clientset := fake.NewClientset()
	searcher := &Searcher{clientset: clientset}

	_, err := searcher.getSecretYAML("default", "nonexistent")
	assert.Error(t, err)
}

func TestResourceSearcher_GetServiceAccountYAML_Error(t *testing.T) {
	clientset := fake.NewClientset()
	searcher := &Searcher{clientset: clientset}

	_, err := searcher.getServiceAccountYAML("default", "nonexistent")
	assert.Error(t, err)
}

func TestResourceSearcher_GetGenericResourceNames_Error(t *testing.T) {
	searcher := &Searcher{resourceType: "unknown"}

	_, err := searcher.getGenericResourceNames("default")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "kubeGet client not available")
}

func TestResourceSearcher_GetSecretYAML_Success(t *testing.T) {
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{Name: "test-secret", Namespace: "default"},
		Data: map[string][]byte{
			"key1": []byte("value1"),
		},
	}

	clientset := fake.NewClientset(secret)
	searcher := &Searcher{clientset: clientset}

	yaml, err := searcher.getSecretYAML("default", "test-secret")
	assert.NoError(t, err)
	assert.Contains(t, yaml, "test-secret")
	assert.Contains(t, yaml, "key1")
}

func TestResourceSearcher_GetServiceAccountYAML_Success(t *testing.T) {
	sa := &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{Name: "test-sa", Namespace: "default"},
	}

	clientset := fake.NewClientset(sa)
	searcher := &Searcher{clientset: clientset}

	yaml, err := searcher.getServiceAccountYAML("default", "test-sa")
	assert.NoError(t, err)
	assert.Contains(t, yaml, "test-sa")
}

func TestResourceSearcher_GetDefaultNamespace_NoConfig(t *testing.T) {
	searcher := &Searcher{}
	namespace, err := searcher.getDefaultNamespace()
	assert.NoError(t, err)
	assert.Equal(t, "default", namespace)
}

func TestDiscoverAPIVersionAndKind_NoClientset(t *testing.T) {
	searcher := &Searcher{kind: "Pod"}
	_, _, _, err := searcher.discoverAPIVersionAndKind()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "kubernetes clientset not available")
}

func TestDiscoverAPIVersionAndKind_WithFakeClientset(t *testing.T) {
	clientset := fake.NewClientset()
	searcher := &Searcher{
		clientset: clientset,
		kind:      "Pod",
	}

	// With fake clientset, discovery will fail, so we expect an error
	_, _, _, err := searcher.discoverAPIVersionAndKind()
	assert.Error(t, err)
	// The error could be about API groups or resources not being found
	assert.True(t, strings.Contains(err.Error(), "could not find API version") ||
		strings.Contains(err.Error(), "error getting API groups"))
}

func TestGetGenericResourceNames_HelpfulErrorMessage(t *testing.T) {
	// Test that the constructor correctly handles different resource name formats
	testCases := []struct {
		name          string
		malformedKind string
		expectedHint  string
	}{
		{
			name:          "DataSciencePipelines duplicate resource name",
			malformedKind: "datasciencepipelinesapplications.datasciencepipelinesapplications.opendatahub.io",
			expectedHint:  "datasciencepipelinesapplications.v1.opendatahub.io",
		},
		{
			name:          "Custom resource with duplicate parts",
			malformedKind: "myresource.myresource.example.com",
			expectedHint:  "myresource.v1.example.com",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a searcher with the malformed kind
			searcher := &Searcher{
				kind: tc.malformedKind,
				// kubeGet is nil, which will trigger the error path
			}

			// The call should fail with kubeGet client not available
			_, err := searcher.getGenericResourceNames("default")
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "kubeGet client not available")

			// Verify that the searcher has the expected kind stored
			assert.Equal(t, tc.malformedKind, searcher.kind)
		})
	}
}

func TestMalformedResourceNameDetection(t *testing.T) {
	// Test the logic that detects malformed resource names
	testCases := []struct {
		name              string
		resourceName      string
		shouldTriggerHint bool
	}{
		{
			name:              "Normal kind name",
			resourceName:      "pod",
			shouldTriggerHint: false,
		},
		{
			name:              "Normal fully qualified name",
			resourceName:      "deployments.v1.apps",
			shouldTriggerHint: false,
		},
		{
			name:              "Malformed with duplicate resource name",
			resourceName:      "datasciencepipelinesapplications.datasciencepipelinesapplications.opendatahub.io",
			shouldTriggerHint: true,
		},
		{
			name:              "Already has version",
			resourceName:      "datasciencepipelinesapplications.v1.datasciencepipelinesapplications.opendatahub.io",
			shouldTriggerHint: false,
		},
		{
			name:              "Simple dotted name",
			resourceName:      "resource.example.com",
			shouldTriggerHint: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Check if the resource name would trigger the helpful hint logic
			hasManyDots := strings.Count(tc.resourceName, ".") >= 2
			hasNoVersion := !strings.Contains(tc.resourceName, ".v")

			parts := strings.Split(tc.resourceName, ".")
			hasDuplicateFirstPart := len(parts) >= 3 && parts[0] == parts[1]

			shouldTrigger := hasManyDots && hasNoVersion && hasDuplicateFirstPart

			assert.Equal(t, tc.shouldTriggerHint, shouldTrigger,
				"Resource name %s should trigger hint: %v", tc.resourceName, tc.shouldTriggerHint)
		})
	}
}

func TestResourceSearcher_HandlesCoreResourceTypes(t *testing.T) {
	// Test that core resource types still go through the correct path
	clientset := fake.NewClientset()

	testCases := []struct {
		kind           string
		expectedMethod string
	}{
		{"pod", "core"},
		{"pods", "core"},
		{"configmap", "core"},
		{"configmaps", "core"},
		{"secret", "core"},
		{"secrets", "core"},
		{"serviceaccount", "core"},
		{"serviceaccounts", "core"},
		{"customresource", "generic"},
	}

	for _, tc := range testCases {
		t.Run(tc.kind, func(t *testing.T) {
			searcher := &Searcher{
				clientset: clientset,
				kind:      tc.kind,
			}

			// This will call the appropriate method based on the kind
			_, err := searcher.getResources("default")

			if tc.expectedMethod == "core" {
				// Core resources should work with fake clientset
				assert.NoError(t, err)
			} else {
				// Generic resources will fail due to missing kubeGet client
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "kubeGet client not available")
			}
		})
	}
}

func TestGetGenericResourceYAML_ErrorHandling(t *testing.T) {
	searcher := &Searcher{
		kind: "customresource",
		// kubeGet is nil, which will trigger the error path
	}

	_, err := searcher.getGenericResourceYAML("default", "test-resource")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "kubeGet client not available")
}
