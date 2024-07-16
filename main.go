package main

import (
	"fmt"
	"net/http"
	"os"

	ChirpyDatabase "github.com/Couches/chirpy-database"
  "github.com/joho/godotenv"  
)

var config apiConfig = apiConfig{}

func main() {
  godotenv.Load()
  config.jwtSecret = os.Getenv("JWT_SECRET")
  
  result := ChirpyDatabase.NewDB("database.json")
  if result.Error != nil {
    return
  }

  database := (*result.Body).(ChirpyDatabase.Database)
  config.Database = database

	serverMux := http.NewServeMux()
	serverMux.Handle("/app/", config.middlewareMetrics(http.StripPrefix("/app/", http.FileServer(http.Dir("./app")))))

	for _, endpoint := range getEndpoints() {
		serverMux.HandleFunc(fmt.Sprintf("%v %v%v", endpoint.method, endpoint.namespace, endpoint.route), methodHandler(endpoint, endpoint.callback, config))
	}

	server := http.Server{
		Addr:    ":8080",
		Handler: serverMux,
	}
	server.ListenAndServe()
}

type apiConfig struct {
	pageVisits    int
	Database  ChirpyDatabase.Database
	jwtSecret     string
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
		// Chirps endpoints
		{
			method:    "POST",
			namespace: "/api",
			route:     "/chirps",
			callback:  endpointCreateChirp,
		},
		{
			method:    "GET",
			namespace: "/api",
			route:     "/chirps/{chirpID}",
			callback:  endpointGetChirp,
		},
		{
			method:    "GET",
			namespace: "/api",
			route:     "/chirps",
			callback:  endpointGetAllChirps,
		},
		// Users endpoints
		{
			method:    "POST",
			namespace: "/api",
			route:     "/users",
			callback:  endpointCreateUser,
		},
		{
			method:    "GET",
			namespace: "/api",
			route:     "/users/{userID}",
			callback:  endpointGetUser,
		},
		{
			method:    "GET",
			namespace: "/api",
			route:     "/users",
			callback:  endpointGetAllUsers,
		},
		{
			method:    "PUT",
			namespace: "/api",
			route:     "/users",
			callback:  endpointUpdateUserLogin,
		},
		// Login endpoints
		{
			method:    "POST",
			namespace: "/api",
			route:     "/login",
			callback:  loginEndpoint,
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
