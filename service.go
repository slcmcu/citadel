package citadel

import "errors"

var (
	ErrNoCommand = errors.New("command not available")
)

// Service represents a resource in a cluster that can be
// consumed
type Service struct {
	Name   string `json:"name,omitempty"`
	Type   string `json:"type,omitempty"`
	Addr   string `json:"addr,omitempty"`   // address to the service, can include port number
	Cpus   int    `json:"cpus,omitempty"`   // number of cpus assiged to the service
	Memory int    `json:"memory,omitempty"` // memory limit in mb
	Action Action `json:"-"`
}

// List returns the services that are children of the current
// service
func (s *Service) List() ([]*Service, error) {
	return s.Action.List(s)
}

func (s *Service) Start() error {
	return s.Action.Start(s)
}

func (s *Service) Stop() error {
	return s.Action.Stop(s)
}
