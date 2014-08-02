package citadel

type Container struct {
	Name        string              `json:"name,omitempty"`
	Image       string              `json:"image,omitempty"`
	Cpus        float64             `json:"cpus,omitempty"`
	Memory      float64             `json:"memory,omitempty"`
	Environment map[string]string   `json:"environment,omitempty"`
	Hostname    string              `json:"hostname,omitempty"`
	Domainname  string              `json:"domain,omitempty"`
	Type        string              `json:"type,omitempty"`
	Labels      []string            `json:"labels,omitempty"`
	UserData    map[string][]string `json:"user_data,omitempty"`
}
