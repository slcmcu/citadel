package citadel

import "time"

// Container represents a running container on the cluster
type Container struct {
	// ID is a unique identifier for a container
	ID string
	// Image is the base image name that the container is created from
	Image string
	// HostID is the host where the container is running
	HostID string
	// Cpus is the number of cpus that the container has reserved
	Cpus int
	// Memory is the amount of memory in mb that the container has reserved
	Memory int
	// StartedAt is the time when the container was launched
	StartedAt time.Time
}
