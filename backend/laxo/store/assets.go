package store

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"mime"
	"net/http"
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
    ShopID: ShopID,
    MurmurHash: MurmurHash,
    OriginalFilename: null.StringFrom(OriginalFilename),
    Extension: null.NewString("", false),
    FileSize: null.IntFrom(FileSize),
    Width: null.IntFrom(Width),
    Height: null.IntFrom(Height),
  }

  asset, err := s.queries.CreateAsset(
    context.Background(),
    param,
  )

  return &asset, err
}

func (s *assetsStore) UpdateProductMedia(assetID string, productID string, status null.String, order null.Int) (*sqlc.ProductsMedia, error) {
  media, err := s.queries.UpdateProductMedia(
    context.Background(),
    sqlc.UpdateProductMediaParams{
      ProductID: productID,
      AssetID: assetID,
      Status: status,
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
      ProductID: productID,
      AssetID: assetID,
      Status: status,
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
      AssetID: assetID,
    },
  )
  if err != nil {
    return nil, err
  }

  return &media, nil
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

  filename := asset.ID + fileExt

  path := s.assetsPath + string(os.PathSeparator)
  path = path + shopToken + string(os.PathSeparator)

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
    ShopID: null.NewString("", false),
    MurmurHash: null.NewString("", false),
    OriginalFilename: null.NewString("", false),
    Extension: null.StringFrom(fileExt),
    FileSize: null.NewInt(0, false),
    Width: null.NewInt(0, false),
    Height: null.NewInt(0, false),
    ID: asset.ID,
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

func (s *assetsStore) SaveNewProductMedia(mID int64, oFilename string, b []byte, shopID string, productID string, shopToken string) error {
  filetype := http.DetectContentType(b)

  ext, err := mime.ExtensionsByType(filetype)
  if err != nil {
    return err
  }

  if len(ext) == 0 {
    return errors.New("no extension found for image")
  }

  fileExt := ""

  for _, e := range ext {
    if e == ".jpg" {
      fileExt = ".jpg"
      break
    }
  }

  if fileExt == "" {
    fileExt = ext[0]
  }

  params := sqlc.CreateProductMediaParams{
    ProductID: productID,
    //OriginalFilename: null.StringFrom(oFilename),
    //MurmurHash: null.IntFrom(mID),
    //Extension: null.StringFrom(fileExt),
  }

  //pMModel, err := s.queries.CreateProductMedia(
  _, err = s.queries.CreateProductMedia(
    context.Background(),
    params,
  )
  if err != nil {
    return err
  }

  //filename := pMModel.ID + fileExt

  path := s.assetsPath + string(os.PathSeparator)
  path = path + shopToken + string(os.PathSeparator)
  path = path + productImageDir + string(os.PathSeparator)

  err = os.MkdirAll(path, os.ModePerm)
  if err != nil {
    return err
  }

  path = path //+ filename

	err = ioutil.WriteFile(path, b, 0644)
	if err != nil {
		return nil
	}
  return nil
}

func (s *assetsStore) GetAssetByMurmur(murmurHash string, shopID string) (*sqlc.Asset, error) {
  params := sqlc.GetAssetByMurmurParams{
    MurmurHash: murmurHash,
    ShopID: shopID,
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
