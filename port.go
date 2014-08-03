package citadel

type Port struct {
	Proto         string `json:"proto,omitempty"`
	Port          int    `json:"port,omitempty"`
	ImagePort int    `json:"container_port,omitempty"`
}
