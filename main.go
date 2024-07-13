package main

import (
	"net/http"
)

var config apiConfig = apiConfig{}

func main() {
	serverMux := http.NewServeMux()
	serverMux.Handle("/app/", config.middlewareMetrics(http.StripPrefix("/app/", http.FileServer(http.Dir("./app")))))

	for _, endpoint := range getEndpoints() {
		serverMux.HandleFunc(endpoint.route, methodHandler(endpoint, endpoint.callback))
	}

	server := http.Server{
		Addr:    ":8080",
		Handler: serverMux,
	}
	server.ListenAndServe()
}

type apiConfig struct {
	pageVisits int
}

type httpEndpoint struct {
	method   string
	route    string
	callback func(http.ResponseWriter, *http.Request)
}

func getEndpoints() []httpEndpoint {
	return []httpEndpoint{
		{
			method:   "GET",
			route:    "/healthz",
			callback: healthCheckEndpoint,
		},
		{
			method:   "GET",
			route:    "/metrics",
			callback: metricsEndpoint,
		},
		{
			method:   "GET",
			route:    "/reset",
			callback: resetEndpoint,
		},
	}
}

func methodHandler(endpoint httpEndpoint, handlerFunc http.HandlerFunc) http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {
    if r.Method != endpoint.method {
      w.WriteHeader(http.StatusMethodNotAllowed)
      w.Write([]byte("405 - Method Not Allowed"))
      return
    }

    handlerFunc(w, r)
  }
}

func (config *apiConfig) middlewareMetrics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		config.pageVisits++
		next.ServeHTTP(w, r)
	})
}

func (config *apiConfig) resetVisits() {
	config.pageVisits = 0
}
