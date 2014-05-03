package citadel

import "time"

type (
	MemoryUsageMetric struct {
		Free  uint64 `json:"free,omitempty"`
		Used  uint64 `json:"used,omitempty"`
		Total uint64 `json:"total,omitempty"`
	}

	CpuMetric struct {
		Nice uint64 `json:"nice,omitempty"`
		Sys  uint64 `json:"sys,omitempty"`
		User uint64 `json:"user,omitempty"`
		Wait uint64 `json:"wait,omitempty"`
	}

	HostMetric struct {
		Name      string             `json:"name,omitempty"`
		Timestamp time.Time          `json:"timestamp,omitempty"`
		Load      map[string]float64 `json:"load,omitempty"`
		Memory    *MemoryUsageMetric `json:"memory,omitempty"`
		Disks     []*Disk            `json:"disks,omitempty"`
		Cpu       *CpuMetric         `json:"cpu,omitempty"`
	}
)
