package master

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/Sirupsen/logrus"

	"citadelapp.io/citadel"
	"citadelapp.io/citadel/repository"
)

// MasterService is an http server to accepting service requests to execute
// on slaves in the cluster
type MasterService struct {
	sync.Mutex
	logger *logrus.Logger
	mux    *http.ServeMux
	repo   repository.Repository
}

func New(logger *logrus.Logger, repo repository.Repository) *MasterService {
	m := &MasterService{
		mux:    http.NewServeMux(),
		logger: logger,
		repo:   repo,
	}
	m.mux.HandleFunc("/run", m.runhandler)
	m.mux.HandleFunc("/stop", m.stophandler)

	return m
}

func (m *MasterService) ListenAndServe(addr string) error {
	return http.ListenAndServe(addr, m.mux)
}

func (m *MasterService) runhandler(w http.ResponseWriter, r *http.Request) {
	m.Lock()
	defer m.Unlock()

}

func (m *MasterService) stophandler(w http.ResponseWriter, r *http.Request) {
	m.Lock()
	defer m.Unlock()

	// find the slave that the service is running on
	// then issue a stop command to the slave

	task, err := m.unmarshalTask(r)
	if err != nil {
		m.httpError(w, err)
		return
	}

	service, err := m.repo.FetchService(task.Name)
	if err != nil {
		m.httpError(w, err)
		return
	}

	if err := service.Stop(); err != nil {
		m.httpError(w, err)
		return
	}
}

func (m *MasterService) unmarshalTask(r *http.Request) (*citadel.Task, error) {
	var t *citadel.Task
	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		return nil, err
	}
	return t, nil
}

func (m *MasterService) httpError(w http.ResponseWriter, err error) {
	m.logger.Error(err)
	http.Error(w, err.Error(), http.StatusInternalServerError)
}
