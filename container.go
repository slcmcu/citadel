package citadel

// Container defines the resources required to launch a container
type Container struct {
	Image   string            `json:"image,omitempty"`   // Image to use for the run
	Cpus    int               `json:"cpus,omitempty"`    // Number of cpus to reserve
	Memory  float64           `json:"memory,omitempty"`  // Amount of memory in bytes to reserve
	Volumes map[string]string `json:"volumes,omitempty"` // Use this volume mapping for storage
	Profile bool              `json:"profile,omitempty"` // Start profiling the container's execution
}

type Containers []*Container

func (cs Containers) Len() int {
	return len(cs)
}

func (cs Containers) Cpus() int {
	i := 0
	for _, c := range cs {
		i += c.Cpus
	}
	return i
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
