package resource

import (
	"context"
	"fmt"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"k8s.io/apimachinery/pkg/runtime/serializer/json"

	"github.com/hbelmiro/go-kube-get/pkg/gokubeget"
)

// Searcher is responsible for searching patterns in Kubernetes resources.
type Searcher struct {
	resourceType  string
	apiVersion    string
	kind          string
	resourceName  string // The plural resource name (e.g., "datasciencepipelinesapplications")
	clientset     kubernetes.Interface
	dynamicClient dynamic.Interface
	config        *rest.Config
	kubeGet       *gokubeget.KubeGet
}

// NewResourceSearcher creates a new ResourceSearcher for the specified resource type.
func NewResourceSearcher(resourceType string) (*Searcher, error) {
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	configOverrides := &clientcmd.ConfigOverrides{}
	config, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides).ClientConfig()
	if err != nil {
		return nil, fmt.Errorf("error creating Kubernetes config: %v", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("error creating Kubernetes clientset: %v", err)
	}

	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("error creating dynamic client: %v", err)
	}

	kubeGet, err := gokubeget.NewKubeGet(config)
	if err != nil {
		// kubeGet is optional for core resources, so we don't return an error here
		kubeGet = nil
	}

	return &Searcher{
		clientset:     clientset,
		dynamicClient: dynamicClient,
		config:        config,
		resourceType:  resourceType,
		kubeGet:       kubeGet,
	}, nil
}

// NewGenericResourceSearcher creates a new ResourceSearcher for generic resources with API version and kind.
func NewGenericResourceSearcher(apiVersion, kind string) (*Searcher, error) {
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	configOverrides := &clientcmd.ConfigOverrides{}
	config, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides).ClientConfig()
	if err != nil {
		return nil, fmt.Errorf("error creating Kubernetes config: %v", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("error creating Kubernetes clientset: %v", err)
	}

	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("error creating dynamic client: %v", err)
	}

	kubeGet, err := gokubeget.NewKubeGet(config)
	if err != nil {
		// kubeGet is optional for core resources, so we don't return an error here
		kubeGet = nil
	}

	return &Searcher{
		clientset:     clientset,
		dynamicClient: dynamicClient,
		config:        config,
		apiVersion:    apiVersion,
		kind:          kind,
		kubeGet:       kubeGet,
	}, nil
}

// NewAutoDiscoveryResourceSearcher creates a new ResourceSearcher that auto-discovers API version and kind.
func NewAutoDiscoveryResourceSearcher(kind string) (*Searcher, error) {
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	configOverrides := &clientcmd.ConfigOverrides{}
	config, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides).ClientConfig()
	if err != nil {
		return nil, fmt.Errorf("error creating Kubernetes config: %v", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("error creating Kubernetes clientset: %v", err)
	}

	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("error creating dynamic client: %v", err)
	}

	kubeGet, err := gokubeget.NewKubeGet(config)
	if err != nil {
		// kubeGet is optional for core resources, so we don't return an error here
		kubeGet = nil
	}

	return &Searcher{
		clientset:     clientset,
		dynamicClient: dynamicClient,
		config:        config,
		kind:          kind,
		kubeGet:       kubeGet,
	}, nil
}

// SearchWithoutNamespace searches for a pattern in resources in the default namespace.
func (s *Searcher) SearchWithoutNamespace(pattern string) ([]Occurrence, error) {
	namespace, err := s.getDefaultNamespace()
	if err != nil {
		return nil, fmt.Errorf("error getting default namespace: %v", err)
	}
	return s.Search(namespace, pattern)
}

// Search searches for a pattern in resources in a specific namespace.
func (s *Searcher) Search(namespace, pattern string) ([]Occurrence, error) {
	if s.clientset == nil {
		return nil, fmt.Errorf("Kubernetes clientset not available")
	}

	resources, err := s.getResources(namespace)
	if err != nil {
		return nil, fmt.Errorf("error getting resources: %v", err)
	}

	var occurrences []Occurrence
	for _, resource := range resources {
		resourceOccurrences := s.searchResource(namespace, resource, pattern)
		occurrences = append(occurrences, resourceOccurrences...)
	}

	return occurrences, nil
}

// searchResource searches for a pattern in a specific resource.
func (s *Searcher) searchResource(namespace, resource, pattern string) []Occurrence {
	yaml, err := s.getResourceYAML(namespace, resource)
	if err != nil {
		return []Occurrence{}
	}

	var occurrences []Occurrence
	lines := strings.Split(yaml, "\n")
	for i, line := range lines {
		if strings.Contains(strings.ToLower(line), strings.ToLower(pattern)) {
			occurrences = append(occurrences, Occurrence{
				Resource: resource,
				Line:     i + 1,
				Content:  line,
			})
		}
	}

	return occurrences
}

