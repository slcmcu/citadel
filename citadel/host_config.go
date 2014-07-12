package main

import (
	"github.com/BurntSushi/toml"
)

type Config struct {
	ID     string   `toml:"id"`
	Cpus   int      `toml:"cpus"`
	Memory int      `toml:"memory"`
	Addr   string   `toml:"addr"`
	Docker string   `toml:"docker"`
	Labels []string `toml:"labels"`
}

func loadConfig(p string) (*Config, error) {
	var config *Config

	if _, err := toml.DecodeFile(p, &config); err != nil {
		return nil, err
	}

	return config, nil
}
