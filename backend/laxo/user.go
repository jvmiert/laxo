package laxo

import (
  "laxo.vn/laxo/laxo/sqlc"
  "github.com/go-ozzo/ozzo-validation/v4"
)

type User struct {
  model *sqlc.User
}

func (u *User) Validate() error {
  return validation.ValidateStruct(&u,
    validation.Field(&u.model.Username, validation.Required, validation.Length(5, 50)),
  )
}

