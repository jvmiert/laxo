package product

import (
	"github.com/hashicorp/go-hclog"
	"github.com/mediocregopher/radix/v4"
)

type Store interface {
}

type Service struct {
  store       Store
  logger      hclog.Logger
  redisClient radix.Client
}

func NewService(store Store, logger hclog.Logger, redisClient radix.Client) Service {
  return Service {
    store: store,
    logger: logger,
    redisClient: redisClient,
  }
}
