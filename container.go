package citadel

// Container defines the resources required to launch a container
type Container struct {
	ID        string            `json:"id,omitempty"`         // Container ID set by slave
	Image     string            `json:"image,omitempty"`      // Image to use for the run
	Cpus      []int             `json:"cpus,omitempty"`       // cpuset to assign the container
	CpuShares int               `json:"cpu_shares,omitempty"` // Number of cpu shares
	Memory    float64           `json:"memory,omitempty"`     // Amount of memory in bytes to reserve
	Volumes   map[string]string `json:"volumes,omitempty"`    // Use this volume mapping for storage
	Profiler  bool              `json:"profiler,omitempty"`   // Start profiling the container's execution
}

type Containers map[string]*Container

func (cs Containers) Len() int {
	return len(cs)
}

// Cpus returns a map where the cpuset is the key and the number
// of containers using that cpu as the value
func (cs Containers) Cpus() Cpus {
	out := make(map[int]int, len(cs))
	for _, c := range cs {
		for _, cpu := range c.Cpus {
			out[cpu]++
		}
	}
	cpus := Cpus{}
	for k, v := range out {
		cpus = append(cpus, &Cpu{ID: k, UsedBy: v})
	}
	return cpus
}

func (cs Containers) Memory() float64 {
	i := 0.0
	for _, c := range cs {
		i += c.Memory
	}
	return i
}

func (cs Containers) Volumes() []string {
	out := []string{}
	for _, c := range cs {
		for k := range c.Volumes {
			out = append(out, k)
		}
	}
	return out
}

func (cs Containers) ContainsImage(i string) bool {
	for _, c := range cs {
		if c.Image == i {
			return true
		}
	}
	return false
}
