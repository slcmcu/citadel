package citadel

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
	HostMetric struct {
		Load   map[string]float64 `json:"load" gorethink:"load"`
		Memory *MemoryUsageMetric `json:"memory" gorethink:"memory"` // bytes
		Disks  []*Disk            `json:"disks" gorethink:"disks"`
	}
)
