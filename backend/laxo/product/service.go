package product

import (
	"encoding/json"
	"errors"

	"github.com/hashicorp/go-hclog"
	"github.com/jackc/pgx/v4"
	"github.com/mediocregopher/radix/v4"
	"gopkg.in/guregu/null.v4"
	"laxo.vn/laxo/laxo"
	"laxo.vn/laxo/laxo/sqlc"
)

type Store interface {
  SaveNewProductToStore(*Product, string) (*sqlc.Product, error)
  GetProductPlatformByProductID(string) (*sqlc.ProductsPlatform, error)
  GetProductPlatformByLazadaID(string) (*sqlc.ProductsPlatform, error)
  CreateProductPlatform(*sqlc.CreateProductPlatformParams) (*sqlc.ProductsPlatform, error)
  UpdateProductToStore(*Product) (*sqlc.Product, error)
  RetrieveShopsByUserID(string) ([]laxo.Shop, error)
  GetProductsByShopID(string) ([]sqlc.Product, error)

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


func (s *Service) GetProductListJSON(pp []Product) ([]byte, error) {
  pList := []json.RawMessage{}

  for _, p := range pp {
    b, err := p.JSON()
    if err != nil {
      return nil, err
    }
    j := json.RawMessage(b)
    pList = append(pList, j)
  }

	productData := map[string]interface{}{
		"products": pList,
		"total": len(pp),
	}

  bytes, err := json.Marshal(productData)
  if err != nil {
    return bytes, err
  }

  return bytes, nil
}

func (s *Service) GetProductsByUserID(userID string) ([]Product, error) {
  var pList []Product

  shops, err := s.store.RetrieveShopsByUserID(userID)
  if err != nil {
    return pList, err
  }

  if len(shops) == 0 {
    return pList, errors.New("user has not setup any shops yet")
  }

  //@TODO: we don't have an active store logic yet so for now we pick the first
  shopID := shops[0].Model.ID

  pModelList, err := s.store.GetProductsByShopID(shopID)
  if err != nil {
    return pList, err
  }

  for _, pModel := range pModelList {
    pList = append(pList, Product{
      Model: &sqlc.Product{
        ID: pModel.ID,
        Name: pModel.Name,
        Description: pModel.Description,
        Msku: pModel.Msku,
        SellingPrice: pModel.SellingPrice,
        CostPrice: pModel.CostPrice,
        ShopID: pModel.ShopID,
        MediaID: pModel.MediaID,
        Created: pModel.Created,
        Updated: pModel.Updated,
      },
    })
  }

  return pList, nil
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
