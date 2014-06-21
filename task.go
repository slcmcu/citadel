package citadel

// Task represents a queue task
type Task struct {
	Id string `json:"id,omitempty" gorethink:"id,omitempty"`
	// command is the command to exec (i.e. start, stop, run, destroy, etc.)
	Command string `json:"command,omitempty" gorethink:"command,omitempty"`
	// hostId is the target for the task
	HostId string `json:"hostId,omitempty" gorethink:"hostId,omitempty"`
	// type is the type of task (i.e. container, image, host, etc.)
	Type string `json:"type,omitempty" gorethink:"type,omitempty"`
	// args are the arguments for the command
	// i.e. ["name=foo", "memory=128", "host=abc"]
	Args map[string]interface{} `json:"args,omitempty" gorethink:"args,omitempty"`
}
