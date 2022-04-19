package laxo

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"time"

	"github.com/mediocregopher/radix/v4"
)

func SetUserSession(u *User) (time.Time, string, error) {
  randomBytes, err := GenerateRandomString(128)

  if err != nil {
    Logger.Error("Couldn't generate random bytes", "error", err)
    return time.Time{}, "", err
  }

  sessionKey := base64.StdEncoding.EncodeToString(randomBytes)

  // Get the seconds till the token expires
  expiresT := time.Now().AddDate(0, 0, AppConfig.AuthCookieExpire)
  expires := time.Until(expiresT)
  expireString := fmt.Sprintf("%.0f", expires.Seconds())

  ctx := context.Background()
  if err := RedisClient.Do(ctx, radix.Cmd(nil, "SETEX", sessionKey, expireString, u.Model.ID)); err != nil {
    Logger.Error("Couldn't set user session in Redis", "error", err)
    return time.Time{}, "", err
  }

  return expiresT, sessionKey, nil
}

func RemoveUserSession(sessionToken string) error {
  ctx := context.Background()
  if err := RedisClient.Do(ctx, radix.Cmd(nil, "DEL", sessionToken)); err != nil {
    Logger.Error("Couldn't remove user session in Redis", "error", err)
    return err
  }
  return nil
}

func SetUserCookie(sessionToken string, w http.ResponseWriter, t time.Time) {
  authCookie := &http.Cookie{
    Name:     AppConfig.AuthCookieName,
    Path:     "/",
    Value:    sessionToken,
    HttpOnly: true,
    Secure:   true,
    Expires:  t,
  }

  http.SetCookie(w, authCookie)
}

func RemoveUserCookie(w http.ResponseWriter) {
  authCookie := &http.Cookie{
    Name:     AppConfig.AuthCookieName,
    Value:    "",
    HttpOnly: true,
    Secure:   true,
    Expires:  time.Unix(0, 0),
  }

  http.SetCookie(w, authCookie)
}

