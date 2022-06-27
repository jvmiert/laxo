package store

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4"
	"gopkg.in/guregu/null.v4"
	"laxo.vn/laxo/laxo/models"
	"laxo.vn/laxo/laxo/shop"
	"laxo.vn/laxo/laxo/sqlc"
)

type shopStore struct {
	*Store
}

func newShopStore(store *Store) shopStore {
	return shopStore{
		store,
	}
}

func (s *shopStore) UpdateLazadaProductPlatformSync(productID string, state bool) error {
	_, err := s.queries.UpdateProductsPlatformSync(context.Background(), sqlc.UpdateProductsPlatformSyncParams{
		ProductID:  productID,
		SyncLazada: null.BoolFrom(state),
	})

	return err
}

func (s *shopStore) CheckProductOwner(productID string, shopID string) (string, error) {
	productID, err := s.queries.CheckProductOwner(context.Background(), sqlc.CheckProductOwnerParams{
		ID:     productID,
		ShopID: shopID,
	})

	return productID, err
}

func (s *shopStore) GetShopByID(shopID string) (*sqlc.Shop, error) {
	sModel, err := s.queries.GetShopByID(context.Background(), shopID)
	return &sModel, err
}

func (s *shopStore) SaveNewShopToStore(shop *shop.Shop, u string) (*sqlc.Shop, error) {
	savedShop, err := s.queries.CreateShop(
		context.Background(),
		sqlc.CreateShopParams{
			ShopName:    shop.Model.ShopName,
			UserID:      u,
			AssetsToken: shop.Model.AssetsToken,
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
			ShopID:       shopID,
			PlatformName: platformName,
		},
	)
}

func (s *shopStore) CreateShopsPlatforms(shopID string, platformName string) (sqlc.ShopsPlatform, error) {
	return s.queries.CreatePlatform(
		context.Background(),
		sqlc.CreatePlatformParams{
			ShopID:       shopID,
			PlatformName: platformName,
		},
	)
}

func (s *shopStore) GetProductsByNameOrSKU(shopID string, name null.String, msku null.String, limit int32, offset int32) ([]sqlc.GetProductsByNameOrSKURow, error) {
	products, err := s.queries.GetProductsByNameOrSKU(
		context.Background(),
		sqlc.GetProductsByNameOrSKUParams{
			ShopID: shopID,
			Name:   name,
			Msku:   msku,
			Limit:  limit,
			Offset: offset,
		},
	)
	if err != nil {
		return nil, err
	}

	return products, nil
}

func (s *shopStore) GetProductsByShopID(shopID string, limit int32, offset int32) ([]sqlc.GetProductsByShopIDRow, error) {
	products, err := s.queries.GetProductsByShopID(
		context.Background(),
		sqlc.GetProductsByShopIDParams{
			ShopID: shopID,
			Limit:  limit,
			Offset: offset,
		},
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

func (s *shopStore) UpdateProductToStore(p *models.Product) (*sqlc.Product, error) {
	params := sqlc.UpdateProductParams{
		Name:             p.Model.Name,
		Description:      p.Model.Description,
		DescriptionSlate: p.Model.DescriptionSlate,
		SellingPrice:     p.Model.SellingPrice,
		CostPrice:        p.Model.CostPrice,
		ShopID:           p.Model.ShopID,
		MediaID:          p.Model.MediaID,
		Updated:          null.TimeFrom(time.Now()),
		ID:               p.Model.ID,
	}
	newModel, err := s.queries.UpdateProduct(
		context.Background(),
		params,
	)

	if err != nil {
		return nil, fmt.Errorf("UpdateProduct: %w", err)
	}

	return &newModel, nil
}

func (s *shopStore) GetProductAssetsByProductID(productID string, shopID string) ([]sqlc.GetProductAssetsByProductIDRow, error) {
	a, err := s.queries.GetProductAssetsByProductID(
		context.Background(),
		sqlc.GetProductAssetsByProductIDParams{
			ID:     productID,
			ShopID: shopID,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("GetProductAssetsByProductID: %w", err)
	}

	return a, nil
}

func (s *shopStore) GetProductDetails(productID string, shopID string) (*sqlc.GetProductDetailsByIDRow, error) {
	p, err := s.queries.GetProductDetailsByID(context.Background(), sqlc.GetProductDetailsByIDParams{ID: productID, ShopID: shopID})
	if err != nil {
		return nil, fmt.Errorf("GetProductDetailsByID: %w", err)
	}

	return &p, nil
}

func (s *shopStore) GetProductByID(productID string) (*sqlc.Product, error) {
	p, err := s.queries.GetProductByID(
		context.Background(),
		productID,
	)

	if err != nil {
		return nil, fmt.Errorf("GetProductByID: %w", err)
	}

	return &p, nil
}

func (s *shopStore) SaveNewProductToStore(p *models.Product, shopID string) (*sqlc.Product, error) {
	var pModel sqlc.Product
	var err error

	if p.Model.ID != "" {
		pModel, err = s.queries.GetProductByID(
			context.Background(),
			p.Model.ID,
		)

		if err != pgx.ErrNoRows && err != nil {
			return nil, fmt.Errorf("GetProductByID: %w", err)
		}
	}

	if pModel.ID != "" && p.Model.Msku.Valid {
		params := sqlc.GetProductByProductMSKUParams{
			Msku:   p.Model.Msku,
			ShopID: shopID,
		}
		pModel, err = s.queries.GetProductByProductMSKU(
			context.Background(),
			params,
		)

		if err != pgx.ErrNoRows && err != nil {
			return nil, fmt.Errorf("GetProductByProductMSKU: %w", err)
		}
	}

	if pModel.ID == "" {
		params := sqlc.CreateProductParams{
			Name:         p.Model.Name,
			Description:  p.Model.Description,
			Msku:         p.Model.Msku,
			SellingPrice: p.Model.SellingPrice,
			CostPrice:    p.Model.CostPrice,
			ShopID:       p.Model.ShopID,
			MediaID:      p.Model.MediaID,
			Updated:      null.TimeFrom(time.Now()),
		}
		pModel, err = s.queries.CreateProduct(
			context.Background(),
			params,
		)
		if err != nil {
			return nil, fmt.Errorf("CreateProduct: %w", err)
		}

		return &pModel, nil
	}

	params := sqlc.UpdateProductParams{
		Name:         p.Model.Name,
		Description:  p.Model.Description,
		SellingPrice: p.Model.SellingPrice,
		CostPrice:    p.Model.CostPrice,
		ShopID:       p.Model.ShopID,
		MediaID:      p.Model.MediaID,
		Updated:      null.TimeFrom(time.Now()),
		ID:           pModel.ID,
	}

	newModel, err := s.queries.UpdateProduct(
		context.Background(),
		params,
	)
	if err != nil {
		return nil, fmt.Errorf("UpdateProduct: %w", err)
	}

	return &newModel, nil
}
