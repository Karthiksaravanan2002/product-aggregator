package props

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Provider ProviderConfig `yaml:"provider"`
}

type ServerConfig struct {
	Port      int    `yaml:"port"`
	TimeoutMs int    `yaml:"timeout_ms"`
	BasePath  string `yaml:"base_path"`
}

type ProviderConfig struct {
	A ProviderEntry `yaml:"a"`
	B ProviderEntry `yaml:"b"`
	C ProviderEntry `yaml:"c"`
}

type ProviderEntry struct {
	Enabled   bool `yaml:"enabled"`
	TimeoutMs int  `yaml:"timeout_ms"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
