package resource

// Occurrence represents an occurrence of a pattern in a Kubernetes resource.
type Occurrence struct {
	Resource  string
	Namespace string
	Line      int
	Content   string
}
