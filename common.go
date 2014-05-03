package citadel

type (
	Host struct {
		Name        string  `json:"-"`
		IPAddress   string  `json:"ip_address,omitempty"`
		Cpus        int     `json:"cpus,omitempty"`
		TotalMemory uint64  `json:"total_memory"`
		Disks       []*Disk `json:"disks,omitempty"`
	}

	Disk struct {
		Path      string `json:"path,omitempty"`
		Used      uint64 `json:"used,omitempty"`
		Free      uint64 `json:"free,omitempty"`
		Files     uint64 `json:"files,omitempty"`
		Available uint64 `json:"available,omitempty"`
		Total     uint64 `json:"total,omitempty"`
	}
)
