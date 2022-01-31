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
var ErrUserNotExist = errors.New("User does not exist")

// Used for returning to frontend
var ValErrWrongPassword = "Password is incorrect"
var ValErrUnknownEmail  = "Email not found"

// When the user db model is empty (not loaded from db)
var ErrModelUnpopulated = errors.New("User model is not retrieved from db")

type UserReturn struct {
	ID         string       `json:"id"`
	Email      string       `json:"email"`
	Created    sql.NullTime `json:"created"`
}

type UserLoginErrorMessage struct {
  Email    string `json:"email,omitempty"`
  Password string `json:"password,omitempty"`
}

func GetUserLoginFailure(emailFailed bool, pwFailed bool) ([]byte, error) {
  r := &UserLoginErrorMessage{}

  if emailFailed {
    r.Email = ValErrUnknownEmail
  }

  if pwFailed {
    r.Password = ValErrWrongPassword
  }

  bytes, err := json.Marshal(r)

  if err != nil {
    return bytes, err
  }

  return bytes, err
}

type LoginRequest struct {
  Email    string `json:"email"`
  Password string `json:"password"`
}

type User struct {
  Model       *sqlc.User
  SessionKey  string
}

func (u *User) CheckPassword(p string) error {
  if u.Model == nil {
    return ErrModelUnpopulated
  }

  err := bcrypt.CompareHashAndPassword([]byte(u.Model.Password), []byte(p))

  if err != nil {
    return err
  }
  return nil
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

func (u *User) JSON() ([]byte, error) {
  var ur UserReturn
  ur.ID         = u.Model.ID
  ur.Email      = u.Model.Email
  ur.Created    = u.Model.Created

  bytes, err := json.Marshal(ur)

  if err != nil {
    Logger.Error("UserReturn marshal error", "error", err)
    return bytes, err
  }

  return bytes, nil
}

func RetrieveUserFromDBbyID(uID string) (*User, error) {
  user, err := Queries.GetUserByID(
    context.Background(),
    uID,
  )

  if err == pgx.ErrNoRows {
    return nil, ErrUserNotExist
  } else if err != nil {
    Logger.Debug("User retrieval error", err)
    return nil, err
  }

  var u User
  u.Model = &user

  return &u, nil
}

func RetrieveUserFromDBbyEmail(email string) (*User, error) {
  user, err := Queries.GetUserByEmail(
    context.Background(),
    email,
  )

  if err == pgx.ErrNoRows {
    return nil, ErrUserNotExist
  } else if err != nil {
    Logger.Debug("User retrieval error", err)
    return nil, err
  }

  var u User
  u.Model = &user

  return &u, nil
}

func SaveNewUserToDB(u *User) error {
  hash, err := bcrypt.GenerateFromPassword([]byte(u.Model.Password), 13)

  if err != nil {
    Logger.Error("Password hash error", "error", err)
    return err
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
    return err
  }

  u.Model = &savedUser

  return nil
}
