package citadel

import "time"

// Status of the container
type Status string

const (
	Created Status = "created"
	Running Status = "running"
	Stopped Status = "stopped"
)

// State represents the containers state
type State struct {
	// The current status of the container
	Status Status
	// StartedAt is the time when the container was launched
	StartedAt time.Time
	// ExitedAt is the time the container exited
	ExitedAt time.Time
	// ExitCode is the code which the container last returned on exit
	ExitCode int
}

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
	// State is the current state of the container
	State State
}
