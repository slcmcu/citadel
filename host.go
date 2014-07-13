package citadel

import (
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
	// Address where the host's api can be reached
	Addr string `json:"addr,omitempty"`

	logger   *logrus.Logger
	docker   *dockerclient.DockerClient
	registry Registry
}

func NewHost(addr string, labels []string, etcdMachines []string, docker *dockerclient.DockerClient, logger *logrus.Logger) (*Host, error) {
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
		Addr:     addr,
		docker:   docker,
		logger:   logger,
		registry: NewRegistry(etcdMachines),
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

// RunContainer takes an application ID to lookup how a container is supposed to be run.
// A container is created with a unique ID and run on the host with the container saved back
// to the central registry
func (h *Host) RunContainer(applicationID string) *Transaction {
	h.mux.Lock()
	defer h.mux.Unlock()

	var (
		tran     = NewTransaction(RunTransaction)
		instance = 0
	)

	app, err := h.registry.FetchApplication(applicationID)
	if err != nil {
		return tran.Error(err)
	}

	config := &dockerclient.ContainerConfig{
		Image:  app.Image,
		Cmd:    app.Args,
		Memory: app.Memory * 1024 * 1024,
		Cpuset: utils.IToCpuset(app.Cpus),
	}

	name := fmt.Sprintf("citadel-%s-%d", app.ID, instance)

	id, err := h.docker.CreateContainer(config, name)
	if err != nil {
		return tran.Error(err)
	}

	c := &Container{
		ID:            id,
		ApplicationID: app.ID,
		HostID:        h.ID,
		Name:          name,
	}

	if err := h.startContainer(app, c); err != nil {
		return tran.Error(err)
	}

	if err := h.registry.SaveContainer(h.ID, c); err != nil {
		return tran.Error(err)
	}

	tran.Containers = append(tran.Containers, c)

	return tran
}

func (h *Host) startContainer(app *Application, c *Container) error {
	var hostConfig *dockerclient.HostConfig

	if app.Ports != nil {
		hostConfig = &dockerclient.HostConfig{
			PortBindings: make(map[string][]dockerclient.PortBinding),
		}

		for _, p := range app.Ports {
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

	info, err := h.docker.InspectContainer(c.ID)
	if err != nil {
		return err
	}

	ports, err := createPorts(info)
	if err != nil {
		return err
	}

	c.Ports = ports

	state := h.getState(info)
	c.State = state
	c.State.StartedAt = time.Now()
	c.State.ExitedAt = time.Time{}

	return nil
}

// StopContainer will stop the running container, remove it from the hosts registry
// and delete the container from docker
func (h *Host) StopContainer(id string) *Transaction {
	h.mux.Lock()
	defer h.mux.Unlock()

	tran := NewTransaction(StopTransaction)

	app, err := h.registry.FetchApplication(id)
	if err != nil {
		return tran.Error(err)
	}

	containers, err := h.registry.FetchContainers(h.ID)
	if err != nil {
		return tran.Error(err)
	}

	for _, c := range containers {
		if c.ApplicationID == app.ID {
			tran.Containers = append(tran.Containers, c)

			if err := h.docker.StopContainer(c.ID, 10); err != nil {
				return tran.Error(err)
			}

			info, err := h.docker.InspectContainer(c.ID)
			if err != nil {
				return tran.Error(err)
			}

			state := h.getState(info)
			c.State.ExitCode = state.ExitCode
			c.State.ExitedAt = time.Now()
			c.State.Status = Stopped

			h.registry.DeleteContainer(h.ID, c.ID)

			if err := h.docker.RemoveContainer(c.ID); err != nil {
				return tran.Error(err)
			}
		}
	}

	return tran
}

// Register ensures that the host can run the given application based on the requirements
func (h *Host) Register(id string) error {
	app, err := h.registry.FetchApplication(id)
	if err != nil {
		return err
	}

	return h.docker.PullImage(app.Image, "latest")
}

// getState inspects the container's state in docker and returns the updated information
func (h *Host) getState(info *dockerclient.ContainerInfo) State {
	state := State{}

	state.ExitCode = info.State.ExitCode

	if info.State.Running {
		state.Status = Running
	} else {
		state.Status = Stopped
	}

	return state
}

func (h *Host) verifyState() error {
	containers, err := h.registry.FetchContainers(h.ID)
	if err != nil {
		return err
	}

	for _, c := range containers {
		info, err := h.docker.InspectContainer(c.ID)
		if err != nil {
			if err == dockerclient.ErrNotFound {
				h.logger.WithField("id", c.ID).Warn("container no longer exists in docker")

				if derr := h.registry.DeleteContainer(h.ID, c.ID); derr != nil {
					h.logger.WithField("error", derr).Warn("error deleting non-existant container")
				}

				continue
			}

			return err
		}

		if c.State.Status == Running && !info.State.Running {
			h.logger.WithField("id", c.ID).Warn("state mismatch")

			c.State.Status = Stopped
			// TODO: make it run again
		}
	}

	return nil
}

func (h *Host) eventHandler(event *dockerclient.Event, _ ...interface{}) {
	switch event.Status {
	case "die":
		h.mux.Lock()
		defer h.mux.Unlock()

		// only restart it if it's a managed container
		container, err := h.registry.FetchContainer(h.ID, event.Id)
		if err != nil {
			h.logger.WithField("error", err).Error("fetch container from registry")

			return
		}

		app, err := h.registry.FetchApplication(container.ApplicationID)
		if err != nil {
			h.logger.WithField("error", err).Error("fetch container's application from registry")

			return
		}

		if app.Type == Service {
			container.State.ExitedAt = time.Now()

			if err := h.startContainer(app, container); err != nil {
				h.logger.WithField("error", err).Error("restarting dead container")
			}

			if err := h.registry.SaveContainer(h.ID, container); err != nil {
				h.logger.WithField("error", err).Error("saving after restart of dead container")
			}
		}

	default:
		h.logger.WithFields(logrus.Fields{
			"type": event.Status,
			"id":   event.Id,
			"from": event.From,
		}).Debug("docker event")
	}
}

func (h *Host) registerHost() error {
	return h.registry.SaveHost(h)
}

func (h *Host) deregisterHost() error {
	return h.registry.DeleteHost(h.ID)
}
