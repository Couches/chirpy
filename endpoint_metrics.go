package main

import (
	"fmt"
	"net/http"
	"os"
)

func metricsEndpoint(responseWriter http.ResponseWriter, request *http.Request, config apiConfig) {
  htmlBytes, err := os.ReadFile("./admin/metrics.html")
  if err != nil {
    http.Error(responseWriter, "There was an error reading metrics.html", http.StatusInternalServerError)
  }

  html := string(htmlBytes)
  html = fmt.Sprintf(html, config.pageVisits)

  responseWriter.Header().Add("Content-Type", "text/html; charset=utf-8")
  responseWriter.WriteHeader(200)
  responseWriter.Write([]byte(html))
}

