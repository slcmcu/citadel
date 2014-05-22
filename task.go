package citadel

// Task represents a task submitting to modify a service
type Task struct {
	// Name is the full path of the service being requested
	Name string `json:"name,omitempty"`

	// Instance is the number of instances of the service to interact with
	Instances int `json:"instances,omitempty"`

	// Service is the data used to run or stop a service
	Service *ServiceData `json:"service,omitempty"`
}
