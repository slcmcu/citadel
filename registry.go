package citadel

type Registry interface {
	// FetchResources returns all the resources registered in the cluster with
	// up to date information about current reservations
	FetchResources() ([]*Resource, error)

	// Close releases any external demands by the registry
	Close() error
}
