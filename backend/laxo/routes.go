package laxo

import (
  "fmt"
  "net/http"
  "errors"
  "time"
)

func HandleGetUser(w http.ResponseWriter, r *http.Request, uID string) {
  fmt.Fprintf(w, "Hello, your uID is: %s\n", uID)
}

func HandleCreateUser(w http.ResponseWriter, r *http.Request) {
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

  if err := u.ValidateNew(); err != nil {
    Logger.Info("User validation error", "error", err)
    http.Error(w, err.Error(), http.StatusUnprocessableEntity)
    return
  }

  if err := SaveNewUserToDB(&u); err != nil {
    http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
    return
  }

  sessionKey, err := SetUserSession(&u)

  if err != nil {
    http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
    return
  }

  u.SessionKey = sessionKey

  // @todo: make this a helper function?
  expires := time.Now().AddDate(0, 0, AppConfig.AuthCookieExpire)

  authCookie := &http.Cookie{
    Name: AppConfig.AuthCookieName,
    Value: u.SessionKey,
    HttpOnly: true,
    Secure: true,
    Expires: expires,
  }

  http.SetCookie(w, authCookie)

  js, err := u.JSON()

  if err != nil {
    http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
    return
  }

  w.Write(js)
}

