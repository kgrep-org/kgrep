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
	clientset := fake.NewSimpleClientset()
	searcher := &Searcher{
		clientset:    clientset,
		resourceType: "pods",
	}

	occurrences := searcher.SearchWithoutNamespace("test")
	assert.Empty(t, occurrences)
}

func TestResourceSearcher_SearchWithNamespace(t *testing.T) {
	clientset := fake.NewSimpleClientset()
	searcher := &Searcher{
		clientset:    clientset,
		resourceType: "pods",
	}

	occurrences := searcher.Search("default", "test")
	assert.Empty(t, occurrences)
}

func TestResourceSearcher_GetPodNames(t *testing.T) {
	pod1 := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{Name: "pod1", Namespace: "default"},
	}
	pod2 := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{Name: "pod2", Namespace: "default"},
	}

	clientset := fake.NewSimpleClientset(pod1, pod2)
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

	clientset := fake.NewSimpleClientset(cm1, cm2)
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

	clientset := fake.NewSimpleClientset(secret1, secret2)
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

	clientset := fake.NewSimpleClientset(sa1, sa2)
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
	clientset := fake.NewSimpleClientset()
	searcher := &Searcher{clientset: clientset}

	_, err := searcher.getSecretYAML("default", "nonexistent")
	assert.Error(t, err)
}

func TestResourceSearcher_GetServiceAccountYAML_Error(t *testing.T) {
	clientset := fake.NewSimpleClientset()
	searcher := &Searcher{clientset: clientset}

	_, err := searcher.getServiceAccountYAML("default", "nonexistent")
	assert.Error(t, err)
}

func TestResourceSearcher_GetGenericResourceNames_Error(t *testing.T) {
	searcher := &Searcher{resourceType: "unknown"}

	_, err := searcher.getGenericResourceNames("default")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "dynamic client not available")
}

func TestResourceSearcher_GetSecretYAML_Success(t *testing.T) {
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{Name: "test-secret", Namespace: "default"},
		Data: map[string][]byte{
			"key1": []byte("value1"),
		},
	}

	clientset := fake.NewSimpleClientset(secret)
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

	clientset := fake.NewSimpleClientset(sa)
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

func TestNewAutoDiscoveryResourceSearcher(t *testing.T) {
	searcher := NewAutoDiscoveryResourceSearcher("Pod")
	assert.NotNil(t, searcher)
	assert.Equal(t, "Pod", searcher.kind)
	assert.Empty(t, searcher.apiVersion)
}

func TestDiscoverAPIVersionAndKind_NoClientset(t *testing.T) {
	searcher := &Searcher{kind: "Pod"}
	_, _, _, err := searcher.discoverAPIVersionAndKind()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "kubernetes clientset not available")
}

func TestDiscoverAPIVersionAndKind_WithFakeClientset(t *testing.T) {
	clientset := fake.NewSimpleClientset()
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
