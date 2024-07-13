package main

import "net/http"
import "fmt"

func metricsEndpoint(responseWriter http.ResponseWriter, request *http.Request) {
  responseWriter.Header().Add("Content-Type", "text/plain; charset=utf-8")
  responseWriter.WriteHeader(200)
  responseWriter.Write([]byte(fmt.Sprintf("Hits: %v\n", config.pageVisits)))
}

