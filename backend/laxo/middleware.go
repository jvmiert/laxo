package laxo

import (
  "mime"
  "net/http"
  "strings"
)

func assureJSON(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
  contentType := r.Header.Get("Content-type")

	for _, v := range strings.Split(contentType, ",") {
		t, _, err := mime.ParseMediaType(v)
		if err != nil {
      Logger.Error("MIME parse error", "error", err)
      http.Error(w, http.StatusText(http.StatusUnsupportedMediaType), http.StatusUnsupportedMediaType)
      return
		}
    if t != "application/json" {
      http.Error(w, http.StatusText(http.StatusUnsupportedMediaType), http.StatusUnsupportedMediaType)
      return
    }
	}
  w.Header().Set("Content-Type", "application/json")
  next(w, r)
}
