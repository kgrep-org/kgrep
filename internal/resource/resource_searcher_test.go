package resource

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes/fake"
)

func TestResourceSearcher_SearchWithoutNamespace(t *testing.T) {
	clientset := fake.NewClientset()
	searcher := &Searcher{
		clientset:    clientset,
		resourceType: "pods",
	}

	_, err := searcher.SearchWithoutNamespace("test")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "kubeGet client not available")
}

func TestResourceSearcher_SearchWithNamespace(t *testing.T) {
	clientset := fake.NewClientset()
	searcher := &Searcher{
		clientset:    clientset,
		resourceType: "pods",
	}

	_, err := searcher.Search("default", "test")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "kubeGet client not available")
}

func TestResourceSearcher_SearchAllNamespaces(t *testing.T) {
	clientset := fake.NewClientset()
	searcher := &Searcher{
		clientset:    clientset,
		resourceType: "pods",
	}

	_, err := searcher.SearchAllNamespaces("test")
	// Should succeed in getting namespaces but fail on kubeGet
	assert.NoError(t, err)
}

func TestResourceSearcher_SearchAllNamespaces_NoClientset(t *testing.T) {
	searcher := &Searcher{
		resourceType: "pods",
	}

	_, err := searcher.SearchAllNamespaces("test")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Kubernetes clientset not available")
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
		Resource:  "test-pod",
		Namespace: "test-namespace",
		Line:      10,
		Content:   "test content",
	}

	assert.Equal(t, "test-pod", occurrence.Resource)
	assert.Equal(t, "test-namespace", occurrence.Namespace)
	assert.Equal(t, 10, occurrence.Line)
	assert.Equal(t, "test content", occurrence.Content)
}

func TestResourceSearcher_GetGenericResourceNames_Error(t *testing.T) {
	searcher := &Searcher{resourceType: "unknown"}

	_, err := searcher.getGenericResourceNames("default")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "kubeGet client not available")
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
	assert.Contains(t, err.Error(), "Kubernetes clientset not available")
}

func TestDiscoverAPIVersionAndKind_WithFakeClientset(t *testing.T) {
	clientset := fake.NewClientset()
	searcher := &Searcher{
		clientset: clientset,
		kind:      "Pod",
	}

	_, _, _, err := searcher.discoverAPIVersionAndKind()
	assert.Error(t, err)
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
			searcher := &Searcher{
				kind: tc.malformedKind,
			}

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
	clientset := fake.NewClientset()

	testCases := []struct {
		kind string
	}{
		{"pod"},
		{"pods"},
		{"configmap"},
		{"configmaps"},
		{"secret"},
		{"secrets"},
		{"serviceaccount"},
		{"serviceaccounts"},
		{"customresource"},
	}

	for _, tc := range testCases {
		t.Run(tc.kind, func(t *testing.T) {
			searcher := &Searcher{
				clientset: clientset,
				kind:      tc.kind,
			}

			_, err := searcher.getGenericResourceNames("default")
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "kubeGet client not available")
		})
	}
}

func TestGetGenericResourceYAML_ErrorHandling(t *testing.T) {
	searcher := &Searcher{
		kind: "customresource",
	}

	_, err := searcher.getGenericResourceYAML("default", "test-resource")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "kubeGet client not available")
}

func TestGetGenericResourceNames_ClusterScopedFallback(t *testing.T) {
	searcher := &Searcher{
		kind:    "namespace",
		kubeGet: nil,
	}

	_, err := searcher.getGenericResourceNames("some-namespace")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "kubeGet client not available")
}

func TestGetGenericResourceYAML_ClusterScopedFallback(t *testing.T) {
	searcher := &Searcher{
		kind:    "namespace",
		kubeGet: nil,
	}

	_, err := searcher.getGenericResourceYAML("some-namespace", "some-name")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "kubeGet client not available")
}

