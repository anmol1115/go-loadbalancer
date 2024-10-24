package config

type Config struct {}

func LoadConfig(path string) *Config {
  return &Config{}
} 

func (c *Config) GetBackendServer() []string {
  return []string{"http://backend1:8080", "http://backend2:8080", "http://backend3:8080"}
}

func (c *Config) GetLoadBalancer() string {
  return "RoundRobin"
}

func (c *Config) GetServerPort() string {
  return ":8080"
}
