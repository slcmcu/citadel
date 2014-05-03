package metrics

type (
	Metric struct {
		Memory *Memory `json:"memory,omitempty"`
		Cpu    *Cpu    `json:"cpu,omitempty"`
		Load1  float64 `json:"load_1,omitempty"`
		Load5  float64 `json:"load_5,omitempty"`
		Load15 float64 `json:"load_15,omitempty"`
	}

	Memory struct {
		Used  uint64 `json:"used,omitempty"`
		Total uint64 `json:"total,omitempty"`
	}

	Cpu struct {
		Nice uint64 `json:"nice,omitempty"`
		Sys  uint64 `json:"sys,omitempty"`
		User uint64 `json:"user,omitempty"`
		Wait uint64 `json:"wait,omitempty"`
	}
)
