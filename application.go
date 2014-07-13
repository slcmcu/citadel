package citadel

type ApplicationType string

const (
	Undefined ApplicationType = ""        // Undefined type, will be treated like a batch task
	Service   ApplicationType = "service" // Long running task
	Batch     ApplicationType = "batch"   // One time short lived task
)

// Port represents a port mapping
type Port struct {
	Proto     string `json:"proto,omitempty"`
	Container int    `json:"container,omitempty"`
	Host      int    `json:"host,omitempty"`
}

type Application struct {
	// ID is the unique id for a specific application
	ID string `json:"id,omitempty"`
	// Image is the base image name that the containers are created from
	Image string `json:"image,omitempty"`
	// Type is the type of application, batch, service, etc
	Type ApplicationType `json:"type,omitempty"`
	// Instances is the number of instances of the application that you want to run
	Instances int `json:"instances,omitempty"`
	// Cpus is the number of cpus that the application has reserved
	Cpus []int `json:"cpus,omitempty"`
	// Memory is the amount of memory in mb that the application has reserved
	Memory int `json:"memory,omitempty"`
	// Ports is a collection of port mappings for an application
	Ports []*Port `json:"ports,omitempty"`
	// Args are additional arguments passed to the entrypoint of the containers
	Args []string `json:"args,omitempty"`
}
