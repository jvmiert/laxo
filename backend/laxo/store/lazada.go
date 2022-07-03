package store

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4"
	"gopkg.in/guregu/null.v4"
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

func (s *lazadaStore) GetLazadaLaxoLinkByAssetIDAndShopID(assetID string, shopID string) (*sqlc.AssetsLazada, error) {
	row, err := s.queries.GetLazadaLaxoLinkByAssetIDAndShopID(
		context.Background(),
		sqlc.GetLazadaLaxoLinkByAssetIDAndShopIDParams{
			ID:     assetID,
			ShopID: shopID,
		},
	)
	if err != nil {
		return nil, err
	}

	return &sqlc.AssetsLazada{AssetID: row.AssetID.String, LazadaUrl: row.LazadaUrl, Created: row.Created}, err
}

func (s *lazadaStore) CreateLazadaLaxoLink(assetID string, URL string) (*sqlc.AssetsLazada, error) {
	asset, err := s.queries.CreateLazadaLaxoAssetLink(
		context.Background(),
		sqlc.CreateLazadaLaxoAssetLinkParams{
			AssetID:   assetID,
			LazadaUrl: null.StringFrom(URL),
		},
	)

	return &asset, err
}

func (s *lazadaStore) GetLazadaPlatformByShopID(shopID string) (*sqlc.PlatformLazada, error) {
	lazInfo, err := s.queries.GetLazadaPlatformByShopID(
		context.Background(),
		shopID,
	)

	return &lazInfo, err
}

func (s *lazadaStore) UpdateLazadaPlatform(pID string, authResp *lazada.AuthResponse) error {
	return s.queries.UpdateLazadaPlatform(
		context.Background(),
		sqlc.UpdateLazadaPlatformParams{
			AccessToken:      authResp.AccessToken,
			RefreshToken:     authResp.RefreshToken,
			RefreshExpiresIn: null.TimeFrom(authResp.DateRefreshExpired),
			AccessExpiresIn:  null.TimeFrom(authResp.DateAccessExpired),
			ID:               pID,
		},
	)
}

func (s *lazadaStore) SaveNewLazadaPlatform(shopID string, authResp *lazada.AuthResponse) (*sqlc.PlatformLazada, error) {
	pModel, err := s.queries.CreateLazadaPlatform(
		context.Background(),
		sqlc.CreateLazadaPlatformParams{
			ShopID:           shopID,
			AccessToken:      authResp.AccessToken,
			Country:          authResp.Country,
			RefreshToken:     authResp.RefreshToken,
			AccountPlatform:  authResp.AccountPlatform,
			Account:          authResp.Account,
			UserIDVn:         authResp.CountryUserInfo[0].UserID,
			SellerIDVn:       authResp.CountryUserInfo[0].SellerID,
			ShortCodeVn:      authResp.CountryUserInfo[0].ShortCode,
			RefreshExpiresIn: null.TimeFrom(authResp.DateRefreshExpired),
			AccessExpiresIn:  null.TimeFrom(authResp.DateAccessExpired),
		},
	)

	return &pModel, err
}

func (s *lazadaStore) GetValidAccessTokenByShopID(shopID string) (string, error) {
	accessToken, err := s.queries.GetValidAccessTokenByShopID(
		context.Background(),
		shopID,
	)
	if err != nil {
		return "", fmt.Errorf("GetValidAccessTokenByShopID: %w", err)
	}

	return accessToken, nil
}

func (s *lazadaStore) GetValidRefreshTokenByShopID(shopID string) (string, error) {
	refreshToken, err := s.queries.GetValidRefreshTokenByShopID(
		context.Background(),
		shopID,
	)
	if err != nil {
		return "", fmt.Errorf("GetValidRefreshTokenByShopID: %w", err)
	}

	return refreshToken, nil
}

