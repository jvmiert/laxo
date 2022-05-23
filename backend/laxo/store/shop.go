package store

import (
	"context"
	"time"

	"github.com/jackc/pgx/v4"
	"gopkg.in/guregu/null.v4"
	"laxo.vn/laxo/laxo/shop"
	"laxo.vn/laxo/laxo/sqlc"
)

type shopStore struct {
  *Store
}

func newShopStore(store *Store) shopStore{
  return shopStore{
    store,
  }
}

func (s *shopStore) GetShopByID(shopID string) (*sqlc.Shop, error) {
  sModel, err := s.queries.GetShopByID(context.Background(), shopID)
  return &sModel, err
}

func (s *shopStore) GetLazadaPlatformByShopID(shopID string) (*sqlc.PlatformLazada, error) {
  lazInfo, err := s.queries.GetLazadaPlatformByShopID(
    context.Background(),
    shopID,
  )

  return &lazInfo, err
}

func (s *shopStore) SaveNewShopToStore(shop *shop.Shop, u string) (*sqlc.Shop, error) {
  savedShop, err := s.queries.CreateShop(
    context.Background(),
    sqlc.CreateShopParams{
      ShopName: shop.Model.ShopName,
      UserID: u,
    },
  )

  return &savedShop, err
}

func (s *shopStore) RetrieveShopsPlatformsByUserID(userID string) ([]sqlc.GetShopsPlatformsByUserIDRow, error) {
  return s.queries.GetShopsPlatformsByUserID(
    context.Background(),
    userID,
  )
}

func (s *shopStore) RetrieveShopsPlatformsByShopID(shopID string) ([]sqlc.ShopsPlatform, error) {
  return s.queries.GetPlatformsByShopID(
    context.Background(),
    shopID,
  )
}

func (s *shopStore) RetrieveSpecificPlatformByShopID(shopID string, platformName string) (sqlc.ShopsPlatform, error) {
  return s.queries.GetSpecificPlatformByShopID(
    context.Background(),
    sqlc.GetSpecificPlatformByShopIDParams{
      ShopID: shopID,
      PlatformName: platformName,
    },
  )
}

func (s *shopStore) CreateShopsPlatforms(shopID string, platformName string) (sqlc.ShopsPlatform, error) {
  return s.queries.CreatePlatform(
    context.Background(),
    sqlc.CreatePlatformParams{
      ShopID: shopID,
      PlatformName: platformName,
    },
  )
}

func (s *shopStore) GetProductsByShopID(shopID string) ([]sqlc.GetProductsByShopIDRow, error) {
  products, err := s.queries.GetProductsByShopID(
    context.Background(),
    shopID,
  )
  if err != nil {
    return nil, err
  }

  return products, nil
}

func (s *shopStore) GetProductPlatformByProductID(productID string) (*sqlc.ProductsPlatform, error) {
  pPlatform, err := s.queries.GetProductPlatformByProductID(
    context.Background(),
    productID,
  )
  if err != nil {
    return nil, err
  }

  return &pPlatform, nil
}


func (s *shopStore) GetProductPlatformByLazadaID(lazadaID string) (*sqlc.ProductsPlatform, error) {
  pPlatform, err := s.queries.GetProductPlatformByLazadaID(
    context.Background(),
    null.StringFrom(lazadaID),
  )
  if err != nil {
    return nil, err
  }

  return &pPlatform, nil
}

func (s *shopStore) CreateProductPlatform(param *sqlc.CreateProductPlatformParams) (*sqlc.ProductsPlatform, error) {
  pPlatform, err := s.queries.CreateProductPlatform(
    context.Background(),
    *param,
  )

  return &pPlatform, err
}

func (s *shopStore) UpdateProductToStore(p *shop.Product) (*sqlc.Product, error) {
  params := sqlc.UpdateProductParams{
    Name: p.Model.Name,
    Description: p.Model.Description,
    Msku: p.Model.Msku,
    SellingPrice: p.Model.SellingPrice,
    CostPrice: p.Model.CostPrice,
    ShopID: null.StringFrom(p.Model.ShopID),
    MediaID: p.Model.MediaID,
    Updated: null.TimeFrom(time.Now()),
    ID: p.Model.ID,
  }
  newModel, err := s.queries.UpdateProduct(
    context.Background(),
    params,
  )

  return &newModel, err
}

func (s *shopStore) SaveNewProductToStore(p *shop.Product, shopID string) (*sqlc.Product, error) {
  var pModel sqlc.Product
  var err error

  if p.Model.ID != "" {
    pModel, err = s.queries.GetProductByID(
      context.Background(),
      p.Model.ID,
    )

    if err != pgx.ErrNoRows && err != nil {
      return nil, err
    }
  }

  if pModel.ID != "" && p.Model.Msku.Valid {
    params := sqlc.GetProductByProductMSKUParams{
      Msku: p.Model.Msku,
      ShopID: shopID,
    }
    pModel, err = s.queries.GetProductByProductMSKU(
      context.Background(),
      params,
    )

    if err != pgx.ErrNoRows && err != nil {
      return nil, err
    }
  }

  if pModel.ID == "" {
    params := sqlc.CreateProductParams{
      Name: p.Model.Name,
      Description: p.Model.Description,
      Msku: p.Model.Description,
      SellingPrice: p.Model.SellingPrice,
      CostPrice: p.Model.CostPrice,
      ShopID: p.Model.ShopID,
      MediaID: p.Model.MediaID,
      Updated: null.TimeFrom(time.Now()),
    }
    pModel, err = s.queries.CreateProduct(
      context.Background(),
      params,
    )
    if err != nil {
      return nil, err
    }

    return &pModel, nil
  }

  params := sqlc.UpdateProductParams{
    Name: p.Model.Name,
    Description: p.Model.Description,
    Msku: p.Model.Description,
    SellingPrice: p.Model.SellingPrice,
    CostPrice: p.Model.CostPrice,
    ShopID: null.StringFrom(p.Model.ShopID),
    MediaID: p.Model.MediaID,
    Updated: null.TimeFrom(time.Now()),
    ID: pModel.ID,
  }

  newModel, err := s.queries.UpdateProduct(
    context.Background(),
    params,
  )
  if err != nil {
    return nil, err
  }

  return &newModel, nil
}
