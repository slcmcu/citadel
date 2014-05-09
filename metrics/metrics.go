package metrics

const (
	ErrorAction   Action = "error"
	MessageAction Action = "message"

	Agent      System = "agent"
	Management System = "management"
	Scheduler  System = "scheduler"
	Docker     System = "docker"
)

type (
	Action string
	System string

	Metric struct {
		Time   float64 `json:"time,omitempty"`
		Memory *Memory `json:"memory,omitempty"`
		Cpu    *Cpu    `json:"cpu,omitempty"`
		Load1  float64 `json:"load_1,omitempty"`
		Load5  float64 `json:"load_5,omitempty"`
		Load15 float64 `json:"load_15,omitempty"`
	}

	Memory struct {
		Used  float64 `json:"used,omitempty"`
		Total float64 `json:"total,omitempty"`
	}

	Cpu struct {
		Nice float64 `json:"nice,omitempty"`
		Sys  float64 `json:"sys,omitempty"`
		User float64 `json:"user,omitempty"`
		Wait float64 `json:"wait,omitempty"`
	}

	Event struct {
		Time      float64 `json:"time,omitempty"`
		Host      string  `json:"host,omitempty"`
		Action    Action  `json:"action,omitempty"`
		Subsystem System  `json:"system,omitempty"`
		Content   string  `json:"content,omitempty"`
		Data      string  `json:"data,omitempty"`
	}
)
