package citadel

import "errors"

var (
	ErrNoCommand = errors.New("command not available")
)

// Service represents a resource in a cluster that can be
// consumed
type Service struct {
	Name     string             `json:"name,omitempty"`
	Type     string             `json:"type,omitempty"`
	Addr     string             `json:"addr,omitempty"`   // address to the service, can include port number
	Cpus     int                `json:"cpus,omitempty"`   // number of cpus assiged to the service
	Memory   int                `json:"memory,omitempty"` // memory limit in mb
	Commands map[string]Command `json:"-"`
}

// List returns the services that are children of the current
// service
func (s *Service) List() ([]*Service, error) {
	cmd := s.Commands["list"]
	if cmd == nil {
		return nil, ErrNoCommand
	}
	return cmd.(ListCommand).List(s)
}

func (s *Service) Start() error {
	cmd := s.Commands["start"]
	if cmd == nil {
		return ErrNoCommand
	}
	return cmd.(StartCommand).Start(s)
}

func (s *Service) Stop() error {
	cmd := s.Commands["stop"]
	if cmd == nil {
		return ErrNoCommand
	}
	return cmd.(StopCommand).Stop(s)
}
