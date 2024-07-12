package main

import "net/http"

func main() {
  serverMux := http.NewServeMux()
  serverMux.Handle("/app/", http.StripPrefix("/app/", http.FileServer(http.Dir("./app"))))
  serverMux.HandleFunc("/healthz", healthCheck)
  server := http.Server {
    Addr: ":8080",
    Handler: serverMux,
  }
  server.ListenAndServe()
}

func healthCheck(responseWriter http.ResponseWriter, request *http.Request) {
  responseWriter.Header().Add("Content-Type", "text/plain; charset=utf-8")
  responseWriter.WriteHeader(200)
  responseWriter.Write([]byte("OK"))
}
