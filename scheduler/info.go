package scheduler

type Info struct {
	Containers     int      `json:"containers,omitempty"`      // Number of containers
	TotalCpus      int      `json:"total_cpus,omitempty"`      // Total number of cpus
	TotalMemory    float64  `json:"total_memory,omitempty"`    // Total RAM in bytes
	Volumes        []string `json:"volumes,omitempty"`         // Volumes residing on the host
	ReservedCpus   int      `json:"reserved_cpus,omitempty"`   // Cpus reserved for currently running containers
	ReservedMemory float64  `json:"reserved_memory,omitempty"` // Memory reserved for currently running containers
}