func (s *lazadaStore) SaveOrUpdateLazadaProductSKU(sku []lazada.ProductsResponseSkus, productID, shopID string) (*sqlc.ProductsSkuLazada, error) {
	if len(sku) == 0 {
		return nil, errors.New("product does not have Lazada SKU")
	}

	priceString := sku[0].Price.String()
	specialPriceString := sku[0].SpecialPrice.String()

	priceFloat, err := sku[0].Price.Float64()
	if err == nil {
		priceString = fmt.Sprintf("%.2f", priceFloat)
	}

	skuModel, err := s.queries.GetLazadaProductSKUByProductID(
		context.Background(),
		productID,
	)

	if err != pgx.ErrNoRows && err != nil {
		return nil, fmt.Errorf("GetLazadaProductSKUByProductID: %w", err)
	}

	var pSKUModel sqlc.ProductsSkuLazada

	if skuModel.ID == "" {
		params := sqlc.CreateLazadaProductSKUParams{
			Status:          sku[0].Status,
			Quantity:        sku[0].Quantity,
			SellerSku:       sku[0].SellerSku.String,
			ShopSku:         sku[0].ShopSku.String,
			SkuID:           sku[0].SkuID,
			Url:             sku[0].URL,
			Price:           null.StringFrom(priceString),
			Available:       sku[0].Available,
			PackageContent:  sku[0].PackageContent,
			PackageWidth:    sku[0].PackageWidth,
			PackageWeight:   sku[0].PackageWeight,
			PackageLength:   sku[0].PackageLength,
			PackageHeight:   sku[0].PackageHeight,
			SpecialPrice:    null.StringFrom(specialPriceString),
			SpecialToTime:   null.TimeFrom(sku[0].SpecialToTime),
			SpecialFromTime: null.TimeFrom(sku[0].SpecialFromTime),
			SpecialFromDate: null.TimeFrom(sku[0].SpecialFromDate),
			SpecialToDate:   null.TimeFrom(sku[0].SpecialToDate),
			ProductID:       productID,
			ShopID:          shopID,
		}

		pSKUModel, err = s.queries.CreateLazadaProductSKU(
			context.Background(),
			params,
		)

		if err != nil {
			return nil, fmt.Errorf("CreateLazadaProductSKU: %w", err)
		}

		return &pSKUModel, nil
	}

	params := sqlc.UpdateLazadaProductSKUParams{
		Status:          sku[0].Status,
		Quantity:        sku[0].Quantity,
		SellerSku:       sku[0].SellerSku.String,
		ShopSku:         sku[0].ShopSku.String,
		SkuID:           sku[0].SkuID,
		Url:             sku[0].URL,
		Price:           null.StringFrom(priceString),
		Available:       sku[0].Available,
		PackageContent:  sku[0].PackageContent,
		PackageWidth:    sku[0].PackageWidth,
		PackageWeight:   sku[0].PackageWeight,
		PackageLength:   sku[0].PackageLength,
		PackageHeight:   sku[0].PackageHeight,
		SpecialPrice:    null.StringFrom(specialPriceString),
		SpecialToTime:   null.TimeFrom(sku[0].SpecialToTime),
		SpecialFromTime: null.TimeFrom(sku[0].SpecialFromTime),
		SpecialFromDate: null.TimeFrom(sku[0].SpecialFromDate),
		SpecialToDate:   null.TimeFrom(sku[0].SpecialToDate),
		ProductID:       productID,
		ShopID:          shopID,
		ID:              skuModel.ID,
	}

	pSKUModel, err = s.queries.UpdateLazadaProductSKU(
		context.Background(),
		params,
	)

	if err != nil {
		return nil, fmt.Errorf("UpdateLazadaProductSKU: %w", err)
	}

	return &pSKUModel, nil
}

func (s *lazadaStore) SaveOrUpdateLazadaProductAttribute(a *lazada.ProductsResponseAttributes, productID string) (*sqlc.ProductsAttributeLazada, error) {
	attr, err := s.queries.GetLazadaProductAttributeByProductID(
		context.Background(),
		productID,
	)

	if err != pgx.ErrNoRows && err != nil {
		return nil, fmt.Errorf("GetLazadaProductAttributeByProductID: %w", err)
	}

	var pAttributeModel sqlc.ProductsAttributeLazada

	if attr.ID == "" {
		paramsAttribute := sqlc.CreateLazadaProductAttributeParams{
			Name:                 a.Name,
			ShortDescription:     a.ShortDescription,
			Description:          a.Description,
			Brand:                a.Brand,
			Model:                a.Model,
			HeadphoneFeatures:    a.HeadphoneFeatures,
			Bluetooth:            a.Bluetooth,
			WarrantyType:         a.WarrantyType,
			Warranty:             a.Warranty,
			Hazmat:               a.Hazmat,
			ExpireDate:           a.ExpireDate,
			BrandClassification:  a.BrandClassification,
			IngredientPreference: a.IngredientPreference,
			LotNumber:            a.LotNumber,
			UnitsHb:              a.UnitsHB,
			FmltSkincare:         a.FmltSkinCare,
			Quantitative:         a.Quantitative,
			SkincareByAge:        a.SkinCareByAge,
			SkinBenefit:          a.SkinBenefit,
			SkinType:             a.SkinType,
			UserManual:           a.UserManual,
			CountryOriginHb:      a.CountryOriginHB,
			ColorFamily:          a.ColorFamily,
			FragranceFamily:      a.FragranceFamily,
			Source:               a.Source,
			ProductID:            productID,
		}
		pAttributeModel, err = s.queries.CreateLazadaProductAttribute(
			context.Background(),
			paramsAttribute,
		)

		if err != nil {
			return nil, fmt.Errorf("CreateLazadaProductAttribute: %w", err)
		}

		return &pAttributeModel, nil
	}

	paramsAttribute := sqlc.UpdateLazadaProductAttributeParams{
		Name:                 a.Name,
		ShortDescription:     a.ShortDescription,
		Description:          a.Description,
		Brand:                a.Brand,
		Model:                a.Model,
		HeadphoneFeatures:    a.HeadphoneFeatures,
		Bluetooth:            a.Bluetooth,
		WarrantyType:         a.WarrantyType,
		Warranty:             a.Warranty,
		Hazmat:               a.Hazmat,
		ExpireDate:           a.ExpireDate,
		BrandClassification:  a.BrandClassification,
		IngredientPreference: a.IngredientPreference,
		LotNumber:            a.LotNumber,
		UnitsHb:              a.UnitsHB,
		FmltSkincare:         a.FmltSkinCare,
		Quantitative:         a.Quantitative,
		SkincareByAge:        a.SkinCareByAge,
		SkinBenefit:          a.SkinBenefit,
		SkinType:             a.SkinType,
		UserManual:           a.UserManual,
		CountryOriginHb:      a.CountryOriginHB,
		ColorFamily:          a.ColorFamily,
		FragranceFamily:      a.FragranceFamily,
		Source:               a.Source,
		ID:                   attr.ID,
	}

	pAttributeModel, err = s.queries.UpdateLazadaProductAttribute(
		context.Background(),
		paramsAttribute,
	)

	if err != nil {
		return nil, fmt.Errorf("UpdateLazadaProductAttribute: %w", err)
	}

	return &pAttributeModel, nil
}

