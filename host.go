package citadel

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/citadel/citadel/utils"
	"github.com/cloudfoundry/gosigar"
	"github.com/coreos/go-etcd/etcd"
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

	logger   *logrus.Logger
	docker   *dockerclient.DockerClient
	registry *etcd.Client
}

func NewHost(labels []string, docker *dockerclient.DockerClient, logger *logrus.Logger) (*Host, error) {
	mem := sigar.Mem{}
	if err := mem.Get(); err != nil {
		return nil, err
	}

	cpus := sigar.CpuList{}
	if err := cpus.Get(); err != nil {
		return nil, err
	}

	id, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	h := &Host{
		ID:       id,
		Cpus:     len(cpus.List),
		Memory:   int(mem.Total / 1024 / 1024),
		Labels:   labels,
		docker:   docker,
		logger:   logger,
		registry: etcd.NewClient([]string{"http://127.0.0.1:4001"}),
	}

	if err := h.verifyState(); err != nil {
		return nil, err
	}

	if err := h.registerHost(); err != nil {
		return nil, err
	}

	docker.StartMonitorEvents(h.eventHandler, nil)

	return h, nil
}

// Close stops the events monitor
func (h *Host) Close() error {
	h.mux.Lock()
	defer h.mux.Unlock()

	h.docker.StopAllMonitorEvents()

	return h.deregisterHost()
}

func (h *Host) Containers() ([]*Container, error) {
	out := []*Container{}

	h.mux.Lock()

	resp, err := h.registry.Get(filepath.Join("/citadel", h.ID, "containers"), true, true)
	if err != nil {
		return nil, err
	}

	h.mux.Unlock()

	for _, node := range resp.Node.Nodes {
		var container *Container
		if err := json.Unmarshal([]byte(node.Value), &container); err != nil {
			return nil, err
		}

		out = append(out, container)
	}

	return out, nil
}

func (h *Host) Container(id string) (*Container, error) {
	h.mux.Lock()

	resp, err := h.registry.Get(filepath.Join("/citadel", h.ID, "containers", id), false, false)
	if err != nil {
		return nil, err
	}

	h.mux.Unlock()

	var container *Container
	if err := json.Unmarshal([]byte(resp.Node.Value), &container); err != nil {
		return nil, err
	}

	return container, nil
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

	data, err := json.Marshal(c)
	if err != nil {
		return err
	}

	if _, err := h.registry.Set(filepath.Join("/citadel", h.ID, "containers", c.ID), string(data), 0); err != nil {
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

	_, err = h.registry.Delete(filepath.Join("/citadel", h.ID, "containers", c.ID), false)

	if nerr := h.docker.RemoveContainer(c.ID); err == nil {
		err = nerr
	}

	return err
}

func (h *Host) inspect(id string) (*Container, error) {
	info, err := h.docker.InspectContainer(id)
	if err != nil {
		return nil, err
	}

	return containerFromDocker(h, info)
}

func (h *Host) verifyState() error {
	containers, err := h.Containers()
	if err != nil {
		return err
	}

	for _, c := range containers {
		info, err := h.docker.InspectContainer(c.ID)
		if err != nil {
			if err == dockerclient.ErrNotFound {
				h.logger.WithField("id", c.ID).Warn("container no longer exists in docker")

				if _, derr := h.registry.Delete(filepath.Join("/citadel", h.ID, "containers", c.ID), false); derr != nil {
					h.logger.WithField("error", derr).Warn("error deleting non-existant container")
				}

				continue
			}

			return err
		}

		if c.State.Status == Running && !info.State.Running {
			h.logger.WithField("id", c.ID).Warn("state mismatch")

			c.State.Status = Stopped
		}
	}

	return nil
}

func (h *Host) eventHandler(event *dockerclient.Event, _ ...interface{}) {
	switch event.Status {
	case "die":
		fromDocker, err := h.inspect(event.Id)
		if err != nil {
			h.logger.WithField("error", err).Error("fetch dead container information")
			return
		}

		h.mux.Lock()

		// only restart it if it's a managed container
		container, err := h.Container(fromDocker.ID)
		if err != nil {
			h.logger.WithField("error", err).Error("fetch container from registry")
			return
		}

		if container.Type == Service {
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

func (h *Host) registerHost() error {
	data, err := json.Marshal(h)
	if err != nil {
		return err
	}

	if _, err := h.registry.Set(filepath.Join("/citadel/hosts", h.ID), string(data), 0); err != nil {
		return err
	}

	return nil
}

func (h *Host) deregisterHost() error {
	_, err := h.registry.Delete(filepath.Join("/citadel/hosts", h.ID), false)
	return err
}
