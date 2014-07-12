package citadel

import (
	"sync"

	"citadelapp.io/citadel/utils"
	"github.com/samalba/dockerclient"
)

// Host represents a host machine on the cluster
// running docker containers
type Host struct {
	// ID is a unique identifier for a host
	ID string `json:"id,omitempty"`
	// Cpus is the number of cpus the host has available
	Cpus int `json:"cpus,omitempty"`
	// Memory is the amount of memory in mb the host has available
	Memory int `json:"memory,omitempty"`
	// Label is specific attributes of a host
	Labels []string `json:"labels,omitempty"`

	docker     *dockerclient.DockerClient `json:"-"`
	containers []*Container
	mux        sync.Mutex
}

func NewHost(id string, cpus, memory int, labels []string, docker *dockerclient.DockerClient) (*Host, error) {
	h := &Host{
		ID:     id,
		Cpus:   cpus,
		Memory: memory,
		Labels: labels,
		docker: docker,
	}

	docker.StartMonitorEvents(h.eventHandler, nil)

	return h, nil
}

func (h *Host) eventHandler(event *dockerclient.Event, _ ...interface{}) {
	switch event.Status {
	case "start":
	case "die":
	case "kill":
	case "stop":
	case "pause":
	case "unpause":
	}
}

// Close stops the events monitor
func (h *Host) Close() error {
	h.docker.StopAllMonitorEvents()

	return nil
}

// GetContainers returns all containers on the host
func (h *Host) GetContainers() ([]*Container, error) {
	h.mux.Lock()
	defer h.mux.Unlock()

	dockerContainers, err := h.docker.ListContainers(true)
	if err != nil {
		return nil, err
	}

	containers := []*Container{}
	for _, dc := range dockerContainers {
		info, err := h.docker.InspectContainer(dc.Id)
		if err != nil {
			return nil, err
		}

		c := &Container{
			ID:     info.Id,
			Image:  utils.CleanImageName(dc.Image),
			HostID: h.ID,
			Cpus:   info.Config.CpuShares, // FIXME: not the right place, this is cpuset
		}

		if info.Config.Memory > 0 {
			c.Memory = info.Config.Memory / 1024 / 1024
		}

		if info.State.Running {
			c.State.Status = Running
		} else {
			c.State.Status = Stopped
		}
		c.State.ExitCode = info.State.ExitCode

		containers = append(containers, c)
	}

	return containers, nil
}

func (h *Host) RunContainer(c *Container) error {
	h.mux.Lock()
	defer h.mux.Unlock()

	config := &dockerclient.ContainerConfig{
		Image:     c.Image,
		Memory:    c.Memory * 1024 * 1024,
		CpuShares: c.Cpus,
	}

	id, err := h.docker.CreateContainer(config, c.ID)
	if err != nil {
		return err
	}

	return h.docker.StartContainer(id, nil)
}

func (h *Host) StopContainer(c *Container) error {
	h.mux.Lock()
	defer h.mux.Unlock()

	return h.docker.StopContainer(c.ID, 10)
}
