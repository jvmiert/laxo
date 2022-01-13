package laxo

import (
  "context"
  "github.com/mediocregopher/radix/v4"
)

var RedisClient radix.Client

func InitRedis() error {
  Logger.Debug("Connecting to Redis")
  client, err := (radix.PoolConfig{}).New(context.Background(), "tcp", "127.0.0.1:6379")
  if err != nil {
    return err
  }

  RedisClient = client
  return nil
}
