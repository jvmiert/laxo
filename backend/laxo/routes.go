package laxo

import (
	"errors"
	"fmt"
	"net/http"
)

func HandleGetUser(w http.ResponseWriter, r *http.Request, uID string) {
  fmt.Fprintf(w, "Hello, your uID is: %s\n", uID)
}

func HandleTest(w http.ResponseWriter, r *http.Request, uID string) {
  wfID, err := startTask("cool shop")

  if err != nil {
    http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
    return
  }

  fmt.Fprint(w, wfID)
}

func HandleVerifyOAuth(w http.ResponseWriter, r *http.Request, uID string) {
  var o OAuthVerifyRequest

  if err := decodeJSONBody(w, r, &o); err != nil {
    var mr *malformedRequest
    if errors.As(err, &mr) {
      http.Error(w, mr.msg, mr.status)
    } else {
      http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
    }
    return
  }

  printer := getLocalePrinter(r)
  if err := o.Verify(uID, printer); err != nil {
    ErrorJSONEncode(w, err, http.StatusUnprocessableEntity)
    return
  }

  fmt.Fprintf(w, "Hello, your uID is: %s\n", uID)
}

func HandleOAuthRedirects(w http.ResponseWriter, r *http.Request, uID string) {
  shopID := r.URL.Query().Get("shopID")
  o := &OAuthRedirectRequest{ShopID: shopID}

  printer := getLocalePrinter(r)
  if err := o.Validate(uID, printer); err != nil {
    ErrorJSONEncode(w, err, http.StatusUnprocessableEntity)
    return
  }

  err := o.GenerateRedirect()

  if err != nil {
    http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
    return
  }

  js, err := o.JSON()

  if err != nil {
    http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
    return
  }

  w.Header().Set("Content-Type", "application/json; charset=utf-8")
  w.Write(js)
}

func HandleGetMyShops(w http.ResponseWriter, r *http.Request, uID string) {
  shops, err := RetrieveShopsPlatformsByUserID(uID)

  if err != nil {
    ErrorJSONEncode(w, err, http.StatusUnauthorized)
    return
  }

  js, err := GenerateShopPlatformList(shops)
  if err != nil {
    http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
    return
  }

  w.Header().Set("Content-Type", "application/json; charset=utf-8")
  w.Write(js)
}

func HandleCreateShop(w http.ResponseWriter, r *http.Request, uID string) {
  var s Shop

  if err := decodeJSONBody(w, r, &s.Model); err != nil {
    var mr *malformedRequest
    if errors.As(err, &mr) {
      http.Error(w, mr.msg, mr.status)
    } else {
      http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
    }
    return
  }

  printer := getLocalePrinter(r)
  if err := s.ValidateNew(printer); err != nil {
    ErrorJSONEncode(w, err, http.StatusUnprocessableEntity)
    return
  }

  if err := SaveNewShopToDB(&s, uID); err != nil {
    http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
    return
  }

  js, err := s.JSON()

  if err != nil {
    http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
    return
  }

  w.Header().Set("Content-Type", "application/json; charset=utf-8")
  w.Write(js)

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

