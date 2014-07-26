package citadel

import (
	"errors"
	"log"
	"sync"
)

var (
	ErrNoSchedulerForType = errors.New("no scheduler installed for container type")
	ErrUnableToSchedule   = errors.New("unable to schedule container")
	ErrSchedulerExists    = errors.New("scheduler exists for type")
)

// ClusterManager manages changes to the state of the cluster
type ClusterManager struct {
	registry        Registry
	executor        Executor
	resourceManager *ResourceManager

	schedulers map[string]Scheduler

	logger *log.Logger
	mux    sync.Mutex
}

// NewClusterManager returns a new cluster manager initialized with the registry
// and a logger
func NewClusterManager(registry Registry, executor Executor, logger *log.Logger) *ClusterManager {
	return &ClusterManager{
		registry:        registry,
		executor:        executor,
		schedulers:      make(map[string]Scheduler),
		resourceManager: newResourceManger(registry),
		logger:          logger,
	}
}

// ScheduleContainer uses the schedulers registered with the cluster and finds
// a resource that is able to run the container.
//
// If not scheduling decision can be made an ErrUnableToSchedule error is returned.
func (m *ClusterManager) ScheduleContainer(c *Container) error {
	m.mux.Lock()
	defer m.mux.Unlock()

	m.logger.Printf("task=%q image=%q cpus=%f memory=%f type=%q\n", "schedule", c.Image, c.Cpus, c.Memory, c.Type)

	// find the correct scheduler for the container's type
	scheduler := m.schedulers[c.Type]

	if scheduler == nil {
		return ErrNoSchedulerForType
	}

	// let the scheduler make a decision about the hosts that it would like the container to
	// be executed on
	resources, err := scheduler.Schedule(c)
	if err != nil {
		return err
	}
	m.logger.Printf("task=%q image=%q resource.count=%d\n", "schedule", c.Image, len(resources))

	// check with the resource manager to ensure that the resources that the scheduler is able
	// to run the container and to place the container on the resource with the best utilization
	// score to maximize effenciency
	placement, err := m.resourceManager.PlaceContainer(resources, c)
	if err != nil {
		return err
	}
	m.logger.Printf("task=%q image=%q placement=%q\n", "schedule", c.Image, placement.Addr)

	// for the selected resource make sure that the resources are reserved for the container
	// and not allocated to anything else
	if err := m.registry.PlaceReservation(placement.ID, c); err != nil {
		return err
	}

	// run the container with the executor on the selected resource
	if err := m.executor.Run(placement, c); err != nil {
		return err
	}

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

	s.setRegistry(m.registry)

	m.schedulers[tpe] = s

	return nil
}
