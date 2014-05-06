package citadel

type Config struct {
	PullInterval   int    `json:"pull_interval"`
	InfluxHost     string `json:"influx_host"`
	InfluxUser     string `json:"influx_user"`
	InfluxPassword string `json:"influx_password"`
	InfluxDatabase string `json:"influx_database"`
}
