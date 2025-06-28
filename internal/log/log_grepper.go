package log

import (
	"bufio"
	"context"
	"fmt"
	"strings"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
)

// LogReader is an interface for reading logs from a pod.
// This allows for swapping a fake implementation during testing.
type Reader interface {
	GetPodLogs(namespace, podName, containerName string) (string, error)
}

// DefaultLogReader is the production implementation of LogReader.
// It uses the real Kubernetes clientset to fetch logs.
type DefaultLogReader struct {
	clientset kubernetes.Interface
}

func (r *DefaultLogReader) GetPodLogs(namespace, podName, containerName string) (string, error) {
	req := r.clientset.CoreV1().Pods(namespace).GetLogs(podName, &corev1.PodLogOptions{
		Container: containerName,
	})

	logs, err := req.Do(context.Background()).Raw()
	if err != nil {
		return "", err
	}

	return string(logs), nil
}

// Grepper searches and filters logs from Kubernetes pods.
type Grepper struct {
	clientset kubernetes.Interface
	config    *rest.Config
	logReader Reader
}

// NewLogGrepper creates a new LogGrepper with a default configuration.
func NewLogGrepper() *Grepper {
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	configOverrides := &clientcmd.ConfigOverrides{}
	config, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides).ClientConfig()
	if err != nil {
		fmt.Printf("Error creating Kubernetes config: %v\n", err)
		// Return a grepper that can still function in a limited capacity (e.g., for tests without a real cluster).
		return &Grepper{}
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Printf("Error creating Kubernetes clientset: %v\n", err)
		return &Grepper{}
	}

	return &Grepper{
		clientset: clientset,
		config:    config,
		logReader: &DefaultLogReader{clientset: clientset},
	}
}

// GrepWithoutNamespace searches for a pattern in logs across all pods in the default namespace.
func (g *Grepper) GrepWithoutNamespace(pattern, sortBy string) []Message {
	namespace, err := g.getDefaultNamespace()
	if err != nil {
		fmt.Printf("Error getting default namespace: %v\n", err)
		return []Message{}
	}
	return g.Grep(namespace, "", pattern, sortBy)
}

// GrepResourceWithoutNamespace searches for a pattern in the logs of a specific resource in the default namespace.
func (g *Grepper) GrepResourceWithoutNamespace(resource, pattern, sortBy string) []Message {
	namespace, err := g.getDefaultNamespace()
	if err != nil {
		fmt.Printf("Error getting default namespace: %v\n", err)
		return []Message{}
	}
	return g.Grep(namespace, resource, pattern, sortBy)
}

// GrepNamespace searches for a pattern in logs across all pods in a specific namespace.
func (g *Grepper) GrepNamespace(namespace, pattern, sortBy string) []Message {
	return g.Grep(namespace, "", pattern, sortBy)
}

// Grep searches for a pattern in logs of a specific resource in a specific namespace.
func (g *Grepper) Grep(namespace, resource, pattern, sortBy string) []Message {
	if g.clientset == nil {
		fmt.Printf("Error: Kubernetes clientset not available\n")
		return []Message{}
	}

	pods, err := g.getPods(namespace, resource)
	if err != nil {
		fmt.Printf("Error getting pods: %v\n", err)
		return []Message{}
	}

	var messages []Message
	for _, pod := range pods {
		podMessages := g.searchPodLogs(pod, pattern)
		messages = append(messages, podMessages...)
	}

	return g.sortMessages(messages, sortBy)
}

// getDefaultNamespace gets the default namespace from kubeconfig.
func (g *Grepper) getDefaultNamespace() (string, error) {
	if g.config == nil {
		// This case is mainly for tests that don't provide a real config.
		// It avoids a panic and allows tests to proceed.
		// In a real scenario, NewLogGrepper would have failed or returned an empty grepper.
		return "default", nil
	}

	// Try to get namespace from kubeconfig
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	configOverrides := &clientcmd.ConfigOverrides{}
	kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)

	namespace, _, err := kubeConfig.Namespace()
	if err != nil {
		// If there's an error getting the namespace, default to "default".
		return "default", nil
	}

	if namespace == "" {
		namespace = "default"
	}
	return namespace, nil
}

