package citadel

import "fmt"

type (
	// Container is a running instance
	Container struct {
		// ID is the container's id
		ID string `json:"id,omitempty"`

		// Name is the container's name
		Name string `json:"name,omitempty"`

		// Image is the configuration from which the container was created
		Image *Image `json:"image,omitempty"`

		// Engine is the engine that is runnnig the container
		Engine *Engine `json:"engine,omitempty"`

		// Ports are the public port mappings for the container
		Ports []*Port `json:"ports,omitempty"`

		// State is the docker state of the container
		State *State `json:"state,omitempty"`
	}
	// State is the container state
	State struct {
		Running   bool   `json:"running,omitempty"`
		Pid       int    `json:"pid,omitempty"`
		ExitCode  int    `json:"exit_code,omitempty"`
		StartedAt string `json:"started_at,omitempty"`
		Ghost     bool   `json:"ghost,omitempty"`
	}
)

func (c *Container) String() string {
	name := c.ID
	if c.Name != "" {
		name = c.Name
	}

	return fmt.Sprintf("container %s Image %s Engine %s", name, c.Image.Name, c.Engine.ID)
}
