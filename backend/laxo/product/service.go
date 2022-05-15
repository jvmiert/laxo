package product

import (
	"github.com/hashicorp/go-hclog"
	"github.com/jackc/pgx/v4"
	"github.com/mediocregopher/radix/v4"
	"gopkg.in/guregu/null.v4"
	"laxo.vn/laxo/laxo/sqlc"
)

type Store interface {
  SaveNewProductToStore(*Product, string) (*sqlc.Product, error)
  GetProductPlatformByProductID(string) (*sqlc.ProductsPlatform, error)
  GetProductPlatformByLazadaID(string) (*sqlc.ProductsPlatform, error)
  CreateProductPlatform(*sqlc.CreateProductPlatformParams) (*sqlc.ProductsPlatform, error)
  UpdateProductToStore(*Product) (*sqlc.Product, error)
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

func (s *Service) GetProductPlatformByProductID(productID string) (*sqlc.ProductsPlatform, error) {
  return s.store.GetProductPlatformByProductID(productID)
}

func (s *Service) GetProductPlatformByLazadaID(productID string) (*sqlc.ProductsPlatform, error) {
  return s.store.GetProductPlatformByLazadaID(productID)
}

func (s *Service) SaveOrUpdateProductToStore(p *Product, shopID string, lazadaID string) (*Product, error) {
  var platform *sqlc.ProductsPlatform
  var pReturn *Product
  var newModel *sqlc.Product
  var err error

  platform, err = s.GetProductPlatformByLazadaID(lazadaID)
  if err != pgx.ErrNoRows && err != nil {
    return nil, err
  }

  // product was not yet saved
  if err == pgx.ErrNoRows {
    newModel, err = s.store.SaveNewProductToStore(p, shopID)
    if err != nil {
      return nil, err
    }

    param := &sqlc.CreateProductPlatformParams{
      ProductID: newModel.ID,
      ProductsLazadaID: null.StringFrom(lazadaID),
    }
    platform, err = s.store.CreateProductPlatform(param)
    if err != nil {
      return nil, err
    }

    pReturn = &Product{
      Model: newModel,
      PlatformModel: platform,
    }

    return pReturn, nil
  }

  p.Model.ID = platform.ProductID

  newModel, err = s.store.UpdateProductToStore(p)
  if err != nil {
    return nil, err
  }

  pReturn = &Product{
    Model: newModel,
    PlatformModel: platform,
  }

  return pReturn, nil
}
