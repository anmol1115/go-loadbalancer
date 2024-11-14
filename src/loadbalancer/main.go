package main

import (
	"context"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"loadBalancer/backend"
	"loadBalancer/config"
	"loadBalancer/strategy"
)

const (
	Attempts int = iota
	Retry
)

func GetAttemptsFromContext(r *http.Request) int {
	if attempts, ok := r.Context().Value(Attempts).(int); ok {
		return attempts
	}

	return 1
}

func GetRetryFromContext(r *http.Request) int {
	if retries, ok := r.Context().Value(Retry).(int); ok {
		return retries
	}

	return 0
}

func handleError(e error) {
	panic(e.Error())
}

func main() {
	config, err := config.LoadConfig("./config.yaml")
  if err != nil {
    handleError(err)
  }
	loadBalancer := strategy.GetLoadBalancer(config.GetLoadBalancer())

	for server, weight := range config.GetBackendServer() {
		serverURL, err := url.Parse(server)
		if err != nil {
			handleError(err)
		}
		proxy := httputil.NewSingleHostReverseProxy(serverURL)
		proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, e error) {
			log.Printf("[%s] %s\n", serverURL.Host, e.Error())
			retries := GetRetryFromContext(r)

			if retries < 3 {
				select {
				case <-time.After(time.Millisecond * 10):
					ctx := context.WithValue(r.Context(), Retry, retries+1)
					proxy.ServeHTTP(w, r.WithContext(ctx))
				}

				return
			}

			strategy.MarkBackendStateus(serverURL, false, loadBalancer)

			attempts := GetAttemptsFromContext(r)
			log.Printf("%s(%s) Attempting Retry %d\n", r.RemoteAddr, r.URL.Path, attempts)
			ctx := context.WithValue(r.Context(), Attempts, attempts+1)
			strategy.Serve(loadBalancer)(w, r.WithContext(ctx))
		}

		strategy.AddBackends(loadBalancer, backend.CreateBackend(*serverURL, proxy, weight))
	}

	go func() {
		t := time.NewTicker(time.Minute * 2)
		for {
			select {
			case <-t.C:
				strategy.HealthCheck(loadBalancer)
			}
		}
	}()

	http.HandleFunc("/", strategy.Serve(loadBalancer))
	http.HandleFunc("/backendstatus", strategy.BackendStatusHandler(loadBalancer))
	if err := http.ListenAndServe(config.GetServerPort(), nil); err != nil {
    handleError(err)
	}
}