// getDefaultNamespace gets the default namespace from kubeconfig.
func (s *Searcher) getDefaultNamespace() (string, error) {
	if s.config == nil {
		return "default", nil
	}

	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	configOverrides := &clientcmd.ConfigOverrides{}
	kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)

	namespace, _, err := kubeConfig.Namespace()
	if err != nil {
		return "default", nil
	}

	if namespace == "" {
		namespace = "default"
	}
	return namespace, nil
}

// getResources gets the resources of a specific type in a namespace.
func (s *Searcher) getResources(namespace string) ([]string, error) {
	// If we have API version and kind, use generic resource handling
	if s.apiVersion != "" && s.kind != "" {
		return s.getGenericResourceNames(namespace)
	}

	// If we only have kind, check if it's a core resource first
	if s.kind != "" && s.apiVersion == "" {
		lowerKind := strings.ToLower(s.kind)
		switch lowerKind {
		case "pod", "pods":
			return s.getPodNames(namespace)
		case "configmap", "configmaps":
			return s.getConfigMapNames(namespace)
		case "secret", "secrets":
			return s.getSecretNames(namespace)
		case "serviceaccount", "serviceaccounts":
			return s.getServiceAccountNames(namespace)
		default:
			// For non-core resources, try auto-discovery
			return s.getGenericResourceNames(namespace)
		}
	}

	switch strings.ToLower(s.resourceType) {
	case "pods":
		return s.getPodNames(namespace)
	case "configmaps":
		return s.getConfigMapNames(namespace)
	case "secrets":
		return s.getSecretNames(namespace)
	case "serviceaccounts":
		return s.getServiceAccountNames(namespace)
	default:
		return s.getGenericResourceNames(namespace)
	}
}

