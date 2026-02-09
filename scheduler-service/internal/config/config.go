package config

import (
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

type Task struct {
	Name    string                 `yaml:"name"`
	Cron    string                 `yaml:"cron"`
	Topic   string                 `yaml:"topic"`
	Message map[string]interface{} `yaml:"message"`
}

type Config struct {
	Tasks []Task `yaml:"tasks"`
}

func Load(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var cfg Config
	err = yaml.Unmarshal(bytes, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
