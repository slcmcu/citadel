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

	scheduler := m.schedulers[c.Type]

	if scheduler == nil {
		return ErrNoSchedulerForType
	}

	resources, err := scheduler.Schedule(c)
	if err != nil {
		return err
	}
	m.logger.Printf("task=%q image=%q resource.count=%d\n", "schedule", c.Image, len(resources))

	placement, err := m.resourceManager.PlaceContainer(resources, c)
	if err != nil {
		return err
	}
	m.logger.Printf("task=%q image=%q placement=%q score=%f\n", "schedule", c.Image, placement.r.Addr, placement.Score)

	if err := m.executor.Run(placement.r, c); err != nil {
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

	m.schedulers[tpe] = s

	return nil
}