// getPodNames gets pod names in a namespace.
func (s *Searcher) getPodNames(namespace string) ([]string, error) {
	pods, err := s.clientset.CoreV1().Pods(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var names []string
	for _, pod := range pods.Items {
		names = append(names, pod.Name)
	}
	return names, nil
}

// getConfigMapNames gets configmap names in a namespace.
func (s *Searcher) getConfigMapNames(namespace string) ([]string, error) {
	configmaps, err := s.clientset.CoreV1().ConfigMaps(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var names []string
	for _, cm := range configmaps.Items {
		names = append(names, cm.Name)
	}
	return names, nil
}

// getSecretNames gets secret names in a namespace.
func (s *Searcher) getSecretNames(namespace string) ([]string, error) {
	secrets, err := s.clientset.CoreV1().Secrets(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var names []string
	for _, secret := range secrets.Items {
		names = append(names, secret.Name)
	}
	return names, nil
}

// getServiceAccountNames gets serviceaccount names in a namespace.
func (s *Searcher) getServiceAccountNames(namespace string) ([]string, error) {
	serviceAccounts, err := s.clientset.CoreV1().ServiceAccounts(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var names []string
	for _, sa := range serviceAccounts.Items {
		names = append(names, sa.Name)
	}
	return names, nil
}

// getGenericResourceNames gets resource names for generic resources.
func (s *Searcher) getGenericResourceNames(namespace string) ([]string, error) {
	if s.kubeGet == nil {
		return nil, fmt.Errorf("kubeGet client not available")
	}

	kind := s.kind

	// Try the kind as-is first with go-kube-get
	_, resources, err := s.kubeGet.Get(context.Background(), kind, namespace)
	if err == nil {
		var names []string
		for _, resource := range resources.Items {
			names = append(names, resource.GetName())
		}
		return names, nil
	}

	// If that failed and it looks like resource.group format (like datasciencepipelinesapplications.opendatahub.io),
	// try just the resource name part
	if strings.Contains(kind, ".") && !strings.Contains(kind, ".v") {
		parts := strings.Split(kind, ".")
		if len(parts) >= 2 {
			resourceName := parts[0]
			_, resources, err := s.kubeGet.Get(context.Background(), resourceName, namespace)
			if err == nil {
				var names []string
				for _, resource := range resources.Items {
					names = append(names, resource.GetName())
				}
				return names, nil
			}
		}
	}

	// Return the original error if nothing worked
	return nil, fmt.Errorf("error getting %s resources: %v", s.kind, err)
}

// getResourceYAML gets the YAML representation of a specific resource.
func (s *Searcher) getResourceYAML(namespace, resource string) (string, error) {
	// If we have API version and kind, use generic resource handling
	if s.apiVersion != "" && s.kind != "" {
		return s.getGenericResourceYAML(namespace, resource)
	}

	// If we only have kind, check if it's a core resource first
	if s.kind != "" && s.apiVersion == "" {
		lowerKind := strings.ToLower(s.kind)
		switch lowerKind {
		case "pod", "pods":
			return s.getPodYAML(namespace, resource)
		case "configmap", "configmaps":
			return s.getConfigMapYAML(namespace, resource)
		case "secret", "secrets":
			return s.getSecretYAML(namespace, resource)
		case "serviceaccount", "serviceaccounts":
			return s.getServiceAccountYAML(namespace, resource)
		default:
			// For non-core resources, try auto-discovery
			return s.getGenericResourceYAML(namespace, resource)
		}
	}

	switch strings.ToLower(s.resourceType) {
	case "pods":
		return s.getPodYAML(namespace, resource)
	case "configmaps":
		return s.getConfigMapYAML(namespace, resource)
	case "secrets":
		return s.getSecretYAML(namespace, resource)
	case "serviceaccounts":
		return s.getServiceAccountYAML(namespace, resource)
	default:
		return "", fmt.Errorf("resource type %s not supported", s.resourceType)
	}
}

// getGenericResourceYAML gets YAML for generic resources.
func (s *Searcher) getGenericResourceYAML(namespace, name string) (string, error) {
	if s.kubeGet == nil {
		return "", fmt.Errorf("kubeGet client not available")
	}

	kind := s.kind

	// Try the kind as-is first with go-kube-get
	_, resources, err := s.kubeGet.Get(context.Background(), kind, namespace)
	if err == nil {
		// Find the specific resource by name
		for _, resource := range resources.Items {
			if resource.GetName() == name {
				return s.objectToYAML(&resource)
			}
		}
		return "", fmt.Errorf("%s %s not found", s.kind, name)
	}

	// If that failed and it looks like resource.group format (like datasciencepipelinesapplications.opendatahub.io),
	// try just the resource name part
	if strings.Contains(kind, ".") && !strings.Contains(kind, ".v") {
		parts := strings.Split(kind, ".")
		if len(parts) >= 2 {
			resourceName := parts[0]
			_, resources, err := s.kubeGet.Get(context.Background(), resourceName, namespace)
			if err == nil {
				// Find the specific resource by name
				for _, resource := range resources.Items {
					if resource.GetName() == name {
						return s.objectToYAML(&resource)
					}
				}
				return "", fmt.Errorf("%s %s not found", s.kind, name)
			}
		}
	}

	// Return error if nothing worked
	return "", fmt.Errorf("error getting %s resources: %v", s.kind, err)
}

// getPodYAML gets pod YAML.
func (s *Searcher) getPodYAML(namespace, name string) (string, error) {
	pod, err := s.clientset.CoreV1().Pods(namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		return "", err
	}
	return s.objectToYAML(pod)
}

// getConfigMapYAML gets configmap YAML.
func (s *Searcher) getConfigMapYAML(namespace, name string) (string, error) {
	configmap, err := s.clientset.CoreV1().ConfigMaps(namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		return "", err
	}
	return s.objectToYAML(configmap)
}

// getSecretYAML gets secret YAML.
func (s *Searcher) getSecretYAML(namespace, name string) (string, error) {
	secret, err := s.clientset.CoreV1().Secrets(namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		return "", err
	}
	return s.objectToYAML(secret)
}

// getServiceAccountYAML gets serviceaccount YAML.
func (s *Searcher) getServiceAccountYAML(namespace, name string) (string, error) {
	serviceAccount, err := s.clientset.CoreV1().ServiceAccounts(namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		return "", err
	}
	return s.objectToYAML(serviceAccount)
}

// objectToYAML converts a runtime.Object to YAML string.
func (s *Searcher) objectToYAML(obj runtime.Object) (string, error) {
	serializer := json.NewSerializerWithOptions(json.DefaultMetaFactory, nil, nil, json.SerializerOptions{Yaml: true})
	data, err := runtime.Encode(serializer, obj)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// discoverAPIVersionAndKind discovers the API version, correct kind, and resource name for a given kind name.
func (s *Searcher) discoverAPIVersionAndKind() (string, string, string, error) {
	if s.clientset == nil {
		return "", "", "", fmt.Errorf("kubernetes clientset not available")
	}

	// Get all API groups
	apiGroups, err := s.clientset.Discovery().ServerGroups()
	if err != nil {
		return "", "", "", fmt.Errorf("error getting API groups: %v", err)
	}

	// Search through all API groups and versions
	for _, group := range apiGroups.Groups {
		for _, version := range group.Versions {
			apiVersion := version.GroupVersion
			if group.Name == "" {
				apiVersion = version.Version // For core v1 resources
			}

			// Get resources for this API version
			resourceList, err := s.clientset.Discovery().ServerResourcesForGroupVersion(apiVersion)
			if err != nil {
				continue // Skip if we can't get resources for this version
			}

			// Look for a resource with matching kind
			for _, resource := range resourceList.APIResources {
				if strings.EqualFold(resource.Kind, s.kind) {
					return apiVersion, resource.Kind, resource.Name, nil
				}
			}
		}
	}

	return "", "", "", fmt.Errorf("could not find API version for kind '%s'", s.kind)
}
