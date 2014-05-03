package main

import (
	"github.com/BurntSushi/toml"
)

type config struct {
	Listen         string   `toml:"listen"`
	PullInterval   int      `toml:"pull_interval"`
	Machines       []string `toml:"machines"`
	InfluxHost     string   `toml:"influx_host"`
	InfluxUser     string   `toml:"influx_user"`
	InfluxPassword string   `toml:"influx_password"`
	InfluxDatabase string   `toml:"influx_database"`
}

func loadConfig(p string) (*config, error) {
	var c *config
	if _, err := toml.DecodeFile(p, &c); err != nil {
		return nil, err
	}
	return c, nil
}
