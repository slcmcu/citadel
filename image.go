package citadel

// Image is a template for running a docker container
type Image struct {
	// Name is the docker image to base the container off of
	Name string `json:"name,omitempty"`

	// Cpus is the number of cpu resources to give to the container
	Cpus float64 `json:"cpus,omitempty"`

	// Memory is the amount of memory in MB for the container
	Memory float64 `json:"memory,omitempty"`

	// Envionrment is the environment vars to set on the container
	Environment map[string]string `json:"environment,omitempty"`

	// Hostname is the host name to set for the container
	Hostname string `json:"hostname,omitempty"`

	// Domainname is the domain name to set for the container
	Domainname string `json:"domain,omitempty"`

	// Args are cli arguments to pass to the image
	Args []string `json:"args,omitempty"`

	// Type is the container type, often service, batch, etc...
	Type string `json:"type,omitempty"`

	// Labels are matched with constraints on the engines
	Labels []string `json:"labels,omitempty"`

	// BindPorts ensures that the container has exclusive access to the specified ports
	BindPorts []*Port `json:"bind_ports,omitempty"`

	// UserData is user defined data that is passed to the container
	UserData map[string][]string `json:"user_data,omitempty"`
}
