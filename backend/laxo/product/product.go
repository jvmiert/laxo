package product

import (
	"encoding/json"

	"laxo.vn/laxo/laxo/sqlc"
)

type Product struct {
  Model          *sqlc.Product              `json:"product"`
  MediaModels    []sqlc.ProductsMedia       `json:"-"`
  PlatformModel  *sqlc.ProductsPlatform     `json:"-"`
}

func (p *Product) JSON() ([]byte, error) {
  bytes, err := json.Marshal(p)

  if err != nil {
    return bytes, err
  }

  return bytes, nil
}
