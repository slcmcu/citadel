package citadel

type Resource struct {
	Addr   string  `json:"addr,omitempty"`
	Cpus   float64 `json:"cpus,omitempty"`
	Memory float64 `json:"memory,omitempty"`
}
