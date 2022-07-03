package store

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/guregu/null.v4"
	"laxo.vn/laxo/laxo/sqlc"
)

var ErrInvalidPath = errors.New("set a valid assets base path")

const productImageDir = "products"

type assetsStore struct {
	assetsPath string
	*Store
}

func newAssetsStore(store *Store, assetsBasePath string) (*assetsStore, error) {
	if assetsBasePath == "" {
		return nil, ErrInvalidPath
	}

	dir := filepath.Dir(assetsBasePath)

	fi, err := os.Stat(dir)
	if os.IsNotExist(err) {
		return nil, ErrInvalidPath
	} else if err != nil {
		return nil, err
	}

	if !fi.IsDir() {
		return nil, ErrInvalidPath
	}

	err = os.MkdirAll(assetsBasePath, 0664)
	if err != nil {
		return nil, err
	}

	return &assetsStore{
		dir,
		store,
	}, nil
}

func (s *assetsStore) CreateNewAsset(ShopID, MurmurHash, OriginalFilename string, FileSize, Width, Height int64) (*sqlc.Asset, error) {
	param := sqlc.CreateAssetParams{
		ShopID:           ShopID,
		MurmurHash:       MurmurHash,
		OriginalFilename: null.StringFrom(OriginalFilename),
		Extension:        null.NewString("", false),
		FileSize:         null.IntFrom(FileSize),
		Width:            null.IntFrom(Width),
		Height:           null.IntFrom(Height),
	}

	asset, err := s.queries.CreateAsset(
		context.Background(),
		param,
	)

	return &asset, err
}

func (s *assetsStore) GetAssetBytesByID(assetID string, shopID string, shopToken string) ([]byte, error) {
	assetModel, err := s.GetAssetByIDAndShopID(assetID, shopID)
	if err != nil {
		return nil, fmt.Errorf("GetAssetByIDAndShopID: %w", err)
	}

	if !assetModel.Extension.Valid {
		return nil, fmt.Errorf("!Extension.Valid: %w", errors.New("invalid image extension"))
	}

	assetDirPath := s.GetAssetDirPath(shopToken)
	filename := assetModel.ID + assetModel.Extension.String

	b, err := os.ReadFile(assetDirPath + filename)
	if err != nil {
		return nil, fmt.Errorf("ReadFile: %w", err)
	}

	return b, err
}

func (s *assetsStore) UnlinkProductMedia(productID string, assetID string) error {
	param := sqlc.DeleteProductMediaParams{
		ProductID: productID,
		AssetID:   assetID,
	}

	return s.queries.DeleteProductMedia(
		context.Background(),
		param,
	)
}

func (s *assetsStore) UpdateProductMedia(assetID string, productID string, status null.String, order null.Int) (*sqlc.ProductsMedia, error) {
	media, err := s.queries.UpdateProductMedia(
		context.Background(),
		sqlc.UpdateProductMediaParams{
			ProductID:  productID,
			AssetID:    assetID,
			Status:     status,
			ImageOrder: order,
		},
	)
	if err != nil {
		return nil, err
	}

	return &media, nil
}

func (s *assetsStore) CreateProductMedia(assetID string, productID string, status string, order int64) (*sqlc.ProductsMedia, error) {
	media, err := s.queries.CreateProductMedia(
		context.Background(),
		sqlc.CreateProductMediaParams{
			ProductID:  productID,
			AssetID:    assetID,
			Status:     status,
			ImageOrder: null.IntFrom(order),
		},
	)
	if err != nil {
		return nil, err
	}

	return &media, nil
}

func (s *assetsStore) GetProductMedia(assetID string, productID string) (*sqlc.ProductsMedia, error) {
	media, err := s.queries.GetProductMedia(
		context.Background(),
		sqlc.GetProductMediaParams{
			ProductID: productID,
			AssetID:   assetID,
		},
	)
	if err != nil {
		return nil, err
	}

	return &media, nil
}

func (s *assetsStore) GetAllAssetsByShopID(shopID string, limit int32, offset int32) ([]sqlc.GetAllAssetsByShopIDRow, error) {
	assets, err := s.queries.GetAllAssetsByShopID(
		context.Background(),
		sqlc.GetAllAssetsByShopIDParams{
			ShopID: shopID,
			Limit:  limit,
			Offset: offset,
		},
	)
	if err != nil {
		return nil, err
	}

	return assets, nil
}

func (s *assetsStore) GetAssetByIDAndShopID(assetID string, shopID string) (*sqlc.Asset, error) {
	asset, err := s.queries.GetAssetByIDAndShopID(
		context.Background(),
		sqlc.GetAssetByIDAndShopIDParams{
			ID:     assetID,
			ShopID: shopID,
		},
	)
	if err != nil {
		return nil, err
	}

	return &asset, nil
}

func (s *assetsStore) GetAssetDirPath(shopToken string) string {
	path := s.assetsPath + string(os.PathSeparator)
	path = path + shopToken + string(os.PathSeparator)

	return path
}

func (s *assetsStore) GetAssetByID(assetID string) (*sqlc.Asset, error) {
	asset, err := s.queries.GetAssetByID(
		context.Background(),
		assetID,
	)
	if err != nil {
		return nil, err
	}

	return &asset, nil
}

func (s *assetsStore) SaveAssetToDisk(b []byte, assetID string, shopToken string, fileExt string) (*sqlc.Asset, error) {
	asset, err := s.GetAssetByID(assetID)
	if err != nil {
		return nil, fmt.Errorf("GetAssetByID: %w", err)
	}

	path := s.GetAssetDirPath(shopToken)
	filename := asset.ID + fileExt

	err = os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return nil, err
	}

	path = path + filename

	err = ioutil.WriteFile(path, b, 0644)
	if err != nil {
		return nil, fmt.Errorf("ioutil.WriteFile: %w", err)
	}

	param := sqlc.UpdateAssetParams{
		ShopID:           null.NewString("", false),
		MurmurHash:       null.NewString("", false),
		OriginalFilename: null.NewString("", false),
		Extension:        null.StringFrom(fileExt),
		FileSize:         null.NewInt(0, false),
		Width:            null.NewInt(0, false),
		Height:           null.NewInt(0, false),
		ID:               asset.ID,
	}

	updatedAsset, err := s.queries.UpdateAsset(
		context.Background(),
		param,
	)
	if err != nil {
		return nil, fmt.Errorf("UpdateAsset: %w", err)
	}

	return &updatedAsset, nil
}

func (s *assetsStore) GetAssetByMurmur(murmurHash string, shopID string) (*sqlc.Asset, error) {
	params := sqlc.GetAssetByMurmurParams{
		MurmurHash: murmurHash,
		ShopID:     shopID,
	}

	asset, err := s.queries.GetAssetByMurmur(
		context.Background(),
		params,
	)
	if err != nil {
		return nil, err
	}

	return &asset, nil
}

func (s *assetsStore) GetAssetByOriginalName(originalName string, shopID string) (*sqlc.Asset, error) {
	params := sqlc.GetAssetByOriginalNameParams{
		OriginalFilename: null.StringFrom(originalName),
		ShopID:           shopID,
	}

	asset, err := s.queries.GetAssetByOriginalName(
		context.Background(),
		params,
	)
	if err != nil {
		return nil, err
	}

	return &asset, nil
}
