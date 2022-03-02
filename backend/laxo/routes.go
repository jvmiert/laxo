package laxo

import (
  "fmt"
  "net/http"
  "errors"
)

func HandleGetUser(w http.ResponseWriter, r *http.Request, uID string) {
  fmt.Fprintf(w, "Hello, your uID is: %s\n", uID)
}

func HandleLogout(w http.ResponseWriter, r *http.Request) {
  c, err := r.Cookie(AppConfig.AuthCookieName)

  if err == http.ErrNoCookie {
    http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
    return
  } else if err != nil {
    http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
    Logger.Error("Error in auth handler function (cookie parsing)", "error", err)
    return
  }

  err = RemoveUserSession(c.Value)
  if err != nil {
    http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
    return
  }

  RemoveUserCookie(w)
  fmt.Fprintf(w, "ok")
}

func HandleLogin(w http.ResponseWriter, r *http.Request) {
  var loginRequest LoginRequest

  if err := decodeJSONBody(w, r, &loginRequest); err != nil {
    var mr *malformedRequest
    if errors.As(err, &mr) {
      http.Error(w, mr.msg, mr.status)
    } else {
      http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
    }
    return
  }

  printer := getLocalePrinter(r)

  user, err := LoginUser(loginRequest.Email, loginRequest.Password, printer)

  if err != nil {
    ErrorJSONEncode(w, err, http.StatusUnauthorized)
    return
  }

  sessionKey, err := SetUserSession(user)

  if err != nil {
    http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
    return
  }

  user.SessionKey = sessionKey

  SetUserCookie(user.SessionKey, w)

  js, err := user.JSON()

  if err != nil {
    http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
    return
  }

  w.Write(js)
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

  printer := getLocalePrinter(r)
  if err := u.ValidateNew(printer); err != nil {
    ErrorJSONEncode(w, err, http.StatusUnprocessableEntity)
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
  SetUserCookie(u.SessionKey, w)

  js, err := u.JSON()

  if err != nil {
    http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
    return
  }

  w.Header().Set("Content-Type", "application/json; charset=utf-8")
  w.Write(js)
}

