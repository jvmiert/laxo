package store

import (
	"context"
	"errors"
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

  productPath := assetsBasePath + string(os.PathSeparator) + productImageDir

  err = os.MkdirAll(productPath, 0664)
  if err != nil {
    return nil, err
  }

  return &assetsStore{
    dir,
    store,
  }, nil
}

func (s *assetsStore) GetProductMediaByProductID(productID string) (*sqlc.ProductsMedia, error) {
  productMedia, err := s.queries.GetProductMediaByProductID(
    context.Background(),
    productID,
  )
  if err != nil {
    return nil, err
  }

  return &productMedia, nil
}

func (s *assetsStore) SaveNewProductMedia(mID int64, oFilename string, b []byte, shopID string, productID string) error {
  params := sqlc.CreateProductMediaParams{
    ProductID: productID,
    OriginalFilename: null.StringFrom(oFilename),
    MurmurHash: null.IntFrom(mID),
  }

  pMModel, err := s.queries.CreateProductMedia(
    context.Background(),
    params,
  )
  if err != nil {
    return err
  }
  filetype := http.DetectContentType(b)

  ext, err := mime.ExtensionsByType(filetype)
  if err != nil {
    return err
  }

  if len(ext) == 0 {
    return errors.New("no extension found for image")
  }

  filename := pMModel.ID + ext[0]

  path := s.assetsPath + string(os.PathSeparator) + productImageDir + string(os.PathSeparator) + filename

	err = ioutil.WriteFile(path, b, 0644)
	if err != nil {
		return nil
	}
  return nil
}

func (s *assetsStore) GetProductMediaByMurmur(murmurHash int64, productID string) (*sqlc.ProductsMedia, error) {
  params := sqlc.GetProductMediaByMurmurParams{
    MurmurHash: null.IntFrom(murmurHash),
    ProductID: productID,
  }

  productMedia, err := s.queries.GetProductMediaByMurmur(
    context.Background(),
    params,
  )
  if err != nil {
    return nil, err
  }

  return &productMedia, nil
}