// getPods gets pods in a namespace, optionally filtered by resource name.
func (g *Grepper) getPods(namespace, resource string) ([]corev1.Pod, error) {
	pods, err := g.clientset.CoreV1().Pods(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	if resource == "" {
		return pods.Items, nil
	}

	// Filter pods by resource name
	var filteredPods []corev1.Pod
	for _, pod := range pods.Items {
		if strings.Contains(pod.Name, resource) {
			filteredPods = append(filteredPods, pod)
		}
	}

	return filteredPods, nil
}

// searchPodLogs now uses the LogReader interface instead of a direct clientset call.
func (g *Grepper) searchPodLogs(pod corev1.Pod, pattern string) []Message {
	var messages []Message

	containers := g.getContainerNames(pod)

	for _, container := range containers {
		// Use the injected LogReader
		logs, err := g.logReader.GetPodLogs(pod.Namespace, pod.Name, container)
		if err != nil {
			fmt.Printf("Error getting logs for pod %s container %s: %v\n", pod.Name, container, err)
			continue
		}

		containerMessages := g.searchLogs(logs, pattern, pod.Name, container)
		messages = append(messages, containerMessages...)
	}

	return messages
}

// getContainerNames gets container names from a pod.
func (g *Grepper) getContainerNames(pod corev1.Pod) []string {
	var containers []string

	// Get containers from pod spec
	for _, container := range pod.Spec.Containers {
		containers = append(containers, container.Name)
	}

	// If no containers found in spec, try to get from status
	if len(containers) == 0 {
		for _, containerStatus := range pod.Status.ContainerStatuses {
			containers = append(containers, containerStatus.Name)
		}
	}

	return containers
}

// searchLogs searches for a pattern in log content.
func (g *Grepper) searchLogs(logs, pattern, podName, containerName string) []Message {
	var messages []Message

	// If the pattern is empty, we match every line.
	if pattern == "" {
		scanner := bufio.NewScanner(strings.NewReader(logs))
		lineNumber := 1
		for scanner.Scan() {
			messages = append(messages, Message{
				PodName:       podName,
				ContainerName: containerName,
				Message:       scanner.Text(),
				LineNumber:    lineNumber,
			})
			lineNumber++
		}
		return messages
	}

	// Otherwise, search for the pattern in each line.
	scanner := bufio.NewScanner(strings.NewReader(logs))
	lineNumber := 1
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(strings.ToLower(line), strings.ToLower(pattern)) {
			messages = append(messages, Message{
				PodName:       podName,
				ContainerName: containerName,
				Message:       line,
				LineNumber:    lineNumber,
			})
		}
		lineNumber++
	}

	return messages
}

// sortMessages sorts messages based on the sortBy parameter.
func (g *Grepper) sortMessages(messages []Message, sortBy string) []Message {
	switch strings.ToUpper(sortBy) {
	case "MESSAGE":
		// Sort by message content
		for i := 0; i < len(messages)-1; i++ {
			for j := i + 1; j < len(messages); j++ {
				if messages[i].Message > messages[j].Message {
					messages[i], messages[j] = messages[j], messages[i]
				}
			}
		}
	case "POD_AND_CONTAINER":
		// Default sort: by pod name, then container name, then line number
		for i := 0; i < len(messages)-1; i++ {
			for j := i + 1; j < len(messages); j++ {
				if messages[i].PodName > messages[j].PodName {
					messages[i], messages[j] = messages[j], messages[i]
				} else if messages[i].PodName == messages[j].PodName {
					if messages[i].ContainerName > messages[j].ContainerName {
						messages[i], messages[j] = messages[j], messages[i]
					} else if messages[i].ContainerName == messages[j].ContainerName {
						if messages[i].LineNumber > messages[j].LineNumber {
							messages[i], messages[j] = messages[j], messages[i]
						}
					}
				}
			}
		}
	}

	return messages
}
