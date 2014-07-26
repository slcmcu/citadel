package citadel

type Registry interface {
	// GetTotalReservations returns the total cpu and memory reserved for the resource by id
	GetTotalReservations(id string) (cpus float64, memory float64, err error)

	// PlaceReservation ensures that the resource has the container's reservation locked
	PlaceReservation(string, *Container) error
}
