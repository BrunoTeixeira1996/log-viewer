package config

import (
	"os"

	"github.com/BurntSushi/toml"
)

type TargetConfig struct {
	Name string `toml:"name"`
	Host string `toml:"host"`
}

type Config struct {
	Targets []TargetConfig `toml:"targets"`
}

func ReadTomlFile(path string) (Config, error) {
	var cfg Config

	input, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}

	if _, err := toml.Decode(string(input), &cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}
