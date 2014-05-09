package scheduler

import "net"

// Container defines the resources required to launch a container
type Container struct {
	Image   string            `json:"image,omitempty"`   // Image to use for the run
	Cpus    int               `json:"cpus,omitempty"`    // Number of cpus to reserve
	Memory  float64           `json:"memory,omitempty"`  // Amount of memory in bytes to reserve
	Volume  map[string]string `json:"volume,omitempty"`  // Use this volume mapping for storage
	Profile bool              `json:"profile,omitempty"` // Start profiling the container's execution
}

// State contains information about the running state of a container
type State struct {
	Container *Container `json:"container,omitempty"` // Container
	IP        net.IP     `json:"ip,omitempty"`        // IP of the running container
}
