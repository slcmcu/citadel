package slave

import (
	"errors"
	"runtime"
	"sync"

	"citadelapp.io/citadel"
	"github.com/Sirupsen/logrus"
	"github.com/samalba/dockerclient"
)

var (
	ErrNotEnoughResources   = errors.New("resources not available to run")
	ErrProfilerNotSupported = errors.New("profiler not supported")
	ErrVolumesNotSupported  = errors.New("persistent storage not supported")
)

// Slave that manages one docker host
type Slave struct {
	sync.RWMutex
	citadel.Resource

	ID         string
	containers citadel.States
	docker     *dockerclient.DockerClient
	log        *logrus.Logger
}

func New(uuid string, logger *logrus.Logger, docker *dockerclient.DockerClient) (*Slave, error) {
	s := &Slave{
		ID:     uuid,
		docker: docker,
		log:    logger,
	}
	s.Cpus = runtime.NumCPU()

	s.docker.StartMonitorEvents(s.eventHandler)

	return s, nil
}

func (s *Slave) Info() (*citadel.Info, error) {
	s.RLock()
	defer s.RUnlock()

	i := &citadel.Info{
		Containers:     s.containers.Len(),
		TotalCpus:      s.Cpus,
		TotalMemory:    s.Memory,
		ReservedCpus:   s.containers.Cpus(),
		ReservedMemory: s.containers.Memory(),
		Volumes:        s.containers.Volumes(),
	}
	return i, nil
}

func (s *Slave) Running() citadel.States {
	s.RLock()
	defer s.RUnlock()
	return s.containers
}

func (s *Slave) Execute(c *citadel.Container) (*citadel.State, error) {
	if err := s.canRun(c); err != nil {
		return nil, err
	}
	if c.Profile {
		// TODO: start profiler for the container
		return nil, ErrProfilerNotSupported
	}

	state := &citadel.State{
		Container: c,
	}
	config := &dockerclient.ContainerConfig{
		Image:     c.Image,
		Memory:    int(c.Memory),
		CpuShares: c.Cpus,
	}

	id, err := s.docker.CreateContainer(config)
	if err != nil {
		return nil, err
	}
	if err := s.docker.StartContainer(id); err != nil {
		return nil, err
	}
	state.ID = id

	s.Lock()
	s.containers[id] = state
	s.Unlock()

	return state, nil
}

func (s *Slave) canRun(c *citadel.Container) error {
	if len(c.Volumes) > 0 {
		return ErrVolumesNotSupported
	}
	info, err := s.Info()
	if err != nil {
		return err
	}
	if !info.CanAllocate(c) {
		return ErrNotEnoughResources
	}
	return nil
}

func (s *Slave) eventHandler(event *dockerclient.Event, args ...interface{}) {
	switch event.Status {
	case "die", "stop", "kill":
		if err := s.docker.RemoveContainer(event.Id); err != nil {
			s.log.WithFields(logrus.Fields{
				"error": err,
				"event": event.Status,
				"id":    event.Id,
			}).Error("cannot remote container")
		}
		s.Lock()
		delete(s.containers, event.Id)
		s.Unlock()
	}
}
