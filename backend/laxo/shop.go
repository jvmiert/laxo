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
	ID         string       `json:"id"`
  UserID     string       `json:"userID"`
	ShopName   string       `json:"shopName"`
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
  sr.ShopName = s.Model.ShopName

  bytes, err := json.Marshal(sr)

  if err != nil {
    Logger.Error("Shop marshal error", "error", err)
    return bytes, err
  }

  return bytes, nil
}

func (s *Shop) ValidateNew(printer *message.Printer) error {
  err := validation.ValidateStruct(s.Model,
    validation.Field(&s.Model.ShopName, validation.Required, validation.Length(3, 300)),
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

func RetrieveShopsFromDBbyUserID(userID string) ([]Shop, error) {
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

