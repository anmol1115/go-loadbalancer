package main

import (
	"log"
	"net/http"
	"net/url"

	"loadBalancer/backend"
	"loadBalancer/config"
	"loadBalancer/strategy"
)

func handleError(e error) {
	panic(e.Error())
}

func main() {
	config := config.LoadConfig("path")
	loadBalancer := strategy.GetLoadBalancer(config.GetLoadBalancer())

	for _, server := range config.GetBackendServer() {
		serverURL, err := url.Parse(server)
		if err != nil {
			handleError(err)
		}

    strategy.AddBackends(loadBalancer, backend.CreateBackend(*serverURL))
	}

  server := http.Server{
    Addr: config.GetServerPort(),
    Handler: http.HandlerFunc(strategy.Serve(loadBalancer)),
  }
  if err := server.ListenAndServe(); err != nil {
    log.Fatal(err)
  }
}
