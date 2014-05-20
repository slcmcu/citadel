package citadel

// Command is an action that can be run on a service
type Command interface {
}

type ListCommand interface {
	List(*Service) ([]*Service, error)
}

type StartCommand interface {
	Start(*Service) error
}

type StopCommand interface {
	Stop(*Service) error
}
