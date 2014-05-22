package master

import (
	"path/filepath"
	"sync"

	"citadelapp.io/citadel"
	"citadelapp.io/citadel/handler"
	"citadelapp.io/citadel/repository"
)

// Service is an http server to accepting service requests to execute
// on slaves in the cluster
type Service struct {
	sync.Mutex

	data *citadel.ServiceData
	repo repository.Repository
}

func New(data *citadel.ServiceData, repo repository.Repository) citadel.Service {
	return &Service{
		data: data,
		repo: repo,
	}
}

func (s *Service) Data() *citadel.ServiceData {
	return s.data
}

func (s *Service) List(t *citadel.Task) ([]*citadel.ServiceData, error) {
	s.Lock()
	defer s.Unlock()

	if t.Name == "" {
		t.Name = "/"
	}

	return s.repo.FetchServices(t.Name)
}

func (s *Service) Run(t *citadel.Task) (*citadel.RunResult, error) {
	s.Lock()
	defer s.Unlock()

	parent, _ := filepath.Split(t.Name)
	service, err := s.repo.FetchService(parent)
	if err != nil {
		return nil, err
	}

	client := handler.NewClient(service)

	return client.Run(t)
}

func (s *Service) Stop(t *citadel.Task) (*citadel.StopResult, error) {
	s.Lock()
	defer s.Unlock()

	parent, _ := filepath.Split(t.Name)
	service, err := s.repo.FetchService(parent)
	if err != nil {
		return nil, err
	}

	client := handler.NewClient(service)

	return client.Stop(t)
}
