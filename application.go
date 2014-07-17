package citadel

type RunType string

const (
	Undefined RunType = ""        // Undefined type, will be treated like a batch task
	Service   RunType = "service" // Long running task
	Batch     RunType = "batch"   // One time short lived task
	Group     RunType = "group"
)

// Port represents a port mapping
type Port struct {
	Proto     string `json:"proto,omitempty"`
	Container int    `json:"container,omitempty"`
	Host      int    `json:"host,omitempty"`
}

type Volume struct {
	Path string `json:"path,omitempty"`
	UID  int    `json:"uid,omitempty"`
	GID  int    `json:"gid,omitempty"`
}

type Config struct {
	// Image is the base image name that the containers are created from
	Image string `json:"image,omitempty"`
	// Type is the type of run strategy, batch, service, etc
	Type RunType `json:"type,omitempty"`
	// Cpus is the number of cpus that the application has reserved
	Cpus []int `json:"cpus,omitempty"`
	// Memory is the amount of memory in mb that the application has reserved
	Memory int `json:"memory,omitempty"`
	// Args are additional arguments passed to the entrypoint of the containers
	Args []string `json:"args,omitempty"`
}

type Application struct {
	// ID is the unique id for a specific application
	ID string `json:"id,omitempty"`
	// Ports is a collection of port mappings for the application
	Ports []*Port `json:"ports,omitempty"`
	// Volumes to be shared across all containers for the application
	Volumes []*Volume `json:"volumes,omitempty"`
	// Containers specifies the exact container requirements for the application
	Containers []*Config `json:"containers,omitempty"`
	// Constraints match with the labels on select hosts to place applications
	Constraints []string `json:"constraints,omitempty"`
}

func (a *Application) totalMemory() int {
	total := 0
	for _, c := range a.Containers {
		total += c.Memory
	}

	return total
}

func (a *Application) totalCpus() int {
	total := 0
	for _, c := range a.Containers {
		total += len(c.Cpus)
	}

	return total
}
