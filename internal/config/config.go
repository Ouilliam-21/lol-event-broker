package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Endpoints EndpointsConfig `yaml:"endpoints"`
	Database  DatabaseConfig  `yaml:"database"`
	Events    EventsConfig    `yaml:"events"`
	Players   PlayersConfig   `yaml:"players"`
}

type EndpointsConfig struct {
	LiveClient string `yaml:"live_client"`
	Droplet    string `yaml:"droplet"`
}

type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Name     string `yaml:"name"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

type EventsConfig struct {
	Watch []string `yaml:"watch"`
}

type PlayersConfig struct {
	Watch []string `yaml:"watch"`
}

func LoadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &config, nil
}

func (c *Config) GetWatchedPlayers() []string {
	return c.Players.Watch
}

func (c *Config) GetWatchedEvents() []string {
	return c.Events.Watch
}
