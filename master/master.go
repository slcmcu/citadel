package master

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"citadelapp.io/citadel"
	"citadelapp.io/citadel/slave"
)

var (
	ErrNoValidOffers = errors.New("no valid offers for tasks")
)

type Repository interface {
	FetchSlaves() ([]*slave.Slave, error)
	FetchContainers(string) (citadel.Containers, error)
}

// Master is the master node in a cluster
type Master struct {
	sync.Mutex

	ID      string        `json:"id,omitempty"`
	Addr    string        `json:"addr,omitempty"`
	Timeout time.Duration `json:"-"`
}

type result struct {
	slaves []*slave.Slave
	err    error
}

func New(uuid, addr string, timeout time.Duration) (*Master, error) {
	m := &Master{
		Timeout: timeout,
	}

	m.Addr = addr
	m.ID = uuid
	return m, nil
}

func (m *Master) Schedule(task *citadel.Task, repo Repository) ([]*slave.Slave, error) {
	m.Lock()
	defer m.Unlock()

	complete := make(chan *result)

	go func() {
		r := &result{slaves: []*slave.Slave{}}

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
	case <-time.After(m.Timeout):
		return nil, fmt.Errorf("no execution before timeout %s", m.Timeout)
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
