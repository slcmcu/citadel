package actions

import (
	"errors"

	"citadelapp.io/citadel"
)

var (
	Actions          = make(map[string]citadel.Action)
	ErrNoServiceType = errors.New("service type not registered")
)

func LoadActions(s *citadel.Service) error {
	action := Actions[s.Type]
	if action == nil {
		return ErrNoServiceType
	}
	s.Action = action
	return nil
}
