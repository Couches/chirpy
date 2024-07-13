package main

import "net/http"

func resetEndpoint(responseWriter http.ResponseWriter, request *http.Request, config apiConfig) {
  responseWriter.Header().Add("Content-Type", "text/plain; charset=utf-8")
  responseWriter.WriteHeader(200)

  config.resetVisits()
}

