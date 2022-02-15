package laxo

import (
  "context"
  "mime"
  "net/http"
  "time"
  "strings"

  "github.com/mediocregopher/radix/v4"
)

type AuthHandlerFunc func(w http.ResponseWriter, r *http.Request, u string)

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

// assureAuth handler function checks if the user has a cookie with a token in its request.
// Then the token is validated by checking if it exists in Redis. If it's valid,
// it will call the auth typed route function with the extra userID parameter. If the
// token does not exist in Redis or the token is not present in the cookie it will
// return a 403 forbidden code.
func assureAuth(handler AuthHandlerFunc) http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {
    c, err := r.Cookie(AppConfig.AuthCookieName)

    if err == http.ErrNoCookie {
      http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
      return
    } else if err != nil {
      http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
      Logger.Error("Error in auth handler function (cookie parsing)", "error", err)
      return
    }

    // Seeing if token is present in Redis
    var uID string

    ctx := context.Background()
    err = RedisClient.Do(ctx, radix.Cmd(&uID, "GET", c.Value))

    if err != nil {
      http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
      Logger.Error("Error in auth handler function (Redis)", "error", err)
      return
    }

    if uID == "" {
      c := &http.Cookie{
        Name:     AppConfig.AuthCookieName,
        Value:    "",
        Expires:  time.Unix(0, 0),
        Secure:   true,
        HttpOnly: true,
      }

      http.SetCookie(w, c)
      http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
      return
    }

    handler(w, r, uID)
  }
}

