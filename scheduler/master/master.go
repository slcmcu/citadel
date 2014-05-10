package master

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"citadelapp.io/citadel"
	"citadelapp.io/citadel/repository"
)

var (
	ErrNoValidOffers = errors.New("no valid offers for tasks")
)

// Master is the master node in a cluster
type Master struct {
	sync.Mutex
	citadel.Master

	timeout time.Duration
}

type result struct {
	slaves []*citadel.Slave
	err    error
}

func New(uuid, addr string, timeout time.Duration) (*Master, error) {
	m := &Master{
		timeout: timeout,
	}
	m.Addr = addr
	m.ID = uuid
	return m, nil
}

func (m *Master) Schedule(task *citadel.Task, repo repository.Repository) ([]*citadel.Slave, error) {
	m.Lock()
	defer m.Unlock()

	complete := make(chan *result)

	go func() {
		r := &result{slaves: []*citadel.Slave{}}

		slaves, err := repo.FetchSlaves()
		if err != nil {
			r.err = fmt.Errorf("unable to fetch slaves %s", err)
			complete <- r
			return
		}

		for _, s := range slaves {
			containers, err := repo.FetchContainers(s.ID)
			if err != nil {
				r.err = fmt.Errorf("unable to fetch containers %s", err)
				complete <- r
				return
			}

			// TODO: make this a plugin
			if !containers.ContainsImage(task.Container.Image) {
				var (
					reservedCpu    = containers.Cpus()
					reservedMemory = containers.Memory()
					allocate       = (s.Cpus-reservedCpu-task.Container.Cpus) > 0 && (s.Memory-reservedMemory-task.Container.Memory) > 0
				)

				if allocate {
					r.slaves = append(r.slaves, s)
				}
			}
		}
		complete <- r
	}()

	select {
	case <-time.After(m.timeout):
		return nil, fmt.Errorf("no execution before timeout %s", m.timeout)
	case r := <-complete:
		if r.err != nil {
			return nil, r.err
		}

		if len(r.slaves) < task.Instances {
			return nil, ErrNoValidOffers
		}
		return r.slaves, nil
	}
}
