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
      Name: a.Name.NullString,
      ShortDescription: a.ShortDescription.NullString,
      Description: a.Description.NullString,
      Brand: a.Brand.NullString,
      Model: a.Model.NullString,
      HeadphoneFeatures: a.HeadphoneFeatures.NullString,
      Bluetooth: a.Bluetooth.NullString,
      WarrantyType: a.WarrantyType.NullString,
      Warranty: a.Warranty.NullString,
      Hazmat: a.Hazmat.NullString,
      ExpireDate: a.ExpireDate.NullString,
      BrandClassification: a.BrandClassification.NullString,
      IngredientPreference: a.IngredientPreference.NullString,
      LotNumber: a.LotNumber.NullString,
      UnitsHb: a.UnitsHB.NullString,
      FmltSkincare: a.FmltSkinCare.NullString,
      Quantitative: a.Quantitative.NullString,
      SkincareByAge: a.SkinCareByAge.NullString,
      SkinBenefit: a.SkinBenefit.NullString,
      SkinType: a.SkinType.NullString,
      UserManual: a.UserManual.NullString,
      CountryOriginHb: a.CountryOriginHB.NullString,
      ColorFamily: a.ColorFamily.NullString,
      FragranceFamily: a.FragranceFamily.NullString,
      Source: a.Source.NullString,
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
    Name: a.Name.NullString,
    ShortDescription: a.ShortDescription.NullString,
    Description: a.Description.NullString,
    Brand: a.Brand.NullString,
    Model: a.Model.NullString,
    HeadphoneFeatures: a.HeadphoneFeatures.NullString,
    Bluetooth: a.Bluetooth.NullString,
    WarrantyType: a.WarrantyType.NullString,
    Warranty: a.Warranty.NullString,
    Hazmat: a.Hazmat.NullString,
    ExpireDate: a.ExpireDate.NullString,
    BrandClassification: a.BrandClassification.NullString,
    IngredientPreference: a.IngredientPreference.NullString,
    LotNumber: a.LotNumber.NullString,
    UnitsHb: a.UnitsHB.NullString,
    FmltSkincare: a.FmltSkinCare.NullString,
    Quantitative: a.Quantitative.NullString,
    SkincareByAge: a.SkinCareByAge.NullString,
    SkinBenefit: a.SkinBenefit.NullString,
    SkinType: a.SkinType.NullString,
    UserManual: a.UserManual.NullString,
    CountryOriginHb: a.CountryOriginHB.NullString,
    ColorFamily: a.ColorFamily.NullString,
    FragranceFamily: a.FragranceFamily.NullString,
    Source: a.Source.NullString,
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
      Status: p.Status.NullString,
      SubStatus: p.SubStatus.NullString,
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
    Status: p.Status.NullString,
    SubStatus: p.SubStatus.NullString,
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
