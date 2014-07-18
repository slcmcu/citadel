package citadel

type Container struct {
	Image  string  `json:"image,omitempty"`
	Cpus   float64 `json:"cpus,omitempty"`
	Memory float64 `json:"memory,omitempty"`
}
