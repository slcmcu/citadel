package citadel

// Task represents a task submitting to modify a service
type Task struct {
	Name      string   `json:"name,omitempty"`
	Instances int      `json:"instances,omitempty"`
	Service   *Service `json:"service,omitempty"`
}
