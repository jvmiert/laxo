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

var ErrUserExists = errors.New("User already exists")
var ErrUserNotExist = errors.New("User does not exist")

// When the user db model is empty (not loaded from db)
var ErrModelUnpopulated = errors.New("User model is not retrieved from db")

var ValidationErrPwReqDigit = "password_requires_digit"
var ValidationErrPwReqLetter = "password_requires_letter"

type UserReturn struct {
	ID         string       `json:"id"`
	Email      string       `json:"email"`
	Created    sql.NullTime `json:"created"`
}

type UserLoginErrorMessage struct {
  Email    string `json:"email,omitempty"`
  Password string `json:"password,omitempty"`
}

func GetUserLoginFailure(emailFailed bool, pwFailed bool, printer *message.Printer) ([]byte, error) {
  r := &UserLoginErrorMessage{}

  if emailFailed {
    r.Email = printer.Sprintf("Email not found")
  }

  if pwFailed {
    r.Password = printer.Sprintf("Password is incorrect")
  }

  bytes, err := json.Marshal(r)

  if err != nil {
    return bytes, err
  }

  return bytes, err
}

func GetUserRegistrationFailure(errs error, printer *message.Printer) validation.Errors {
  // @TODO:
  //   - Create a switch statement with all the possible validation error codes and return correct translated string
  //   - Create a proper object to return to frontend (map, marshall)
  //     - https://github.com/go-ozzo/ozzo-validation#validation-errors

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
    default:
      errorString = "unknown error"
    }

    ozzoError.SetMessage(errorString)

    errMap[key] = ozzoError

    Logger.Debug("GetUserRegFailure", "key", key, "code", code, "error", err, "params", params)
    Logger.Debug("Translate result", errorString)
  }

  b, _ := json.Marshal(errMap)
  Logger.Debug("Marshal result", "bytes", string(b))
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
    validation.Field(&u.Model.Email, validation.Required, validation.Length(3, 300), is.Email),
    validation.Field(&u.Model.Password, validation.Required, validation.Length(8, 128),
      validation.Match(regexp.MustCompile(`\d`)).ErrorObject(validation.NewError(ValidationErrPwReqDigit, ValidationErrPwReqDigit)),
      validation.Match(regexp.MustCompile(`[^\d]`)).ErrorObject(validation.NewError(ValidationErrPwReqDigit, ValidationErrPwReqDigit))),
  )

  if err != nil {
    Logger.Debug("Original error:", "error", err)
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
    // @TODO: validation.NewError()
    //        we need to create a new validation error and
    //        return that instead of a go error.
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
    },
  )

  if err != nil {
    Logger.Error("Save user to DB error", "error", err)
    return err
  }

  u.Model = &savedUser

  return nil
}
