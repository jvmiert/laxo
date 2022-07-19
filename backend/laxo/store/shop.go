package store

import (
	"context"
	"fmt"
	"strings"
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

func (s *shopStore) GetProductByProductMSKU(Msku string, shopID string) (*sqlc.Product, error) {
	params := sqlc.GetProductByProductMSKUParams{
		Msku:   null.StringFrom(Msku),
		ShopID: shopID,
	}

	pModel, err := s.queries.GetProductByProductMSKU(
		context.Background(),
		params,
	)

	return &pModel, err
}

func (s *shopStore) UpdateProductImageOrderRequest(productID string, request *models.ProductImageOrderRequest) error {
	var b strings.Builder
	paramCount := 1
	paramList := []interface{}{}

	b.WriteString("UPDATE products_media AS t SET ")
	b.WriteString("image_order = c.image_order ")

	b.WriteString("FROM (VALUES ")

	for i, v := range request.Assets {
		if i+1 == len(request.Assets) {
			fmt.Fprintf(&b, "($%d, $%d::integer) ", paramCount, paramCount+1)
		} else {
			fmt.Fprintf(&b, "($%d, $%d::integer), ", paramCount, paramCount+1)
		}
		paramList = append(paramList, v.AssetID, v.Order)
		paramCount += 2
	}

	b.WriteString(") as c(asset_id, image_order) ")
	b.WriteString("WHERE c.asset_id = t.asset_id AND t.product_id = ")
	fmt.Fprintf(&b, "$%d", paramCount)

	paramList = append(paramList, productID)

	_, err := s.pglClient.Exec(context.Background(), b.String(), paramList...)
	if err != nil {
		return err
	}

	return nil
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

//@TODO: The early returns here are pretty ugly. Should refactor into different functions.
func (s *shopStore) SaveNewProductToStore(p *models.Product, shopID string) (*sqlc.Product, error) {
	var pModel sqlc.Product
	var err error

	// No Laxo product ID is present
	if p.Model.ID != "" {
		// Retrieve Laxo product information from DB
		pModel, err = s.queries.GetProductByID(
			context.Background(),
			p.Model.ID,
		)
		if err != pgx.ErrNoRows && err != nil {
			return nil, fmt.Errorf("GetProductByID: %w", err)
		}
	}

	// No Laxo product information was found, model has a merchant SKU
	if pModel.ID != "" && p.Model.Msku.Valid {
		params := sqlc.GetProductByProductMSKUParams{
			Msku:   p.Model.Msku,
			ShopID: shopID,
		}
		// Attempt to retrieve product information from DB with merchant SKU
		pModel, err = s.queries.GetProductByProductMSKU(
			context.Background(),
			params,
		)
		if err != pgx.ErrNoRows && err != nil {
			return nil, fmt.Errorf("GetProductByProductMSKU: %w", err)
		}
	}

	// No Laxo product information was found after MSKU search
	if pModel.ID == "" {
		params := sqlc.CreateProductParams{
			Name:             p.Model.Name,
			Description:      p.Model.Description,
			Msku:             p.Model.Msku,
			SellingPrice:     p.Model.SellingPrice,
			CostPrice:        p.Model.CostPrice,
			ShopID:           p.Model.ShopID,
			MediaID:          p.Model.MediaID,
			Updated:          null.TimeFrom(time.Now()),
			DescriptionSlate: p.Model.DescriptionSlate,
		}

		// We create a new Laxo product
		pModel, err = s.queries.CreateProduct(
			context.Background(),
			params,
		)
		if err != nil {
			return nil, fmt.Errorf("CreateProduct: %w", err)
		}

		// Return new Laxo product
		return &pModel, nil
	}

	// We either found a valid Laxo product with ID or MSKU
	// above or we already returned after making a new product.
	// Now we're handeling updating an existing product.
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

	// Update the existing Laxo product
	newModel, err := s.queries.UpdateProduct(
		context.Background(),
		params,
	)
	if err != nil {
		return nil, fmt.Errorf("UpdateProduct: %w", err)
	}

	return &newModel, nil
}