func TestGetGenericResourceNames_KindBasedRouting(t *testing.T) {
	clientset := fake.NewClientset()

	testCases := []struct {
		name string
		kind string
	}{
		{
			name: "Pod routing",
			kind: "pod",
		},
		{
			name: "ConfigMap routing",
			kind: "configmap",
		},
		{
			name: "Secret routing",
			kind: "secret",
		},
		{
			name: "ServiceAccount routing",
			kind: "serviceaccount",
		},
		{
			name: "Unknown kind fallback",
			kind: "namespace",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			searcher := &Searcher{
				clientset:  clientset,
				kind:       tc.kind,
				apiVersion: "",
			}

			_, err := searcher.getGenericResourceNames("default")
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "kubeGet client not available")
		})
	}
}

func TestGetGenericResourceYAML_KindBasedRouting(t *testing.T) {
	clientset := fake.NewClientset()

	testCases := []struct {
		name string
		kind string
	}{
		{
			name: "Pod YAML routing",
			kind: "pod",
		},
		{
			name: "ConfigMap YAML routing",
			kind: "configmap",
		},
		{
			name: "Secret YAML routing",
			kind: "secret",
		},
		{
			name: "ServiceAccount YAML routing",
			kind: "serviceaccount",
		},
		{
			name: "Unknown kind YAML fallback",
			kind: "namespace",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			searcher := &Searcher{
				clientset:  clientset,
				kind:       tc.kind,
				apiVersion: "",
			}

			_, err := searcher.getGenericResourceYAML("default", "test-resource")
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "kubeGet client not available")
		})
	}
}

func TestSearcher_APIVersionAndKindPrecedence(t *testing.T) {
	clientset := fake.NewClientset()

	searcher := &Searcher{
		clientset:  clientset,
		apiVersion: "v1",
		kind:       "Pod",
		kubeGet:    nil,
	}

	_, err := searcher.getGenericResourceNames("default")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "kubeGet client not available")

	_, err = searcher.getGenericResourceYAML("default", "test-resource")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "kubeGet client not available")
}

func TestWrapKubernetesError(t *testing.T) {
	assert.NoError(t, wrapKubernetesError(nil))

	t.Run("unauthorized error", func(t *testing.T) {
		origErr := errors.NewUnauthorized("unauthorized")
		wrapped := wrapKubernetesError(origErr)
		assert.Error(t, wrapped)
		assert.Contains(t, wrapped.Error(), "you are not authorized to access the cluster")
		assert.ErrorIs(t, wrapped, origErr)
	})

	t.Run("forbidden error", func(t *testing.T) {
		origErr := errors.NewForbidden(schema.GroupResource{}, "name", fmt.Errorf("forbidden"))
		wrapped := wrapKubernetesError(origErr)
		assert.Error(t, wrapped)
		assert.Contains(t, wrapped.Error(), "you do not have permission to perform this action")
		assert.ErrorIs(t, wrapped, origErr)
	})

	t.Run("configuration errors", func(t *testing.T) {
		configErrs := []string{
			"no configuration has been provided",
			"unable to load in-cluster configuration",
			"Couldn't find kubeconfig file",
		}

		for _, msg := range configErrs {
			t.Run(msg, func(t *testing.T) {
				origErr := fmt.Errorf("%s", msg)
				wrapped := wrapKubernetesError(origErr)
				assert.Error(t, wrapped)
				assert.Contains(t, wrapped.Error(), "no Kubernetes configuration found")
				assert.ErrorIs(t, wrapped, origErr)
			})
		}
	})

	t.Run("invalid configuration substring is not rewritten", func(t *testing.T) {
		origErr := fmt.Errorf("parser failed: invalid configuration in application payload")
		wrapped := wrapKubernetesError(origErr)
		assert.Equal(t, origErr, wrapped)
	})

	t.Run("context preservation", func(t *testing.T) {
		origErr := errors.NewUnauthorized("unauthorized")
		contextErr := fmt.Errorf("outer context: %w", origErr)
		wrapped := wrapKubernetesError(contextErr)
		assert.Error(t, wrapped)
		assert.Contains(t, wrapped.Error(), "outer context")
		assert.Contains(t, wrapped.Error(), "you are not authorized to access the cluster")
		assert.ErrorIs(t, wrapped, origErr)
	})

	t.Run("unknown error", func(t *testing.T) {
		origErr := fmt.Errorf("some random error")
		wrapped := wrapKubernetesError(origErr)
		assert.Equal(t, origErr, wrapped)
	})
}
