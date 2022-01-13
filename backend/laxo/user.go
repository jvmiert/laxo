package laxo

import (
  "encoding/json"
  "context"
  "regexp"
  "errors"

  "laxo.vn/laxo/laxo/sqlc"
  "database/sql"
  "github.com/jackc/pgx/v4"
  "github.com/go-ozzo/ozzo-validation/v4"
  "github.com/go-ozzo/ozzo-validation/v4/is"
  "golang.org/x/crypto/bcrypt"
)

var ErrUserExists = errors.New("User already exists")

type UserReturn struct {
	ID         string       `json:"id"`
	Email      string       `json:"email"`
	Created    sql.NullTime `json:"created"`
  SessionKey string       `json:"sessionKey,omitempty"`
}

func (ur *UserReturn) JSON() ([]byte, error) {
  bytes, err := json.Marshal(ur)

  if err != nil {
    Logger.Error("UserReturn marshal error", "error", err)
    return bytes, err
  }

  return bytes, nil
}

type User struct {
  Model *sqlc.User
}

func (u *User) ValidateNew() error {
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

  // No rows exist yet with this email, we pass validation
  if err == pgx.ErrNoRows {
    return nil
  }

  // We returned a result, email exists
  if err == nil {
    return ErrUserExists
  } else {
    return err
  }
}

func SaveNewUserToDB(u *User) (UserReturn, error) {
  var ur UserReturn

  hash, err := bcrypt.GenerateFromPassword([]byte(u.Model.Password), 13)

  if err != nil {
    Logger.Error("Password hash error", "error", err)
    return ur, err
  }

  u.Model.Password = string(hash)

  savedUser, err := Queries.CreateUser(
    context.Background(),
    sqlc.CreateUserParams{
      Email: u.Model.Email,
      Password: u.Model.Password,
    },
  )

  if err != nil {
    Logger.Error("Save user to DB error", "error", err)
    return ur, err
  }

  ur.ID      = savedUser.ID
  ur.Email   = savedUser.Email
  ur.Created = savedUser.Created

  u.Model = &savedUser

  return ur, nil
}
