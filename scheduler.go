package citadel

type Scheduler interface {
	Schedule(*Application, []*Host) ([]*Host, error)
}
