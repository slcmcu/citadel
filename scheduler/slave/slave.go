package slave

import (
	"errors"
	"runtime"
	"sync"

	"citadelapp.io/citadel/scheduler"
	"github.com/Sirupsen/logrus"
	"github.com/samalba/dockerclient"
)

var (
	log = logrus.New()
)

var (
	ErrNotEnoughResources   = errors.New("resources not available to run")
	ErrProfilerNotSupported = errors.New("profiler not supported")
)

// Slave that manages one docker host
type Slave struct {
	sync.RWMutex

	Cpus    int               `json:"cpus,omitempty"`
	Memory  float64           `json:"memory,omitempty"`
	Volumes map[string]string `json:"volumes,omitempty"`

	containers scheduler.States
	docker     *dockerclient.DockerClient
}

func NewSlave(docker *dockerclient.DockerClient) (*Slave, error) {
	s := &Slave{
		Cpus:   runtime.NumCPU(),
		docker: docker,
	}

	s.docker.StartMonitorEvents(s.eventHandler)

	return s, nil
}

func (s *Slave) Info() (*scheduler.Info, error) {
	s.RLock()
	defer s.RUnlock()

	i := &scheduler.Info{
		Containers:     s.containers.Len(),
		TotalCpus:      s.Cpus,
		TotalMemory:    s.Memory,
		ReservedCpus:   s.containers.Cpus(),
		ReservedMemory: s.containers.Memory(),
		Volumes:        s.containers.Volumes(),
	}
	return i, nil
}

func (s *Slave) Running() scheduler.States {
	s.RLock()
	defer s.RUnlock()
	return s.containers
}

func (s *Slave) Execute(c *scheduler.Container) (*scheduler.State, error) {
	if err := s.canRun(c); err != nil {
		return nil, err
	}
	if c.Profile {
		// TODO: start profiler for the container
		return nil, ErrProfilerNotSupported
	}

	state := &scheduler.State{
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

	s.Lock()
	s.containers[id] = state
	s.Unlock()

	return state, nil
}

func (s *Slave) canRun(c *scheduler.Container) error {
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
			log.WithFields(logrus.Fields{
				"error": err,
				"event": event.Status,
				"id":    event.Id,
			}).Error("cannot remote container")
		}
	}
}
