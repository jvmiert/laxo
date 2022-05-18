package store

import (
	"context"
	"time"

	"github.com/jackc/pgx/v4"
	"gopkg.in/guregu/null.v4"
	"laxo.vn/laxo/laxo/product"
	"laxo.vn/laxo/laxo/sqlc"
)

type productStore struct {
  *Store
}

func newProductStore(store *Store) productStore{
  return productStore{
    store,
  }
}

func (s *productStore) GetProductsByShopID(shopID string) ([]sqlc.Product, error) {
  products, err := s.queries.GetProductsByShopID(
    context.Background(),
    shopID,
  )
  if err != nil {
    return nil, err
  }

  return products, nil
}

func (s *productStore) GetProductPlatformByProductID(productID string) (*sqlc.ProductsPlatform, error) {
  pPlatform, err := s.queries.GetProductPlatformByProductID(
    context.Background(),
    productID,
  )
  if err != nil {
    return nil, err
  }

  return &pPlatform, nil
}


func (s *productStore) GetProductPlatformByLazadaID(lazadaID string) (*sqlc.ProductsPlatform, error) {
  pPlatform, err := s.queries.GetProductPlatformByLazadaID(
    context.Background(),
    null.StringFrom(lazadaID),
  )
  if err != nil {
    return nil, err
  }

  return &pPlatform, nil
}

func (s *productStore) CreateProductPlatform(param *sqlc.CreateProductPlatformParams) (*sqlc.ProductsPlatform, error) {
  pPlatform, err := s.queries.CreateProductPlatform(
    context.Background(),
    *param,
  )

  return &pPlatform, err
}

func (s *productStore) UpdateProductToStore(p *product.Product) (*sqlc.Product, error) {
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

func (s *productStore) SaveNewProductToStore(p *product.Product, shopID string) (*sqlc.Product, error) {
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
      s.logger.Error("CreateProduct error", "params", params)
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
