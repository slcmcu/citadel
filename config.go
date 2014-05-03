package citadel

import (
	"github.com/BurntSushi/toml"
)

type Config struct {
	Listen         string   `toml:"listen"`
	PullInterval   int      `toml:"pull_interval"`
	Machines       []string `toml:"machines"`
	InfluxHost     string   `toml:"influx_host"`
	InfluxUser     string   `toml:"influx_user"`
	InfluxPassword string   `toml:"influx_password"`
	InfluxDatabase string   `toml:"influx_database"`
}

func LoadConfig(p string) (*Config, error) {
	var c *Config
	if _, err := toml.DecodeFile(p, &c); err != nil {
		return nil, err
	}
	return c, nil
}
