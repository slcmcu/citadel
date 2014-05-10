package scheduler

type Info struct {
	Containers     int      `json:"containers,omitempty"`      // Number of containers
	TotalCpus      int      `json:"total_cpus,omitempty"`      // Total number of cpus
	TotalMemory    float64  `json:"total_memory,omitempty"`    // Total RAM in bytes
	Volumes        []string `json:"volumes,omitempty"`         // Volumes residing on the host
	ReservedCpus   int      `json:"reserved_cpus,omitempty"`   // Cpus reserved for currently running containers
	ReservedMemory float64  `json:"reserved_memory,omitempty"` // Memory reserved for currently running containers
}

func (i *Info) AvailableCpus() int {
	return i.TotalCpus - i.ReservedCpus
}

func (i *Info) AvailableMemory() float64 {
	return i.TotalMemory - i.ReservedMemory
}

func (i *Info) CanAllocate(c *Container) bool {
	return (i.AvailableCpus()-c.Cpus) >= 0 && (i.AvailableMemory()-c.Memory) >= 0
}
