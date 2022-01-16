package laxo

import (
  "fmt"
  "time"
  "context"
  "encoding/base64"

  "github.com/mediocregopher/radix/v4"
)

func SetUserSession(u *User) (string, error) {
  randomBytes, err := GenerateRandomString(128)

  if err != nil {
    Logger.Error("Couldn't generate random bytes", "error", err)
    return "", err
  }

  sessionKey := base64.StdEncoding.EncodeToString(randomBytes)

  // Get the seconds till the token expires
  expires := time.Until(time.Now().AddDate(0, 0, AppConfig.AuthCookieExpire))
  expireString := fmt.Sprintf("%.0f", expires.Seconds())

  ctx := context.Background()
  if err := RedisClient.Do(ctx, radix.Cmd(nil, "SETEX", sessionKey, expireString, u.Model.ID)); err != nil {
    Logger.Error("Couldn't set user session in Redis", "error", err)
    return "", err
  }

  return sessionKey, nil
}
