package master

import (
	"citadelapp.io/citadel"
	"citadelapp.io/citadel/slave"
)

type Placement struct {
	Slave     *slave.Slave       `json:"slave,omitempty"`
	Container *citadel.Container `json:"container,omitempty"`
}

type Schedule struct {
	Task       *citadel.Task `json:"task,omitempty"`
	Placements []*Placement  `json:"placements,omitempty"`
}

func PlaceOnSlave(s *slave.Slave, task *citadel.Task, cpus citadel.Cpus) *Placement {
	p := &Placement{
		Slave:     s,
		Container: task.Container,
	}

	if len(p.Container.Cpus) == 0 {
		setCpus(p, cpus)
	}
	return p
}

// place the task on the least used cpus
// by default we will place the task on 50% of the cpus
//
// TODO: this should be a plugin
func setCpus(p *Placement, cpus citadel.Cpus) {
	var (
		total = p.Slave.Cpus / 2
		used  = make(map[int]bool)
		set   = []int{}
	)

	if total == 0 {
		total = 1
	}

	citadel.SortCpusByCount(cpus)

	for i := 0; i < total; i++ {
		for _, c := range cpus {
			if !used[c.ID] {
				used[c.ID] = true
				set = append(set, c.ID)
			}
		}
	}
	p.Container.Cpus = set
}
