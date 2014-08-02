package citadel

import (
	"errors"
	"log"
	"strconv"
	"strings"
	"sync"

	"github.com/rcrowley/go-metrics"
	"github.com/samalba/dockerclient"
)

var (
	ErrNoSchedulerForType = errors.New("no scheduler installed for container type")
	ErrUnableToSchedule   = errors.New("unable to schedule container")
	ErrSchedulerExists    = errors.New("scheduler exists for type")
)

// ClusterManager manages changes to the state of the cluster
type ClusterManager struct {
	engines         map[string]*Engine
	resourceManager *ResourceManager

	schedulers map[string]Scheduler

	// timer records the time it takes to comple a scheduling operation
	timer metrics.Timer

	logger *log.Logger
	mux    sync.Mutex
}

// NewClusterManager returns a new cluster manager initialized with the registry
// and a logger
func NewClusterManager(engines []*Engine, logger *log.Logger) *ClusterManager {
	m := &ClusterManager{
		engines:         make(map[string]*Engine),
		schedulers:      make(map[string]Scheduler),
		resourceManager: newEngineManger(logger),
		logger:          logger,
		timer:           metrics.NewTimer(),
	}

	for _, e := range engines {
		m.engines[e.ID] = e
	}

	metrics.Register("citadel-timer", m.timer)

	return m
}

// ScheduleContainer uses the schedulers registered with the cluster and finds
// a resource that is able to run the container.
//
// If not scheduling decision can be made an ErrUnableToSchedule error is returned.
func (m *ClusterManager) ScheduleContainer(c *Container) (*Transaction, error) {
	if err := ValidateContainer(c); err != nil {
		return nil, err
	}

	m.mux.Lock()

	t := newTransaction(c)
	defer func() {
		t.Close()
		m.mux.Unlock()
		m.timer.UpdateSince(t.Started)
	}()

	m.logger.Printf("task=%q image=%q cpus=%f memory=%f type=%q\n", "schedule", c.Image, c.Cpus, c.Memory, c.Type)

	// find the correct scheduler for the container's type
	scheduler := m.schedulers[c.Type]

	if scheduler == nil {
		return nil, ErrNoSchedulerForType
	}

	accepted := []*Engine{}

	for _, e := range m.engines {
		// ensure that we preload all the containers for an engine to be used in the scheduling decison
		if err := e.loadContainers(); err != nil {
			return nil, err
		}

		canrun, err := scheduler.Schedule(c, e)
		if err != nil {
			return nil, err
		}

		if canrun {
			accepted = append(accepted, e)
		}
	}

	m.logger.Printf("task=%q image=%q resource.count=%d\n", "schedule", c.Image, len(accepted))

	// check with the resource manager to ensure that the engines that the scheduler is able
	// to run the container and to place the container on the resource with the best utilization
	// score to maximize effenciency
	engine, err := m.resourceManager.PlaceContainer(c, accepted)
	if err != nil {
		return nil, err
	}

	m.logger.Printf("task=%q image=%q placement=%q\n", "schedule", c.Image, engine.Addr)

	t.place(engine)

	if err := m.runContainer(t); err != nil {
		return nil, err
	}

	return t, nil
}

// AddEngine adds a new engine to the cluster for use
func (m *ClusterManager) AddEngine(e *Engine) error {
	m.mux.Lock()
	defer m.mux.Unlock()

	m.engines[e.ID] = e

	return nil

}

// RemoveEngine removes the engine from the cluster
func (m *ClusterManager) RemoveEngine(e *Engine) error {
	m.mux.Lock()
	defer m.mux.Unlock()

	delete(m.engines, e.ID)

	return nil
}

// RegisterScheduler registers the scheduler for a specific container type within the
// cluster.  An ErrSchedulerExists error is returned if the cluster already has a
// scheduler registered for that specific type.
func (m *ClusterManager) RegisterScheduler(tpe string, s Scheduler) error {
	m.mux.Lock()
	defer m.mux.Unlock()

	m.logger.Printf("task=%q type=%q", "register\n", tpe)

	if _, exists := m.schedulers[tpe]; exists {
		return ErrSchedulerExists
	}

	m.schedulers[tpe] = s

	return nil
}

func (m *ClusterManager) runContainer(t *Transaction) error {
	if err := t.Container.Run(t.Placement.Engine); err != nil {
		return err
	}

	info, err := t.Placement.Engine.client.InspectContainer(t.Container.Name)
	if err != nil {
		return err
	}

	t.Placement.InternalIP = info.NetworkSettings.IpAddress

	for pp, b := range info.NetworkSettings.Ports {
		parts := strings.Split(pp, "/")
		rawPort, proto := parts[0], parts[1]

		port, err := strconv.Atoi(b[0].HostPort)
		if err != nil {
			return err
		}

		containerPort, err := strconv.Atoi(rawPort)
		if err != nil {
			return err
		}

		t.Placement.Ports = append(t.Placement.Ports, &Port{
			Proto:         proto,
			Port:          port,
			ContainerPort: containerPort,
		})
	}

	return nil
}

// ListContainers returns all the running containers in the cluster
func (m *ClusterManager) ListContainers() ([]*Container, error) {
	m.mux.Lock()
	defer m.mux.Unlock()

	containers := []*Container{}

	for _, engine := range m.engines {
		if err := engine.loadContainers(); err != nil {
			return nil, err
		}

		containers = append(containers, engine.containers...)
	}

	return containers, nil
}

// RemoveContainer will iterate over all the engines in the cluster and first kill
// the container then remove it complete from the engine
func (m *ClusterManager) RemoveContainer(c *Container) error {
	m.mux.Lock()
	defer m.mux.Unlock()

	for _, engine := range m.engines {
		if err := c.Kill(engine); err != nil {
			if err == dockerclient.ErrNotFound {
				continue
			}

			return err
		}

		if err := engine.client.RemoveContainer(c.Name); err != nil {
			return err
		}
	}

	return nil
}
