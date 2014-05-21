package citadel

// Action is a standard set of actions that a service can perform
type Action interface {
	List(*Service) ([]*Service, error)
	Start(*Service) error
	Stop(*Service) error
}
