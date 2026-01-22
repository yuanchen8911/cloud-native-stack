package serializer

// URI scheme constants for output destinations
const (
	// ConfigMapURIScheme is the URI scheme for Kubernetes ConfigMap destinations.
	// Format: cm://namespace/configmap-name
	ConfigMapURIScheme = "cm://"

	// StdoutURI is the special URI indicating output should be written to stdout.
	StdoutURI = "-"
)
