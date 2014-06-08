package citadel

// Host is a single host within the cluster
type Host struct {
	ID     string `json:"id,omitempty"`
	Addr   string `json:"addr,omitempty"`
	Cpus   int    `json:"cpus,omitempty"`
	Memory int    `json:"memory,omitempty"`
}

// NewHost returns a new host
func NewHost(id, addr string, cpus, memory int) *Host {
	return &Host{
		ID:     id,
		Addr:   addr,
		Cpus:   cpus,
		Memory: memory,
	}
}
