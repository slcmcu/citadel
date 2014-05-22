package slave

import (
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

	return slave, nil
}

func (s *Service) Data() *citadel.ServiceData {
	return s.data
}

func (s *Service) List(t *citadel.Task) ([]*citadel.ServiceData, error) {
	return s.client.List(t)
}

func (s *Service) Run(t *citadel.Task) (*citadel.RunResult, error) {
	return s.client.Run(t)
}

func (s *Service) Stop(t *citadel.Task) (*citadel.StopResult, error) {
	return s.client.Stop(t)
}
