package laxo

import (
  "context"
  "regexp"

  "laxo.vn/laxo/laxo/sqlc"
  "github.com/jackc/pgx/v4"
  "github.com/go-ozzo/ozzo-validation/v4"
  "github.com/go-ozzo/ozzo-validation/v4/is"
)

type User struct {
  Model *sqlc.User
}

func (u *User) Validate() error {
  err := validation.ValidateStruct(u.Model,
    validation.Field(&u.Model.Email, validation.Required, validation.Length(3, 300), is.Email),
    validation.Field(&u.Model.Password, validation.Required, validation.Length(8, 128),
      validation.Match(regexp.MustCompile(`\d`)).Error("Password must contain a digit"),
      validation.Match(regexp.MustCompile(`[^\d]`)).Error("Password must have a letter")),
  )

  if err != nil {
    return err
  }

  // Making sure email doesn't exist yet
  _, err = Queries.GetUserByEmail(context.Background(), u.Model.Email)
  if err != pgx.ErrNoRows {
    return err
  }

  return nil
}

