package master

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"citadelapp.io/citadel"
)

var (
	ErrNoValidOffers = errors.New("no valid offers for tasks")
)

// Master is the master node in a cluster
type Master struct {
	sync.Mutex

	ID      string        `json:"id,omitempty"`
	Addr    string        `json:"addr,omitempty"`
	Timeout time.Duration `json:"-"`
}

func New(uuid, addr string, timeout time.Duration) (*Master, error) {
	m := &Master{
		Timeout: timeout,
	}

	m.Addr = addr
	m.ID = uuid
	return m, nil
}

func (m *Master) Schedule(task *citadel.Task, repo Repository) (*Schedule, error) {
	m.Lock()
	defer m.Unlock()

	// get the list of available slaves that could run the task
	placements := []*Placement{}

	slaves, err := repo.FetchSlaves()
	if err != nil {
		return nil, fmt.Errorf("unable to fetch slaves %s", err)
	}

	for _, s := range slaves {
		containers, err := repo.FetchContainers(s.ID)
		if err != nil {
			return nil, fmt.Errorf("unable to fetch containers %s", err)
		}

		// TODO: make this a plugin
		var (
			reservedMemory = containers.Memory()
			cpus           = containers.Cpus()
			allocate       = (s.Memory - reservedMemory - task.Container.Memory) > 0
		)

		if allocate {
			placements = append(placements, PlaceOnSlave(s, task, cpus))
		}
	}

	if len(placements) < task.Instances {
		return nil, ErrNoValidOffers
	}

	s := &Schedule{
		Task:       task,
		Placements: placements,
	}
	return s, nil
}
