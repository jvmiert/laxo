package laxo

import (
  "encoding/json"
  "context"
  "regexp"
  "errors"
  "strings"

  "laxo.vn/laxo/laxo/sqlc"
  "database/sql"
  "golang.org/x/text/message"
  "golang.org/x/text/number"
  "github.com/jackc/pgx/v4"
  "github.com/go-ozzo/ozzo-validation/v4"
  "github.com/go-ozzo/ozzo-validation/v4/is"
  "golang.org/x/crypto/bcrypt"
)

var ErrUserNotExist = errors.New("User does not exist")

// When the user db model is empty (not loaded from db)
var ErrModelUnpopulated = errors.New("User model is not retrieved from db")

var ValidationErrPwReqDigit = "password_requires_digit"
var ValidationErrPwReqLetter = "password_requires_letter"

type UserReturn struct {
	ID         string       `json:"id"`
	Email      string       `json:"email"`
	Created    sql.NullTime `json:"created"`
	Fullname   string       `json:"fullname"`
}

func GetUserRegistrationFailure(errs error, printer *message.Printer) validation.Errors {
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
    case is.ErrEmail.Code():
      errorString = printer.Sprintf("must be a valid email address")
    case ValidationErrPwReqDigit:
      errorString = printer.Sprintf("must contain a digit")
    case ValidationErrPwReqLetter:
      errorString = printer.Sprintf("must contain a letter")
    default:
      errorString = printer.Sprintf("unknown validation error")
    }

    newError := ozzoError.SetMessage(errorString)
    errMap[key] = newError
  }
  return errMap
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

func (u *User) ValidateNew(printer *message.Printer) error {
  err := validation.ValidateStruct(u.Model,
    validation.Field(&u.Model.Fullname, validation.Required),
    validation.Field(&u.Model.Email, validation.Required, validation.Length(3, 300), is.Email),
    validation.Field(&u.Model.Password, validation.Required, validation.Length(8, 128),
      validation.Match(regexp.MustCompile(`\d`)).ErrorObject(validation.NewError(ValidationErrPwReqDigit, ValidationErrPwReqDigit)),
      validation.Match(regexp.MustCompile(`[^\d]`)).ErrorObject(validation.NewError(ValidationErrPwReqLetter, ValidationErrPwReqLetter))),
  )

  if err != nil {
    return GetUserRegistrationFailure(err, printer)
  }

  // Making sure email doesn't exist yet
  lowerEmail := strings.ToLower(u.Model.Email)

  _, err = Queries.GetUserByEmail(context.Background(), strings.TrimSpace(lowerEmail))

  // No rows exist yet with this email, we pass validation
  if err == pgx.ErrNoRows {
    return nil
  }

  // We returned a result, email exists
  if err == nil {
    err = validation.Errors{
      "email": validation.NewError(
        "already_exists",
        printer.Sprintf("user already exists")),
    }
    return err
  } else {
    return err
  }
}

func (u *User) JSON() ([]byte, error) {
  var ur UserReturn
  ur.ID         = u.Model.ID
  ur.Email      = u.Model.Email
  ur.Created    = u.Model.Created
  ur.Fullname   = u.Model.Fullname

  bytes, err := json.Marshal(ur)

  if err != nil {
    Logger.Error("UserReturn marshal error", "error", err)
    return bytes, err
  }

  return bytes, nil
}

func LoginUser(email string, password string, printer *message.Printer) (*User, error) {
  user, err := RetrieveUserFromDBbyEmail(email)
  if err == ErrUserNotExist {
    err = validation.Errors{
      "email": validation.NewError(
        "not_exists",
        printer.Sprintf("Email not found")),
    }
    return nil, err
  }


  if err = user.CheckPassword(password); err != nil {
  err = validation.Errors{
    "password": validation.NewError(
      "pw_incorrect",
      printer.Sprintf("Password is incorrect")),
  }
  return nil, err
  }

  return user, nil
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
  lowerEmail := strings.ToLower(email)
  user, err := Queries.GetUserByEmail(
    context.Background(),
    strings.TrimSpace(lowerEmail),
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

  lowerEmail := strings.ToLower(u.Model.Email)
  savedUser, err := Queries.CreateUser(
    context.Background(),
    sqlc.CreateUserParams{
      Email: strings.TrimSpace(lowerEmail),
      Password: u.Model.Password,
      Fullname: u.Model.Fullname,
    },
  )

  if err != nil {
    Logger.Error("Save user to DB error", "error", err)
    return err
  }

  u.Model = &savedUser

  return nil
}
