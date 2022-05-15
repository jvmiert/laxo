package store

import (
	"context"
	"time"

	"github.com/jackc/pgx/v4"
	"laxo.vn/laxo/laxo/lazada"
	"laxo.vn/laxo/laxo/sqlc"
)

type lazadaStore struct {
  *Store
}

func newLazadaStore(store *Store) lazadaStore {
  return lazadaStore{
    store,
  }
}

func (s *lazadaStore) GetValidTokenByShopID(shopID string) (string, error) {
  accessToken, err := s.queries.GetValidAccessTokenByShopID(
    context.Background(),
    shopID,
  )
  if err != nil {
    return "", err
  }

  return accessToken, nil
}

func (s *lazadaStore) SaveOrUpdateLazadaProductAttribute(a *lazada.ProductsResponseAttributes, productID string) (*sqlc.ProductsAttributeLazada, error) {
  attr, err := s.queries.GetLazadaProductAttributeByProductID(
    context.Background(),
    productID,
  )

  if err != pgx.ErrNoRows && err != nil {
    return nil, err
  }

  var pAttributeModel sqlc.ProductsAttributeLazada

  if attr.ID == "" {
    paramsAttribute := sqlc.CreateLazadaProductAttributeParams{
      Name: a.Name,
      ShortDescription: a.ShortDescription,
      Description: a.Description,
      Brand: a.Brand,
      Model: a.Model,
      HeadphoneFeatures: a.HeadphoneFeatures,
      Bluetooth: a.Bluetooth,
      WarrantyType: a.WarrantyType,
      Warranty: a.Warranty,
      Hazmat: a.Hazmat,
      ExpireDate: a.ExpireDate,
      BrandClassification: a.BrandClassification,
      IngredientPreference: a.IngredientPreference,
      LotNumber: a.LotNumber,
      UnitsHb: a.UnitsHB,
      FmltSkincare: a.FmltSkinCare,
      Quantitative: a.Quantitative,
      SkincareByAge: a.SkinCareByAge,
      SkinBenefit: a.SkinBenefit,
      SkinType: a.SkinType,
      UserManual: a.UserManual,
      CountryOriginHb: a.CountryOriginHB,
      ColorFamily: a.ColorFamily,
      FragranceFamily: a.FragranceFamily,
      Source: a.Source,
      ProductID: productID,
    }
    pAttributeModel, err = s.queries.CreateLazadaProductAttribute(
      context.Background(),
      paramsAttribute,
    )

    if err != nil {
      return nil, err
    }

    return &pAttributeModel, nil
  }

  paramsAttribute := sqlc.UpdateLazadaProductAttributeParams{
    Name: a.Name,
    ShortDescription: a.ShortDescription,
    Description: a.Description,
    Brand: a.Brand,
    Model: a.Model,
    HeadphoneFeatures: a.HeadphoneFeatures,
    Bluetooth: a.Bluetooth,
    WarrantyType: a.WarrantyType,
    Warranty: a.Warranty,
    Hazmat: a.Hazmat,
    ExpireDate: a.ExpireDate,
    BrandClassification: a.BrandClassification,
    IngredientPreference: a.IngredientPreference,
    LotNumber: a.LotNumber,
    UnitsHb: a.UnitsHB,
    FmltSkincare: a.FmltSkinCare,
    Quantitative: a.Quantitative,
    SkincareByAge: a.SkinCareByAge,
    SkinBenefit: a.SkinBenefit,
    SkinType: a.SkinType,
    UserManual: a.UserManual,
    CountryOriginHb: a.CountryOriginHB,
    ColorFamily: a.ColorFamily,
    FragranceFamily: a.FragranceFamily,
    Source: a.Source,
    ID: attr.ID,
  }

  pAttributeModel, err = s.queries.UpdateLazadaProductAttribute(
    context.Background(),
    paramsAttribute,
  )

  if err != nil {
    return nil, err
  }

  return &pAttributeModel, nil
}

func (s *lazadaStore) SaveOrUpdateLazadaProduct(p *lazada.ProductsResponseProducts, shopID string) (*sqlc.ProductsLazada, *sqlc.ProductsAttributeLazada, error) {
  qParam := sqlc.GetLazadaProductByLazadaIDParams{
    LazadaID: p.ItemID,
    ShopID: shopID,
  }

  var pModel sqlc.ProductsLazada
  var pModelAttributes *sqlc.ProductsAttributeLazada
  var err error

  pModel, err = s.queries.GetLazadaProductByLazadaID(context.Background(), qParam)

  if err != pgx.ErrNoRows && err != nil {
    return nil, nil, err
  }

  if pModel.ID == "" {
    params := sqlc.CreateLazadaProductParams{
      LazadaID: p.ItemID,
      LazadaPrimaryCategory: p.PrimaryCategory,
      Created: p.Created,
      Updated: time.Now(),
      Status: p.Status,
      SubStatus: p.SubStatus,
      ShopID: shopID,
    }

    pModel, err = s.queries.CreateLazadaProduct(
      context.Background(),
      params,
    )

    if err != nil {
      return nil, nil, err
    }

    pModelAttributes, err = s.SaveOrUpdateLazadaProductAttribute(&p.Attributes, pModel.ID)
    if err != nil {
      return nil, nil, err
    }

    return &pModel, pModelAttributes, nil
  }

  params := sqlc.UpdateLazadaProductParams {
    LazadaID: p.ItemID,
    LazadaPrimaryCategory: p.PrimaryCategory,
    Created: p.Created,
    Updated: time.Now(),
    Status: p.Status,
    SubStatus: p.SubStatus,
    ID: pModel.ID,
  }

  pModel, err = s.queries.UpdateLazadaProduct(
    context.Background(),
    params,
  )

  if err != nil {
    return nil, nil, err
  }

  pModelAttributes, err = s.SaveOrUpdateLazadaProductAttribute(&p.Attributes, pModel.ID)
  if err != nil {
    return nil, nil, err
  }

  return &pModel, pModelAttributes, nil
}
