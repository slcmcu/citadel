package slave

import (
	"path/filepath"

	"citadelapp.io/citadel"
	"citadelapp.io/citadel/repository"
)

type Service struct {
	repo   repository.Repository
	data   *citadel.ServiceData
	client citadel.Service
}

func New(data *citadel.ServiceData, client citadel.Service, repo repository.Repository) (citadel.Service, error) {
	if err := repo.SaveService(data.Name, data); err != nil {
		return nil, err
	}

	slave := &Service{
		repo:   repo,
		data:   data,
		client: client,
	}

	if err := slave.loadExisting(); err != nil {
		return nil, err
	}

	return slave, nil
}

func (s *Service) Data() *citadel.ServiceData {
	return s.data
}

func (s *Service) List(t *citadel.Task) ([]*citadel.ServiceData, error) {
	return s.client.List(t)
}

func (s *Service) Run(t *citadel.Task) (*citadel.RunResult, error) {
	result, err := s.client.Run(t)
	if err != nil {
		return nil, err
	}

	if err := s.repo.SaveService(filepath.Join(s.data.Name, t.Name), t.Service); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *Service) Stop(t *citadel.Task) (*citadel.StopResult, error) {
	return s.client.Stop(t)
}

func (s *Service) loadExisting() error {
	services, err := s.client.List(&citadel.Task{})
	if err != nil {
		return err
	}

	for _, service := range services {
		if err := s.repo.SaveService(filepath.Join(s.data.Name, service.Name), service); err != nil {
			return err
		}
	}

	return nil
}
