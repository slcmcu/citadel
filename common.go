package citadel

import "time"

type (
	Disk struct {
		Path      string `json:"path" gorethink:"path"`
		Used      uint64 `json:"used" gorethink:"used"`
		Free      uint64 `json:"free" gorethink:"free"`
		Files     uint64 `json:"files" gorethink:"files"`
		Available uint64 `json:"available" gorethink:"available"`
		Total     uint64 `json:"total" gorethink:"total"`
	}
	Host struct {
		Name      string             `json:"name" gorethink:"name"`
		IPAddress string             `json:"ip_address" gorethink:"ip_address"`
		Cpus      int                `json:"cpus" gorethink:"cpus"`
		Memory    *MemoryUsageMetric `json:"memory" gorethink:"memory"` // bytes
		Disks     []*Disk            `json:"disks" gorethink:"disks"`
	}
	MemoryUsageMetric struct {
		Free  uint64 `json:"free" gorethink:"free"`
		Used  uint64 `json:"used" gorethink:"used"`
		Total uint64 `json:"total" gorethink:"total"`
	}
	CpuMetric struct {
		Nice uint64 `json:"nice" gorethink:"nice"`
		Sys  uint64 `json:"sys" gorethink:"sys"`
		User uint64 `json:"user" gorethink:"user"`
		Wait uint64 `json:"wait" gorethink:"wait"`
	}
	HostMetric struct {
		Name      string             `json:"name" gorethink:"name"`
		Timestamp time.Time          `json:"timestamp" gorethink:"timestamp"`
		Load      map[string]float64 `json:"load" gorethink:"load"`
		Memory    *MemoryUsageMetric `json:"memory" gorethink:"memory"` // bytes
		Disks     []*Disk            `json:"disks" gorethink:"disks"`
		Cpu       *CpuMetric         `json:"cpu" gorethink:"cpu"`
	}
)
