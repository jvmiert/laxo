package laxo

import (
	"context"

	"github.com/mediocregopher/radix/v4"
)

func (s *Server) InitRedis(uri string) error {
  s.Logger.Infow("Connecting to Redis...",
    "uri", uri,
  )

  client, err := (radix.PoolConfig{
    Size: 10,
  }).New(context.Background(), "tcp", uri)
  if err != nil {
    return err
  }

  s.RedisClient = client
  return nil
}
