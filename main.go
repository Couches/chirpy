package main

import (
	"fmt"
	"net/http"

	ChirpDatabase "github.com/Couches/chirp-database"
)

var config apiConfig = apiConfig{
	Database: *ChirpDatabase.CreateDatabase[chirp]("database.json"),
}

func main() {
	serverMux := http.NewServeMux()
	serverMux.Handle("/app/", config.middlewareMetrics(http.StripPrefix("/app/", http.FileServer(http.Dir("./app")))))

	for _, endpoint := range getEndpoints() {
		serverMux.HandleFunc(fmt.Sprintf("%v%v", endpoint.namespace, endpoint.route), methodHandler(endpoint, endpoint.callback, config))
	}

	server := http.Server{
		Addr:    ":8080",
		Handler: serverMux,
	}
	server.ListenAndServe()
}

type apiConfig struct {
	pageVisits int
	Database   ChirpDatabase.Database[chirp]
}

type httpEndpoint struct {
	method    string
	namespace string
	route     string
	callback  func(http.ResponseWriter, *http.Request, apiConfig)
}

func getEndpoints() []httpEndpoint {
	return []httpEndpoint{
		{
			method:    "GET",
			namespace: "/api",
			route:     "/healthz",
			callback:  healthCheckEndpoint,
		},
		{
			method:    "GET",
			namespace: "/admin",
			route:     "/metrics",
			callback:  metricsEndpoint,
		},
		{
			method:    "GET",
			namespace: "/api",
			route:     "/reset",
			callback:  resetEndpoint,
		},
		{
			method:    "POST",
			namespace: "/api",
			route:     "/chirps",
			callback:  chirpsEndpoint,
		},
	}
}

func methodHandler(endpoint httpEndpoint, handlerFunc func(http.ResponseWriter, *http.Request, apiConfig), config apiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != endpoint.method {
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte("405 - Method Not Allowed"))
			return
		}

		handlerFunc(w, r, config)
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
