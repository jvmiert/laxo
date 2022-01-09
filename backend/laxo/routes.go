package laxo

import (
  "fmt"
  "net/http"
  "errors"
)

func handleLogin(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintf(w, "Hello, you've requested: %s\n", r.URL.Path)
}

func handleCreateUser(w http.ResponseWriter, r *http.Request) {
  var u User

  if err := decodeJSONBody(w, r, &u.Model); err != nil {
    var mr *malformedRequest
    if errors.As(err, &mr) {
      http.Error(w, mr.msg, mr.status)
    } else {
      http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
    }
    return
  }

  if err := u.Validate(); err != nil {
    Logger.Info("User validation error", "error", err)
    http.Error(w, err.Error(), http.StatusUnprocessableEntity)
  }

  fmt.Fprint(w, "OK")
}

