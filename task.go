package citadel

type Task struct {
	Container *Container `json:"container,omitempty"`
	Instances int        `json:"instances,omitempty"`
}
