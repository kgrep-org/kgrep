package log

// Message represents a log message from a Kubernetes pod.
type Message struct {
	PodName       string
	ContainerName string
	LineNumber    int
	Message       string
}
