package laxo

import (
	"context"
	"encoding/json"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"golang.org/x/text/message"
	"golang.org/x/text/number"
	"laxo.vn/laxo/laxo/sqlc"
)

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
  bytes, err := json.Marshal(s)

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

