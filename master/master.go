package master

import (
	"sync"

	"citadelapp.io/citadel"
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

	return repo.FetchServices(t.Name)
}

func (s *Service) Run(t *citadel.Task) (interface{}, error) {
	s.Lock()
	defer s.Unlock()

	panic("not implemented")
}

func (s *Service) Stop(t *citadel.Task) (interface{}, error) {
	s.Lock()
	defer s.Unlock()

	// this should return a handle to the agent
	service, err := m.repo.FetchService(task.Name)
	if err != nil {
		m.httpError(w, err)
		return
	}

	if err := actions.LoadActions(service); err != nil {
		m.httpError(w, err)
		return
	}

	if err := service.Stop(); err != nil {
		m.httpError(w, err)
		return
	}
}
