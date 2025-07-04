package resource

import (
	"context"
	"fmt"
	"strings"

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
		kind:          resourceType,
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

	resources, err := s.getGenericResourceNames(namespace)
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
	yaml, err := s.getGenericResourceYAML(namespace, resource)
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

// getGenericResourceNames gets resource names for generic resources.
func (s *Searcher) getGenericResourceNames(namespace string) ([]string, error) {
	if s.kubeGet == nil {
		return nil, fmt.Errorf("kubeGet client not available")
	}

	kind := s.kind

	_, resources, err := s.kubeGet.Get(context.Background(), kind, namespace)
	if err == nil {
		var names []string
		for _, resource := range resources.Items {
			names = append(names, resource.GetName())
		}
		return names, nil
	}

	if namespace != "" {
		_, resources, err := s.kubeGet.Get(context.Background(), kind, "")
		if err == nil {
			var names []string
			for _, resource := range resources.Items {
				names = append(names, resource.GetName())
			}
			return names, nil
		}
	}

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
			if namespace != "" {
				_, resources, err := s.kubeGet.Get(context.Background(), resourceName, "")
				if err == nil {
					var names []string
					for _, resource := range resources.Items {
						names = append(names, resource.GetName())
					}
					return names, nil
				}
			}
		}
	}

	return nil, fmt.Errorf("error getting %s resources: %v", s.kind, err)
}

// getGenericResourceYAML gets YAML for generic resources.
func (s *Searcher) getGenericResourceYAML(namespace, name string) (string, error) {
	if s.kubeGet == nil {
		return "", fmt.Errorf("kubeGet client not available")
	}

	kind := s.kind

	_, resources, err := s.kubeGet.Get(context.Background(), kind, namespace)
	if err == nil {
		for _, resource := range resources.Items {
			if resource.GetName() == name {
				return s.objectToYAML(&resource)
			}
		}
		return "", fmt.Errorf("%s %s not found", s.kind, name)
	}

	if namespace != "" {
		_, resources, err := s.kubeGet.Get(context.Background(), kind, "")
		if err == nil {
			for _, resource := range resources.Items {
				if resource.GetName() == name {
					return s.objectToYAML(&resource)
				}
			}
			return "", fmt.Errorf("%s %s not found", s.kind, name)
		}
	}

	if strings.Contains(kind, ".") && !strings.Contains(kind, ".v") {
		parts := strings.Split(kind, ".")
		if len(parts) >= 2 {
			resourceName := parts[0]
			_, resources, err := s.kubeGet.Get(context.Background(), resourceName, namespace)
			if err == nil {
				for _, resource := range resources.Items {
					if resource.GetName() == name {
						return s.objectToYAML(&resource)
					}
				}
				return "", fmt.Errorf("%s %s not found", s.kind, name)
			}
			if namespace != "" {
				_, resources, err := s.kubeGet.Get(context.Background(), resourceName, "")
				if err == nil {
					for _, resource := range resources.Items {
						if resource.GetName() == name {
							return s.objectToYAML(&resource)
						}
					}
					return "", fmt.Errorf("%s %s not found", s.kind, name)
				}
			}
		}
	}

	return "", fmt.Errorf("error getting %s resources: %v", s.kind, err)
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

	apiGroups, err := s.clientset.Discovery().ServerGroups()
	if err != nil {
		return "", "", "", fmt.Errorf("error getting API groups: %v", err)
	}

	for _, group := range apiGroups.Groups {
		for _, version := range group.Versions {
			apiVersion := version.GroupVersion
			if group.Name == "" {
				apiVersion = version.Version
			}

			resourceList, err := s.clientset.Discovery().ServerResourcesForGroupVersion(apiVersion)
			if err != nil {
				continue
			}

			for _, resource := range resourceList.APIResources {
				if strings.EqualFold(resource.Kind, s.kind) {
					return apiVersion, resource.Kind, resource.Name, nil
				}
			}
		}
	}

	return "", "", "", fmt.Errorf("could not find API version for kind '%s'", s.kind)
}
