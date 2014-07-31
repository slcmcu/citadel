package citadel

type Resource struct {
	ID             string   `json:"id,omitempty"`
	Addr           string   `json:"addr,omitempty"`
	Cpus           float64  `json:"cpus,omitempty"`
	Memory         float64  `json:"memory,omitempty"`
	Labels         []string `json:"labels,omitempty"`
	ReservedCpus   float64  `json:"reserved_cpus,omitempty"`
	ReservedMemory float64  `json:"reserved_memory,omitempty"`
}
