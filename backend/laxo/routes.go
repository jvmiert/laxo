package laxo

import (
  "fmt"
  "net/http"
  "errors"
)

func HandleGetUser(w http.ResponseWriter, r *http.Request, uID string) {
  fmt.Fprintf(w, "Hello, your uID is: %s\n", uID)
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

  user, err := RetrieveUserFromDBbyEmail(loginRequest.Email)
  if err == ErrUserNotExist {
    printer := getLocalePrinter(r)
    js, errJS := GetUserLoginFailure(true, false, printer)

    if errJS != nil {
      http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
      return
    }

    ErrorJSON(w, js, http.StatusUnauthorized)
    return
  } else if err != nil {
    http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
    return
  }


  if err = user.CheckPassword(loginRequest.Password); err != nil {
    printer := getLocalePrinter(r)
    js, errJS := GetUserLoginFailure(false, true, printer)

    if errJS != nil {
      http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
      return
    }

    ErrorJSON(w, js, http.StatusUnauthorized)
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
  SetUserCookie(u.SessionKey, w)

  js, err := u.JSON()

  if err != nil {
    http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
    return
  }

  w.Write(js)
}

