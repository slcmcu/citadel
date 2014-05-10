package citadel

type Config struct {
	PollInterval   int      `json:"pull_interval"`
	InfluxHost     string   `json:"influx_host"`
	InfluxUser     string   `json:"influx_user"`
	InfluxPassword string   `json:"influx_password"`
	InfluxDatabase string   `json:"influx_database"`
	Namespace      string   `json:"namespace"`
	MasterTimeout  string   `json:"master_timeout,omitempty"`
	SlaveTTL       int      `json:"slave_ttl,omitempty"`
	Natsd          []string `json:"natsd,omitempty"`
}
