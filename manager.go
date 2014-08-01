package citadel

import (
	"errors"
	"log"
	"sync"
	"time"

	"github.com/rcrowley/go-metrics"
)

var (
	ErrNoSchedulerForType = errors.New("no scheduler installed for container type")
	ErrUnableToSchedule   = errors.New("unable to schedule container")
	ErrSchedulerExists    = errors.New("scheduler exists for type")
)

// ClusterManager manages changes to the state of the cluster
type ClusterManager struct {
	engines         []*Docker
	resourceManager *ResourceManager

	schedulers map[string]Scheduler

	timer metrics.Timer

	logger *log.Logger
	mux    sync.Mutex
}

// NewClusterManager returns a new cluster manager initialized with the registry
// and a logger
func NewClusterManager(engines []*Docker, logger *log.Logger) *ClusterManager {
	m := &ClusterManager{
		engines:         engines,
		schedulers:      make(map[string]Scheduler),
		resourceManager: newDockerManger(logger),
		logger:          logger,
		timer:           metrics.NewTimer(),
	}

	metrics.Register("citadel-timer", m.timer)

	return m
}

// ScheduleContainer uses the schedulers registered with the cluster and finds
// a resource that is able to run the container.
//
// If not scheduling decision can be made an ErrUnableToSchedule error is returned.
func (m *ClusterManager) ScheduleContainer(c *Container) (*Docker, error) {
	now := time.Now()
	defer m.timer.UpdateSince(now)

	m.logger.Printf("task=%q image=%q cpus=%f memory=%f type=%q\n", "schedule", c.Image, c.Cpus, c.Memory, c.Type)

	if err := ValidateContainer(c); err != nil {
		return nil, err
	}

	m.mux.Lock()
	defer m.mux.Unlock()

	var err error

	// find the correct scheduler for the container's type
	scheduler := m.schedulers[c.Type]

	if scheduler == nil {
		return nil, ErrNoSchedulerForType
	}

	engines := m.engines

	// let the scheduler make a decision about the hosts that it would like the container to
	// be executed on
	if engines, err = scheduler.Schedule(engines, c); err != nil {
		return nil, err
	}
	m.logger.Printf("task=%q image=%q resource.count=%d\n", "schedule", c.Image, len(engines))

	// check with the resource manager to ensure that the engines that the scheduler is able
	// to run the container and to place the container on the resource with the best utilization
	// score to maximize effenciency
	placement, err := m.resourceManager.PlaceContainer(engines, c)
	if err != nil {
		return nil, err
	}
	m.logger.Printf("task=%q image=%q placement=%q\n", "schedule", c.Image, placement.Addr)

	return placement, nil
}

func (m *ClusterManager) Engines() []*Docker {
	return m.engines
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
