package citadel

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/citadel/citadel/utils"
	"github.com/cloudfoundry/gosigar"
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

	// containers that were started with citadel
	managedContainers map[string]*Container

	logger *logrus.Logger
	docker *dockerclient.DockerClient
}

func NewHost(id string, labels []string, docker *dockerclient.DockerClient, logger *logrus.Logger) (*Host, error) {
	mem := sigar.Mem{}
	if err := mem.Get(); err != nil {
		return nil, err
	}

	cpus := sigar.CpuList{}
	if err := cpus.Get(); err != nil {
		return nil, err
	}

	h := &Host{
		ID:                id,
		Cpus:              len(cpus.List),
		Memory:            int(mem.Total / 1024 / 1024),
		Labels:            labels,
		docker:            docker,
		logger:            logger,
		managedContainers: make(map[string]*Container),
	}

	if err := h.loadState(); err != nil {
		return nil, err
	}

	if err := h.verifyState(); err != nil {
		return nil, err
	}

	docker.StartMonitorEvents(h.eventHandler, nil)

	return h, nil
}

// Close stops the events monitor and saves the host's state to disk
func (h *Host) Close() error {
	h.mux.Lock()
	defer h.mux.Unlock()

	h.docker.StopAllMonitorEvents()

	if err := h.saveState(); err != nil {
		return err
	}

	return nil
}

func (h *Host) Containers() []*Container {
	out := []*Container{}

	h.mux.Lock()

	for _, c := range h.managedContainers {
		out = append(out, c)
	}

	h.mux.Unlock()

	return out
}

func (h *Host) Container(id string) (*Container, error) {
	h.mux.Lock()

	c, exists := h.managedContainers[id]

	h.mux.Unlock()

	if !exists {
		return nil, fmt.Errorf("container %s does not exist", id)
	}

	return c, nil
}

func (h *Host) RunContainer(c *Container) error {
	h.mux.Lock()
	defer h.mux.Unlock()

	if _, exists := h.managedContainers[c.ID]; exists {
		return fmt.Errorf("container %s is already managed", c.ID)
	}

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

	h.managedContainers[c.ID] = c

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
	c.State.StartedAt = time.Now()
	c.State.ExitedAt = time.Time{}

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
	c.State.ExitedAt = time.Now()

	delete(h.managedContainers, c.ID)

	return h.docker.RemoveContainer(c.ID)
}

func (h *Host) inspect(id string) (*Container, error) {
	info, err := h.docker.InspectContainer(id)
	if err != nil {
		return nil, err
	}

	return containerFromDocker(h, info)
}

func (h *Host) saveState() error {
	h.logger.Debug("saving host state")

	f, err := os.OpenFile("host-state.json", os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	if err := json.NewEncoder(f).Encode(h.managedContainers); err != nil {
		return err
	}

	h.logger.Debug("host state saved")

	return nil
}

func (h *Host) loadState() error {
	f, err := os.Open("host-state.json")
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}

		return err
	}
	defer f.Close()

	h.logger.Debug("loading host state from disk")

	if err := json.NewDecoder(f).Decode(&h.managedContainers); err != nil {
		return err
	}

	return nil
}

func (h *Host) verifyState() error {
	for id, c := range h.managedContainers {
		info, err := h.docker.InspectContainer(id)
		if err != nil {
			if err == dockerclient.ErrNotFound {
				h.logger.WithField("id", id).Warn("container no longer exists in docker")

				delete(h.managedContainers, id)

				continue
			}

			return err
		}

		if c.State.Status == Running && !info.State.Running {
			h.logger.WithField("id", id).Warn("state mismatch")

			c.State.Status = Stopped
		}
	}

	return nil
}

func (h *Host) eventHandler(event *dockerclient.Event, _ ...interface{}) {
	switch event.Status {
	case "die":
		container, err := h.inspect(event.Id)
		if err != nil {
			h.logger.WithField("error", err).Error("fetch dead container information")
			return
		}
		h.mux.Lock()

		// only restart it if it's a managed container
		if c, exists := h.managedContainers[container.ID]; exists && c.Type == Service {
			container.State.ExitedAt = time.Now()

			if err := h.startContainer(container); err != nil {
				h.logger.WithField("error", err).Error("restarting dead container")
			}
		}

		h.mux.Unlock()
	default:
		h.logger.WithFields(logrus.Fields{
			"type": event.Status,
			"id":   event.Id,
			"from": event.From,
		}).Debug("docker event")
	}
}
