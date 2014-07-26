package citadel

type Scheduler interface {
	// Schedule returns valid resources for a Container to run on
	Schedule([]*Resource, *Container) ([]*Resource, error)
}