func (s *lazadaStore) SaveOrUpdateLazadaProduct(p *lazada.ProductsResponseProducts, shopID string) (*sqlc.ProductsLazada, *sqlc.ProductsAttributeLazada, *sqlc.ProductsSkuLazada, error) {
	qParam := sqlc.GetLazadaProductByLazadaIDParams{
		LazadaID: p.ItemID,
		ShopID:   shopID,
	}

	var pModel sqlc.ProductsLazada
	var pModelAttributes *sqlc.ProductsAttributeLazada
	var pModelSKU *sqlc.ProductsSkuLazada
	var err error

	pModel, err = s.queries.GetLazadaProductByLazadaID(context.Background(), qParam)

	if err != pgx.ErrNoRows && err != nil {
		return nil, nil, nil, fmt.Errorf("GetLazadaProductByLazadaID: %w", err)
	}

	if pModel.ID == "" {
		params := sqlc.CreateLazadaProductParams{
			LazadaID:              p.ItemID,
			LazadaPrimaryCategory: p.PrimaryCategory,
			Created:               p.Created,
			Updated:               time.Now(),
			Status:                p.Status,
			SubStatus:             p.SubStatus,
			ShopID:                shopID,
		}

		pModel, err = s.queries.CreateLazadaProduct(
			context.Background(),
			params,
		)

		if err != nil {
			return nil, nil, nil, fmt.Errorf("CreateLazadaProduct: %w", err)
		}

		pModelAttributes, err = s.SaveOrUpdateLazadaProductAttribute(&p.Attributes, pModel.ID)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("SaveOrUpdateLazadaProductAttribute: %w", err)
		}

		pModelSKU, err = s.SaveOrUpdateLazadaProductSKU(p.Skus, pModel.ID, shopID)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("SaveOrUpdateLazadaProductSKU: %w", err)
		}

		return &pModel, pModelAttributes, pModelSKU, nil
	}

	params := sqlc.UpdateLazadaProductParams{
		LazadaID:              p.ItemID,
		LazadaPrimaryCategory: p.PrimaryCategory,
		Created:               p.Created,
		Updated:               time.Now(),
		Status:                p.Status,
		SubStatus:             p.SubStatus,
		ID:                    pModel.ID,
	}

	pModel, err = s.queries.UpdateLazadaProduct(
		context.Background(),
		params,
	)

	if err != nil {
		return nil, nil, nil, fmt.Errorf("UpdateLazadaProduct: %w", err)
	}

	pModelAttributes, err = s.SaveOrUpdateLazadaProductAttribute(&p.Attributes, pModel.ID)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("SaveOrUpdateLazadaProductAttribute: %w", err)
	}

	pModelSKU, err = s.SaveOrUpdateLazadaProductSKU(p.Skus, pModel.ID, shopID)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("SaveOrUpdateLazadaProductSKU: %w", err)
	}

	return &pModel, pModelAttributes, pModelSKU, nil
}
