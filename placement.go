package citadel

// Placement defines where the container was run and any runtime updates
// on port allocation and ips for the container
type Placement struct {
	// Engine is the engine that is runnnig the container
	Engine *Engine `json:"engine",omitempty"`

	// Ports are the public port mappings for the container
	Ports []*Port `json:"ports,omitempty"`
}
