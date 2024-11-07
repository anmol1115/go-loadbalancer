package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Port    int    `yaml:"port"`
	Mode    string `yaml:"mode"`
	Servers []struct {
		URL string `yaml:"url"`
	} `yaml:"servers"`
}

func LoadConfig(path string) (*Config, error) {
  data, err := os.ReadFile(path)
  if err != nil {
    return nil, fmt.Errorf("Error opening file: %w", err)
  }

  var config Config
  if err := yaml.Unmarshal(data, &config); err != nil {
    return nil, fmt.Errorf("Error unmarshaling the yaml: %w", err)
  }

  return &config, nil
}

func (c *Config) GetBackendServer() []string {
  var servers []string
  for _, s := range c.Servers {
    servers = append(servers, s.URL)
  }

  return servers
}

func (c *Config) GetLoadBalancer() string {
	return c.Mode
}

func (c *Config) GetServerPort() string {
  return fmt.Sprintf(":%d", c.Port)
}
