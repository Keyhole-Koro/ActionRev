package main

import (
  "encoding/json"
  "log"
  "net/http"
  "os"
)

type rootResponse struct {
  Name   string `json:"name"`
  Status string `json:"status"`
}

func main() {
  mux := http.NewServeMux()

  mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    _ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
  })

  mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    _ = json.NewEncoder(w).Encode(rootResponse{
      Name:   "synthify-backend",
      Status: "running",
    })
  })

  port := os.Getenv("PORT")
  if port == "" {
    port = "8080"
  }

  log.Printf("synthify backend listening on :%s", port)
  if err := http.ListenAndServe(":"+port, mux); err != nil {
    log.Fatal(err)
  }
}
