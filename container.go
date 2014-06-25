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
	Status Status `json:"status,omitempty"`
	// StartedAt is the time when the container was launched
	StartedAt time.Time `json:"started_at,omitempty`
	// ExitedAt is the time the container exited
	ExitedAt time.Time `json:"exited_at,omitempty"`
	// ExitCode is the code which the container last returned on exit
	ExitCode int `json:"exit_code,omitempty"`
}

// Container represents a running container on the cluster
type Container struct {
	// ID is a unique identifier for a container
	ID string `json:"id,omitempty"`
	// Image is the base image name that the container is created from
	Image string `json:"image,omitempty"`
	// HostID is the host where the container is running
	HostID string `json:"host_id,omitempty"`
	// Cpus is the number of cpus that the container has reserved
	Cpus int `json:"cpus,omitempty"`
	// Memory is the amount of memory in mb that the container has reserved
	Memory int `json:"memory,omitempty"`
	// State is the current state of the container
	State State `json:"state,omitempty"`
}
