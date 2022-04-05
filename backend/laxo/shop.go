package laxo

import (
	"context"
	"encoding/json"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"golang.org/x/text/message"
	"golang.org/x/text/number"
	"laxo.vn/laxo/laxo/sqlc"
)

type ShopReturn struct {
	ID         string           `json:"id"`
  UserID     string           `json:"userID"`
	Name       string           `json:"name"`
  Platforms  []PlatformReturn `json:"platforms"`
}

type PlatformReturn struct {
	ID         string   `json:"id"`
	Name       string   `json:"name"`
  Created    int64    `json:"created"`
}

func GetShopCreateFailure(errs error, printer *message.Printer) validation.Errors {
  errMap := errs.(validation.Errors)

  var errorString string
  for key, err := range errMap {
    ozzoError := err.(validation.Error)
    code := ozzoError.Code()
    params := ozzoError.Params()

    switch code {
    case validation.ErrRequired.Code():
      errorString = printer.Sprintf("cannot be blank")
    case validation.ErrLengthOutOfRange.Code():
      errorString = printer.Sprintf("the length must be between %v and %v", number.Decimal(params["min"]), number.Decimal(params["max"]))
    default:
      errorString = printer.Sprintf("unknown validation error")
    }

    newError := ozzoError.SetMessage(errorString)
    errMap[key] = newError
  }
  return errMap
}

type Shop struct {
  Model  *sqlc.Shop
}

func (s *Shop) JSON() ([]byte, error) {
  var sr ShopReturn
  sr.ID       = s.Model.ID
  sr.UserID   = s.Model.UserID
  sr.Name = s.Model.ShopName

  bytes, err := json.Marshal(sr)

  if err != nil {
    Logger.Error("Shop marshal error", "error", err)
    return bytes, err
  }

  return bytes, nil
}

func (s *Shop) ValidateNew(printer *message.Printer) error {
  err := validation.ValidateStruct(s.Model,
    validation.Field(&s.Model.ShopName, validation.Required, validation.Length(6, 300)),
  )

  if err != nil {
    return GetShopCreateFailure(err, printer)
  }
  return nil
}

func GenerateShopList(ss []Shop) ([]byte, error) {
  sList := []json.RawMessage{}
  for _, s := range ss {
    b, err := s.JSON()
    if err != nil {
      return nil, err
    }
    j := json.RawMessage(b)
    sList = append(sList, j)
  }

	shopData := map[string]interface{}{
		"shops": sList,
		"total": len(sList),
	}

  bytes, err := json.Marshal(shopData)

  if err != nil {
    Logger.Error("Shop list marshal error", "error", err)
    return bytes, err
  }
  return bytes, nil
}

func GenerateShopPlatformList(ss []sqlc.GetShopsPlatformsByUserIDRow) ([]byte, error) {
  srList := []ShopReturn{}
  srMap := make(map[string]struct{})
  pMap := make(map[string][]PlatformReturn)

  for _, s := range ss {
    if _, present := srMap[s.ID]; !present {
      var sr ShopReturn
      sr.ID       = s.ID
      sr.UserID   = s.UserID
      sr.Name = s.ShopName

      srList = append(srList, sr)

      srMap[s.ID] = struct{}{}
    }

    created := int64(0)

    if s.PlatformCreated.Valid {
      created = s.PlatformCreated.Time.Unix()
    }

    if s.PlatformID.Valid {
      pMap[s.ID] = append(pMap[s.ID], PlatformReturn{
        ID: s.PlatformID.String,
        Name: s.PlatformName.String,
        Created: created,
      })
    }
  }

  srListU := []ShopReturn{}
  for _, s := range srList {
    if _, present := pMap[s.ID]; present {
      s.Platforms = pMap[s.ID]
      srListU = append(srListU, s)
    } else {
      s.Platforms = make([]PlatformReturn, 0)
      srListU = append(srListU, s)
    }

  }

	shopData := map[string]interface{}{
		"shops": srListU,
		"total": len(srListU),
	}

  bytes, err := json.Marshal(shopData)

  if err != nil {
    Logger.Error("Shop platform list marshal error", "error", err)
    return bytes, err
  }
  return bytes, nil
}

func RetrieveShopsByUserID(userID string) ([]Shop, error) {
  shops, err := Queries.GetShopsByUserID(
    context.Background(),
    userID,
  )

  if err != nil {
    Logger.Debug("Shop retrieval error", err)
    return nil, err
  }

  var sReturn []Shop

  for _, s := range shops {
    sReturn = append(sReturn, Shop{Model: &s})
  }

  return sReturn, nil
}

func RetrieveShopsPlatformsByUserID(userID string) ([]sqlc.GetShopsPlatformsByUserIDRow, error) {
  shops, err := Queries.GetShopsPlatformsByUserID(
    context.Background(),
    userID,
  )

  if err != nil {
    Logger.Debug("Shop retrieval error", err)
    return nil, err
  }

  return shops, nil
}

func SaveNewShopToDB(s *Shop, u string) error {
  savedShop, err := Queries.CreateShop(
    context.Background(),
    sqlc.CreateShopParams{
      ShopName: s.Model.ShopName,
      UserID: u,
    },
  )

  if err != nil {
    Logger.Error("Save shop to DB error", "error", err)
    return err
  }

  s.Model = &savedShop

  return nil
}

