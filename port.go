package citadel

type Port struct {
	Proto string `json:"proto,omitempty"`
	Port  int    `json:"port,omitempty"`
}
