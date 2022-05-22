package shop

import (
	"encoding/json"

	"gopkg.in/guregu/null.v4"
	"laxo.vn/laxo/laxo/sqlc"
)

type ProductPlatformInformation struct {
  ID           string          `json:"id"`
  Name         null.String     `json:"name"`
  ProductURL   null.String     `json:"productURL"`
}

type Product struct {
  Model           *sqlc.Product                   `json:"product"`
  MediaModels     []sqlc.ProductsMedia            `json:"-"`
  PlatformModel   *sqlc.ProductsPlatform          `json:"-"`
  MediaList       []string                        `json:"mediaList"`
  Platforms       []ProductPlatformInformation    `json:"platforms"`
}

func (p *Product) JSON() ([]byte, error) {
  bytes, err := json.Marshal(p)

  if err != nil {
    return bytes, err
  }

  return bytes, nil
}
