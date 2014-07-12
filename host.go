package citadel

import (
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/Sirupsen/logrus"
	"github.com/citadel/citadel/utils"
	"github.com/samalba/dockerclient"
)

// Host represents a host machine on the cluster
// running docker containers
type Host struct {
	mux sync.Mutex

	// ID is a unique identifier for a host
	ID string `json:"id,omitempty"`
	// Cpus is the number of cpus the host has available
	Cpus int `json:"cpus,omitempty"`
	// Memory is the amount of memory in mb the host has available
	Memory int `json:"memory,omitempty"`
	// Label is specific attributes of a host
	Labels []string `json:"labels,omitempty"`

	docker *dockerclient.DockerClient `json:"-"`

	// containers that were started with citadel
	managedContainers map[string]*Container
	logger            *logrus.Logger
}

func NewHost(id string, cpus, memory int, labels []string, docker *dockerclient.DockerClient, logger *logrus.Logger) (*Host, error) {
	h := &Host{
		ID:                id,
		Cpus:              cpus,
		Memory:            memory,
		Labels:            labels,
		docker:            docker,
		logger:            logger,
		managedContainers: make(map[string]*Container),
	}

	docker.StartMonitorEvents(h.eventHandler, nil)

	return h, nil
}

func (h *Host) eventHandler(event *dockerclient.Event, _ ...interface{}) {
	switch event.Status {
	case "start":
	case "die":
		container, err := h.inspect(event.Id)
		if err != nil {
			h.logger.WithField("error", err).Error("fetch dead container information")
			return
		}
		h.mux.Lock()

		// only restart it if it's a managed container
		if c, exists := h.managedContainers[container.ID]; exists && c.Type == Service {
			if err := h.startContainer(container); err != nil {
				h.logger.WithField("error", err).Error("restarting dead container")
			}
		}

		h.mux.Unlock()
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
		c, err := h.inspect(dc.Id)
		if err != nil {
			return nil, err
		}

		managed, exists := h.managedContainers[c.ID]
		if exists {
			c.Type = managed.Type
		}

		containers = append(containers, c)
	}

	return containers, nil
}

func (h *Host) RunContainer(c *Container) error {
	h.mux.Lock()
	defer h.mux.Unlock()

	config := &dockerclient.ContainerConfig{
		Image:  c.Image,
		Memory: c.Memory * 1024 * 1024,
		Cpuset: utils.IToCpuset(c.Cpus),
	}

	if _, err := h.docker.CreateContainer(config, c.ID); err != nil {
		return err
	}

	if err := h.startContainer(c); err != nil {
		return err
	}

	return nil
}

func (h *Host) startContainer(c *Container) error {
	var hostConfig *dockerclient.HostConfig

	if c.Ports != nil {
		hostConfig = &dockerclient.HostConfig{
			PortBindings: make(map[string][]dockerclient.PortBinding),
		}

		for _, p := range c.Ports {
			proto := "tcp"
			if p.Proto != "" {
				proto = p.Proto
			}

			hostConfig.PortBindings[fmt.Sprintf("%d/%s", p.Container, proto)] = []dockerclient.PortBinding{
				{
					HostPort: fmt.Sprint(p.Host),
				},
			}
		}
	}

	if err := h.docker.StartContainer(c.ID, hostConfig); err != nil {
		return err
	}

	current, err := h.inspect(c.ID)
	if err != nil {
		return err
	}

	c.State = current.State

	h.managedContainers[c.ID] = c

	return nil
}

func (h *Host) StopContainer(c *Container) error {
	h.mux.Lock()
	defer h.mux.Unlock()

	if err := h.docker.StopContainer(c.ID, 10); err != nil {
		return err
	}

	// update the state on the original container so that when it is
	// returned it has the latest information
	current, err := h.inspect(c.ID)
	if err != nil {
		return err
	}
	c.State = current.State

	delete(h.managedContainers, c.ID)

	return h.docker.RemoveContainer(c.ID)
}

func (h *Host) inspect(id string) (*Container, error) {
	info, err := h.docker.InspectContainer(id)
	if err != nil {
		return nil, err
	}

	c := &Container{
		ID:     strings.TrimPrefix(info.Name, "/"),
		Image:  info.Image,
		HostID: h.ID,
		Cpus:   utils.CpusetTOI(info.Config.Cpuset),
	}

	// if cpuset is not specified then the container has all the cpus on the host
	if c.Cpus == 0 {
		c.Cpus = h.Cpus
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

	if info.HostConfig != nil && info.HostConfig.PortBindings != nil {
		for cp, bindings := range info.HostConfig.PortBindings {
			var (
				container int
				proto     string
			)

			if _, err := fmt.Sscanf(cp, "%d/%s", &container, &proto); err != nil {
				return nil, err
			}

			for _, b := range bindings {
				hostPort, err := strconv.Atoi(b.HostPort)
				if err != nil {
					return nil, err
				}

				c.Ports = append(c.Ports, &Port{
					Proto:     proto,
					Container: container,
					Host:      hostPort,
				})
			}
		}
	}

	return c, nil
}
