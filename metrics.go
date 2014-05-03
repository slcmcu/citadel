package citadel

type (
	MemoryUsageMetric struct {
		Used  uint64 `json:"used,omitempty"`
		Total uint64 `json:"total,omitempty"`
	}

	CpuMetric struct {
		Nice uint64 `json:"nice,omitempty"`
		Sys  uint64 `json:"sys,omitempty"`
		User uint64 `json:"user,omitempty"`
		Wait uint64 `json:"wait,omitempty"`
	}
)
