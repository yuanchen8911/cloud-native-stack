package snapshotter

const (
	// APIDomain is the API domain for snapshot resources
	APIDomain = "cns.nvidia.com"

	// APIVersion is the current API version for snapshots
	APIVersion = "v1alpha1"

	// FullAPIVersion is the complete API version string
	FullAPIVersion = APIDomain + "/" + APIVersion

	// Kind is the resource kind for snapshots
	Kind = "Snapshot"
)
