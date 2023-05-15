package options

// Deploy holds docker stack deploy options
type Deploy struct {
	Composefiles     []string
	Namespace        string
	ResolveImage     string
	SendRegistryAuth bool
	Prune            bool
	Envfiles         []string
	NoInterpolate    bool
}
