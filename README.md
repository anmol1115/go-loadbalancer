# GO Load Balancer

This project is a simple load balancer implemented in Golang. It supports multiple load balancing strategies, including **Random**, **Round Robin**, and **Weighted Round Robin**.

## Features

- Supports three different load balancing strategies:
  - **Random**: Selects a backend server randomly.
  - **Round Robin**: Selects backend servers in a circular sequence.
  - **Weighted Round Robin**: Selects backend servers based on their assigned weights.
- Keeps track of unavailable backends and routes traffic accordingly.
- Builtin health checks to ensure inclusion of all alive backends.

## Configuration

The load balancer can be configured using a YAML file. Below are the different supported configurations.

### 1. Random Strategy

```yaml
port: 8080
mode: "Random"
servers:
  - url: "http://backend1:8080"
  - url: "http://backend2:8080"
  - url: "http://backend3:8080"
```

### 2. Round Robin Strategy
```yaml
port: 8080
mode: "RoundRobin"
servers:
  - url: "http://backend1:8080"
  - url: "http://backend2:8080"
  - url: "http://backend3:8080"
```

### 3. Weighted Round Robin Strategy
```yaml
port: 8080
mode: "WeightedRoundRobin"
servers:
  - url: "http://backend1:8080"
    weight: 2
  - url: "http://backend2:8080"
    weight: 2
  - url: "http://backend3:8080"
```

# Description:
- *port*: The port on which the load balancer will run.
- *mode*: The load balancing strategy.
- *servers*: A list of backend server URLs.
