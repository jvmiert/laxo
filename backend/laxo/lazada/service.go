package lazada

import (
	"github.com/hashicorp/go-hclog"
)

type Store interface {
  SaveOrUpdateProduct(ProductsResponseProducts, string) error
}

type Service struct {
  store Store
  logger hclog.Logger
}

func NewService(store Store, logger hclog.Logger) Service {
  return Service {
    store: store,
    logger: logger,
  }
}

func (s *Service) SaveOrUpdateProduct(p ProductsResponseProducts, shopID string) error {
  if err := s.store.SaveOrUpdateProduct(p, shopID); err != nil {
    return err
  }

  return nil
}
