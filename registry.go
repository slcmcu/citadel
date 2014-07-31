package citadel

type Registry interface {
	// FetchDockers returns all the resources registered in the cluster with
	// up to date information about current reservations
	FetchDockers() ([]*Docker, error)

	// Close releases any external demands by the registry
	Close() error
}
