package slave

import (
	"errors"
	"runtime"
	"sync"

	"citadelapp.io/citadel"
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

	ID     string  `json:"id,omitempty"`
	IP     string  `json:"ip,omitempty"`
	Cpus   int     `json:"cpus,omitempty"`
	Memory float64 `json:"memory,omitempty"`

	containers citadel.Containers
	docker     *dockerclient.DockerClient
}

func New(uuid string, docker *dockerclient.DockerClient) (*Slave, error) {
	s := &Slave{
		docker:     docker,
		containers: citadel.Containers{},
	}
	s.Cpus = runtime.NumCPU()
	s.Memory = 1024 * 8000
	s.ID = uuid

	return s, nil
}

func (s *Slave) Execute(c *citadel.Container) error {
	if c.Profiler {
		// TODO: start profiler for the container
		return ErrProfilerNotSupported
	}

	config := &dockerclient.ContainerConfig{
		Image:     c.Image,
		Memory:    int(c.Memory),
		CpuShares: c.CpuShares,
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

func (s *Slave) PullImage(image string) error {
	return s.docker.PullImage(image, "latest")
}

func (s *Slave) RemoveContainer(id string) error {
	s.Lock()
	delete(s.containers, id)
	s.Unlock()

	return s.docker.RemoveContainer(id)
}
