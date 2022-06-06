package shop

import (
	"encoding/json"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"laxo.vn/laxo/laxo/sqlc"
)

type ShopReturn struct {
	ID         string           `json:"id"`
  UserID     string           `json:"userID"`
	Name       string           `json:"name"`
  AssetsToken string           `json:"assetsToken"`
  Platforms  []PlatformReturn `json:"platforms"`
}

type PlatformReturn struct {
	ID         string   `json:"id"`
	Name       string   `json:"name"`
  Created    int64    `json:"created"`
}


type Shop struct {
  Model  *sqlc.Shop
}

func (s *Shop) JSON() ([]byte, error) {
  var sr ShopReturn
  sr.ID       = s.Model.ID
  sr.UserID   = s.Model.UserID
  sr.Name = s.Model.ShopName
  sr.AssetsToken = s.Model.AssetsToken

  bytes, err := json.Marshal(sr)

  if err != nil {
    return bytes, err
  }

  return bytes, nil
}

func (s *Shop) ValidateNew() error {
  return validation.ValidateStruct(s.Model,
    validation.Field(&s.Model.ShopName, validation.Required, validation.Length(6, 300)),
  )
}
