package citadel

type Scheduler interface {
	Schedule(*Container) ([]*Resource, error)
}
