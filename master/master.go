package master

import (
	"fmt"
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

	// TODO: scheduling logic
	// right now we will just get all of our services and pick one at random
	services, err := s.repo.FetchServices(s.data.Name)
	if err != nil {
		return nil, err
	}

	if len(services) == 0 {
		return nil, fmt.Errorf("no services registered to run task")
	}

	client := handler.NewClient(services[0])

	// FIXME: master should save the new rundata to the service
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
