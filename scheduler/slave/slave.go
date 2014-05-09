package slave

import (
	"runtime"

	"citadelapp.io/citadel/scheduler"
)

// Slave that manages one docker host
type Slave struct {
	Cpus    int               `json:"cpus,omitempty"`
	Memory  float64           `json:"memory,omitempty"`
	Volumes map[string]string `json:"volumes,omitempty"`

	containers []*scheduler.Container
}

func NewSlave() (*Slave, error) {
	s := &Slave{
		Cpus: runtime.NumCPU(),
	}
	return s, nil
}

func (s *Slave) Info() (*scheduler.Info, error) {
	i := &scheduler.Info{
		Containers:  len(s.containers),
		TotalCpus:   s.Cpus,
		TotalMemory: s.Memory,
	}

	return i, nil
}
