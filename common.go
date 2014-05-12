package citadel

type (
	Disk struct {
		Path      string  `json:"path,omitempty"`
		Used      float64 `json:"used,omitempty"`
		Free      float64 `json:"free,omitempty"`
		Files     float64 `json:"files,omitempty"`
		Available float64 `json:"available,omitempty"`
		Total     float64 `json:"total,omitempty"`
	}
)
