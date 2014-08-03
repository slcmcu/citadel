package citadel

type Container struct {
	// ID is the container's id
	ID string `json:"id,omitempty"`

	// Image is the configuration from which the container was created
	Image *Image `json:"image,omitempty"`

	// Engine is the engine that is runnnig the container
	Engine *Engine `json:"engine",omitempty"`

	// Ports are the public port mappings for the container
	Ports []*Port `json:"ports,omitempty"`
}
