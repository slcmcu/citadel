package citadel

type Registry interface {
	// FetchResources returns all the resources registered in the cluster
	FetchResources() ([]*Resource, error)

	// GetTotalReservations returns the total cpu and memory reserved for the resource by id
	GetTotalReservations(id string) (cpus float64, memory float64, err error)

	// PlaceReservation ensures that the resource has the container's reservation locked
	PlaceReservation(string, *Container) error

	// ReleaseReservation releases any resource allocations for the resource
	ReleaseReservation(string, *Container) error
}
