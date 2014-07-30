package citadel

type Container struct {
	Name        string            `json:"name"`
	Image       string            `json:"image,omitempty"`
	Cpus        float64           `json:"cpus,omitempty"`
	Memory      float64           `json:"memory,omitempty"`
	Environment map[string]string `json:"environment"`
	Hostname    string            `json:"hostname"`
	Domainname  string            `json:"domain"`
	Type        string            `json:"type,omitempty"`
	Labels      []string          `json:"labels,omitempty"`
}
