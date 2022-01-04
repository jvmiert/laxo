package laxo

import (
  "context"
  "fmt"
  "encoding/json"
  "net/http"
)

func handleLogin(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintf(w, "Hello, you've requested: %s\n", r.URL.Path)
}

func handleTest(w http.ResponseWriter, r *http.Request) {
  user, err := Queries.GetUser(context.Background(), 1)
  if err != nil {
    Logger.Error("Failed to retrieve user", "error", err)
  }

  js, err := json.Marshal(user)
  if err != nil {
    Logger.Error("Failed to marshal user json", "error", err)
  }

  w.Header().Set("Content-Type", "application/json")
  w.Write(js)
}

