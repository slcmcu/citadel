package citadel

import "log"

// ClusterManager manages changes to the state of the cluster
type ClusterManager struct {
	registry Registry
	executor Executor

	schedulers map[string]Scheduler

	logger *log.Logger
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
	switch c.Type {
	}

	return nil
}
