package citadel

type (
	Slave struct {
		ID     string  `json:"id,omitempty"`
		IP     string  `json:"ip,omitempty"`
		Cpus   int     `json:"cpus,omitempty"`
		Memory float64 `json:"memory,omitempty"`
	}

	Disk struct {
		Path      string  `json:"path,omitempty"`
		Used      float64 `json:"used,omitempty"`
		Free      float64 `json:"free,omitempty"`
		Files     float64 `json:"files,omitempty"`
		Available float64 `json:"available,omitempty"`
		Total     float64 `json:"total,omitempty"`
	}
)
