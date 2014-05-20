package citadel

type NilCommand struct {
}

func (n *NilCommand) List(s *Service) ([]*Service, error) {
	return []*Service{}, nil
}

func (n *NilCommand) Start(s *Service) error {
	return nil
}

func (n *NilCommand) Stop(s *Service) error {
	return nil
}
