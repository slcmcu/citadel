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
	registry Registry
	executor Executor

	schedulers map[string]Scheduler

	logger *log.Logger
	mux    sync.Mutex
}

// NewClusterManager returns a new cluster manager initialized with the registry
// and a logger
func NewClusterManager(registry Registry, executor Executor, logger *log.Logger) *ClusterManager {
	return &ClusterManager{
		registry:   registry,
		executor:   executor,
		schedulers: make(map[string]Scheduler),
		logger:     logger,
	}
}

// ScheduleContainer uses the schedulers registered with the cluster and finds
// a resource that is able to run the container.
//
// If not scheduling decision can be made an ErrUnableToSchedule error is returned.
func (m *ClusterManager) ScheduleContainer(c *Container) error {
	m.mux.Lock()
	defer m.mux.Unlock()

	scheduler := m.schedulers[c.Type]

	if scheduler == nil {
		return ErrNoSchedulerForType
	}

	return nil
}

// RegisterScheduler registers the scheduler for a specific container type within the
// cluster.  An ErrSchedulerExists error is returned if the cluster already has a
// scheduler registered for that specific type.
func (m *ClusterManager) RegisterScheduler(tpe string, s Scheduler) error {
	m.mux.Lock()
	defer m.mux.Unlock()

	if _, exists := m.schedulers[tpe]; exists {
		return ErrSchedulerExists
	}

	m.schedulers[tpe] = s

	return nil
}
