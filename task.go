package citadel

// Task represents a queue task
type Task struct {
	ID string `json:"id,omitempty"`
	// command is the command to exec (i.e. start, stop, run, destroy, etc.)
	Command string `json:"command,omitempty"`
	// hostId is the target for the task
	Host string `json:"host,omitempty"`
	// type is the type of task (i.e. container, image, host, etc.)
	Type string `json:"type,omitempty"`

	Image       string `json:"image,omitempty"`
	Cpus        int    `json:"cpus,omitempty"`
	Memory      int    `json:"memory,omitempty"`
	Instances   int    `json:"instances,omitempty"`
	ContainerID string `json:"container_id,omitempty"`
}
