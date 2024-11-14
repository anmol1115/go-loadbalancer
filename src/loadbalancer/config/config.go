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
		URL    string `yaml:"url"`
		Weight *int   `yaml:"weight,omitempty"`
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

func (c *Config) GetBackendServer() map[string]int {
	servers := make(map[string]int)
	
  for _, server := range c.Servers {
    switch c.Mode {
    case "RoundRobin", "Random":
      servers[server.URL] = 1

    case "WeightedRoundRobin":
      wt := 1
      if server.Weight != nil {
        wt = *server.Weight
        if wt == 0 {
          wt = 1
        }
      }
      servers[server.URL] = wt

    default:
  }
  }
  return servers
}

func (c *Config) GetLoadBalancer() string {
	return c.Mode
}

func (c *Config) GetServerPort() string {
	return fmt.Sprintf(":%d", c.Port)
}
