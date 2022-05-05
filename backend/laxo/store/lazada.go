package store

import (
	"context"

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


func (s *lazadaStore) SaveOrUpdateProductAttirube(a *lazada.ProductsResponseAttributes, productID string) error {
  attr, err := s.queries.GetLazadaProductAttributeByProductID(
    context.Background(),
    productID,
  )

  if err != pgx.ErrNoRows && err != nil {
    return err
  }


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
    _, err = s.queries.CreateLazadaProductAttribute(
      context.Background(),
      paramsAttribute,
    )

    if err != nil {
      return err
    }

    return nil
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

  err = s.queries.UpdateLazadaProductAttribute(
    context.Background(),
    paramsAttribute,
  )

  if err != nil {
    return err
  }

  return nil
}

func (s *lazadaStore) SaveOrUpdateProduct(p lazada.ProductsResponseProducts, shopID string) error {
  qParam := sqlc.GetLazadaProductByLazadaIDParams{
    LazadaID: p.ItemID,
    ShopID: shopID,
  }

  pModel, err := s.queries.GetLazadaProductByLazadaID(context.Background(), qParam)

  if err != pgx.ErrNoRows && err != nil {
    return err
  }

  if pModel.ID == "" {
    params := sqlc.CreateLazadaProductParams{
      LazadaID: p.ItemID,
      LazadaPrimaryCategory: p.PrimaryCategory,
      Created: p.Created,
      Updated: p.Updated,
      Status: p.Status,
      SubStatus: p.SubStatus,
      ShopID: shopID,
    }

    pModel, err = s.queries.CreateLazadaProduct(
      context.Background(),
      params,
    )

    if err != nil {
      return err
    }

    if err = s.SaveOrUpdateProductAttirube(&p.Attributes, pModel.ID); err != nil {
      return err
    }

    return nil
  }

  params := sqlc.UpdateLazadaProductParams {
    LazadaID: p.ItemID,
    LazadaPrimaryCategory: p.PrimaryCategory,
    Created: p.Created,
    Updated: p.Updated,
    Status: p.Status,
    SubStatus: p.SubStatus,
    ID: pModel.ID,
  }

  err = s.queries.UpdateLazadaProduct(
    context.Background(),
    params,
  )

  if err != nil {
    return err
  }

  if err = s.SaveOrUpdateProductAttirube(&p.Attributes, pModel.ID); err != nil {
    return err
  }

  return nil
}
