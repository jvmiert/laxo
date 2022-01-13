package laxo

import (
  "context"
  "github.com/mediocregopher/radix/v4"
  "encoding/base64"
)

func SetUserSession(u *User) (string, error) {
  randomBytes, err := GenerateRandomString(128)

  if err != nil {
    Logger.Error("Couldn't generate random bytes", "error", err)
    return "", err
  }

  sessionKey := base64.StdEncoding.EncodeToString(randomBytes)

  ctx := context.Background()
  if err := RedisClient.Do(ctx, radix.Cmd(nil, "SET", sessionKey, u.Model.ID)); err != nil {
    Logger.Error("Couldn't set user session in Redis", "error", err)
    return "", err
  }

  return sessionKey, nil
}
