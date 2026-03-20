package config

import (
	"os"
	"strconv"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server ServerConfig `yaml:"server"`
	NATS   NATSConfig   `yaml:"nats"`
}

type ServerConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type NATSConfig struct {
	URL      string `yaml:"url"`
	Token    string `yaml:"token"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

func Load(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var cfg Config
	if err := yaml.NewDecoder(f).Decode(&cfg); err != nil {
		return nil, err
	}

	applyEnv(&cfg)
	return &cfg, nil
}

func applyEnv(cfg *Config) {
	if v := os.Getenv("VIBECHAT_HOST"); v != "" {
		cfg.Server.Host = v
	}
	if v := os.Getenv("VIBECHAT_PORT"); v != "" {
		if port, err := strconv.Atoi(v); err == nil {
			cfg.Server.Port = port
		}
	}
	if v := os.Getenv("VIBECHAT_NATS_URL"); v != "" {
		cfg.NATS.URL = v
	}
	if v := os.Getenv("VIBECHAT_NATS_TOKEN"); v != "" {
		cfg.NATS.Token = v
	}
	if v := os.Getenv("VIBECHAT_NATS_USER"); v != "" {
		cfg.NATS.User = v
	}
	if v := os.Getenv("VIBECHAT_NATS_PASSWORD"); v != "" {
		cfg.NATS.Password = v
	}
}
