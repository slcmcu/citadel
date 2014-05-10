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
	citadel.Slave

	ID         string
	containers citadel.Containers
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

func (s *Slave) Execute(c *citadel.Container) error {
	if err := s.canRun(c); err != nil {
		return err
	}
	if c.Profile {
		// TODO: start profiler for the container
		return ErrProfilerNotSupported
	}

	config := &dockerclient.ContainerConfig{
		Image:     c.Image,
		Memory:    int(c.Memory),
		CpuShares: c.Cpus,
	}

	id, err := s.docker.CreateContainer(config)
	if err != nil {
		return err
	}
	if err := s.docker.StartContainer(id); err != nil {
		return err
	}
	c.ID = id

	s.Lock()
	s.containers[id] = c
	s.Unlock()

	return nil
}

func (s *Slave) canRun(c *citadel.Container) error {
	if len(c.Volumes) > 0 {
		return ErrVolumesNotSupported
	}

	s.RLock()
	defer s.RUnlock()

	var (
		reservedCpu    = s.containers.Cpus()
		reservedMemory = s.containers.Memory()
		// TODO: make this a plugin
		allocate = (s.Cpus-reservedCpu-c.Cpus) > 0 && (s.Memory-reservedMemory-c.Memory) > 0
	)

	if !allocate {
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
