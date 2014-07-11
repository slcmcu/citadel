package citadel

import (
	"citadelapp.io/citadel/utils"
	"github.com/samalba/dockerclient"
)

// Host represents a host machine on the cluster
// running docker containers
type Host struct {
	// ID is a unique identifier for a host
	ID string `json:"id,omitempty"`
	// Region is the region/datacenter that the host is running in
	Region string `json:"region,omitempty"`
	// Addr is the ip and port to connect to a specific host
	Addr string `json:"addr,omitempty"`
	// Cpus is the number of cpus the host has available
	Cpus int `json:"cpus,omitempty"`
	// Memory is the amount of memory in mb the host has available
	Memory int `json:"memory,omitempty"`

	Docker *dockerclient.DockerClient
}

// GetContainers returns all containers on the host
func (h *Host) GetContainers() ([]*Container, error) {
	dockerContainers, err := h.Docker.ListContainers(true)
	if err != nil {
		return nil, err
	}

	containers := []*Container{}

	for _, dc := range dockerContainers {
		info, err := h.Docker.InspectContainer(dc.Id)
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

func (h *Host) CreateContainer(task *Task) (string, error) {
	config := &dockerclient.ContainerConfig{
		Image:     task.Image,
		Memory:    task.Memory * 1024 * 1024,
		CpuShares: task.Cpus,
	}

	return h.Docker.CreateContainer(config, "")
}

func (h *Host) StartContainer(id string) error {
	return h.Docker.StartContainer(id, nil)
}

func (h *Host) StopContainer(id string) error {
	return h.Docker.StopContainer(id, 10)
}

func (h *Host) DeleteContainer(id string) error {
	return h.Docker.RemoveContainer(id)
}
