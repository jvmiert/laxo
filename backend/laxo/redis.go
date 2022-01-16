package laxo

import (
  "context"
  "github.com/mediocregopher/radix/v4"
)

var RedisClient radix.Client

func InitRedis(uri string) error {
  Logger.Debug("Connecting to Redis", "uri", uri)
  client, err := (radix.PoolConfig{}).New(context.Background(), "tcp", uri)
  if err != nil {
    return err
  }

  RedisClient = client
  return nil
}
