package assets

import (
	"errors"
	"io"
	"net/http"
	"path"

	"github.com/jackc/pgx/v4"
	"github.com/twmb/murmur3"
	"laxo.vn/laxo/laxo"
	"laxo.vn/laxo/laxo/lazada"
	"laxo.vn/laxo/laxo/sqlc"
)

type Store interface {
  GetProductMediaByProductID(string) (*sqlc.ProductsMedia, error)
  GetProductMediaByMurmur(int64, string) (*sqlc.ProductsMedia, error)
  SaveNewProductMedia(int64, string, []byte, string, string, string) error
}

type Service struct {
  store       Store
  logger      *laxo.Logger
  server      *laxo.Server
}

func NewService(store Store, logger *laxo.Logger, server *laxo.Server) Service {
  return Service {
    store: store,
    logger: logger,
    server: server,
  }
}

func (s *Service) GetProductMediaByProductID(p string) (*sqlc.ProductsMedia, error) {
  return s.store.GetProductMediaByProductID(p)
}

func (s *Service) GetProductMediaByMurmur(m int64, p string) (*sqlc.ProductsMedia, error) {
  return s.store.GetProductMediaByMurmur(m, p)
}

func (s *Service) ExtractImagesListFromProductResponse(p *lazada.ProductsResponseProducts) ([]string, error) {
  var imageList []string

  for _, v := range p.Images {
    if v.Valid {
      imageList = append(imageList, v.String)
    }
  }

  for _, v := range p.MarketImages {
    if v.Valid {
      imageList = append(imageList, v.String)
    }
  }

  if len(p.Skus) > 0 {
    for _, v := range p.Skus[0].Images {
      if v.Valid {
        imageList = append(imageList, v.String)
      }
    }

    for _, v := range p.Skus[0].MarketImages {
      if v.Valid {
        imageList = append(imageList, v.String)
      }
    }
  }

  return imageList, nil
}

func (s *Service) SaveProductImages(i []string, shopID string, productID string, shopToken string) error {
  for _, v := range i {
    resp, err := http.Get(v)
    if err != nil {
      return err
    }
    defer resp.Body.Close()

    if resp.StatusCode != 200 {
      return errors.New("lazada response code is not 200")
    }

    b, err := io.ReadAll(resp.Body)
    if err != nil {
      return err
    }

    h64 := murmur3.New64()
    h64.Write(b)

    mID := int64(h64.Sum64())

    _, err = s.GetProductMediaByMurmur(mID, productID)
    if err != pgx.ErrNoRows && err != nil {
      return err
    }

    if err == pgx.ErrNoRows {
      filename := path.Base(v)
      err = s.store.SaveNewProductMedia(mID, filename, b, shopID, productID, shopToken)
      if err != nil {
        return err
      }
    }
  }
  return nil
}
