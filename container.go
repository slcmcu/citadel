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
	// ExitCode is the code which the container last returned on exit
	ExitCode int `json:"exit_code"`
	// StartedAt is the time the container was last started
	StartedAt time.Time `json:"started_at,omitempty"`
	// ExitedAt is the time the container last exited.  It is set to zero when a
	// container is started again
	ExitedAt time.Time `json:"exited_at,omitempty"`
}

// Container represents a running container on the cluster along with runtime specific data
type Container struct {
	// ID is a unique identifier for a container
	ID string `json:"id,omitempty"`
	// ApplicationID is the id of the application that this container is based off of
	ApplicationID string `json:"app_id,omitempty"`
	// HostID is the host where the container is running
	HostID string `json:"host_id,omitempty"`
	// State is the current state of the container
	State State `json:"state,omitempty"`
	// Ports is a collection of port mappings for an application
	Ports []*Port `json:"ports,omitempty"`
}
