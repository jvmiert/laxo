package laxo

import (
  "fmt"
  "net/http"
)

func handleLogin(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintf(w, "Hello, you've requested: %s\n", r.URL.Path)
}

