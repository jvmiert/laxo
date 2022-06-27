package user

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/jackc/pgx/v4"
	"github.com/mediocregopher/radix/v4"
	"golang.org/x/text/message"
	"golang.org/x/text/number"
	"laxo.vn/laxo/laxo"
	"laxo.vn/laxo/laxo/sqlc"
)

var ErrInvalidPW = errors.New("saved password does not match submitted")

// When the user db model is empty (not loaded from db)
var ErrModelUnpopulated = errors.New("user model is not retrieved from db")

var ValidationErrPwReqDigit = "password_requires_digit"
var ValidationErrPwReqLetter = "password_requires_letter"

type Store interface {
	RetrieveUserFromDBbyEmail(string) (*User, error)
	SaveNewUserToDB(sqlc.CreateUserParams) (*User, error)
	RetrieveUserFromDBbyID(string) (*User, error)
}

type Service struct {
	store  Store
	logger *laxo.Logger
	server *laxo.Server
}

func NewService(store Store, logger *laxo.Logger, server *laxo.Server) Service {
	return Service{
		store:  store,
		logger: logger,
		server: server,
	}
}

func (s *Service) GetUserRegistrationFailure(errs error, printer *message.Printer) validation.Errors {
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

func (s *Service) CheckPassword(u *User, p string) error {
	if u.Model == nil {
		return ErrModelUnpopulated
	}

	match, err := ComparePasswordAndHash(p, u.Model.Password)

	if err != nil {
		return err
	}

	if !match {
		return ErrInvalidPW
	}

	return nil
}

func (s *Service) ValidateNew(u *User, printer *message.Printer) error {
	err := validation.ValidateStruct(u.Model,
		validation.Field(&u.Model.Fullname, validation.Required),
		validation.Field(&u.Model.Email, validation.Required, validation.Length(3, 300), is.Email),
		validation.Field(&u.Model.Password, validation.Required, validation.Length(8, 128),
			validation.Match(regexp.MustCompile(`\d`)).ErrorObject(validation.NewError(ValidationErrPwReqDigit, ValidationErrPwReqDigit)),
			validation.Match(regexp.MustCompile(`[^\d]`)).ErrorObject(validation.NewError(ValidationErrPwReqLetter, ValidationErrPwReqLetter))),
	)

	if err != nil {
		return s.GetUserRegistrationFailure(err, printer)
	}

	// Making sure email doesn't exist yet
	lowerEmail := strings.ToLower(u.Model.Email)

	_, err = s.store.RetrieveUserFromDBbyEmail(strings.TrimSpace(lowerEmail))
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

func (s *Service) LoginUser(email string, password string, printer *message.Printer) (*User, error) {
	user, err := s.store.RetrieveUserFromDBbyEmail(email)
	if err == pgx.ErrNoRows {
		err = validation.Errors{
			"email": validation.NewError(
				"not_exists",
				printer.Sprintf("Email not found")),
		}
		return nil, err
	}

	if err = s.CheckPassword(user, password); err != nil {
		err = validation.Errors{
			"password": validation.NewError(
				"pw_incorrect",
				printer.Sprintf("Password is incorrect")),
		}
		return nil, err
	}

	return user, nil
}

func (s *Service) RetrieveUserFromDBbyID(userID string) (*User, error) {
	return s.store.RetrieveUserFromDBbyID(userID)
}

func (s *Service) SaveNewUserToDB(u *User) (*User, error) {
	hash, err := CreateHash(u.Model.Password)

	if err != nil {
		return nil, err
	}

	u.Model.Password = hash

	lowerEmail := strings.ToLower(u.Model.Email)

	savedUser, err := s.store.SaveNewUserToDB(
		sqlc.CreateUserParams{
			Email:    strings.TrimSpace(lowerEmail),
			Password: u.Model.Password,
			Fullname: u.Model.Fullname,
		},
	)

	if err != nil {
		return nil, err
	}

	u.Model = savedUser.Model

	return u, nil
}

func (s *Service) SetUserSession(u *User) (time.Time, string, error) {
	randomBytes, err := laxo.GenerateRandomString(128)

	if err != nil {
		s.logger.Errorw("Couldn't generate random bytes",
			"error", err,
		)
		return time.Time{}, "", err
	}

	sessionKey := base64.StdEncoding.EncodeToString(randomBytes)

	// Get the seconds till the token expires
	expiresT := time.Now().AddDate(0, 0, s.server.Config.AuthCookieExpire)
	expires := time.Until(expiresT)
	expireString := fmt.Sprintf("%.0f", expires.Seconds())

	ctx := context.Background()
	if err := s.server.RedisClient.Do(ctx, radix.Cmd(nil, "SETEX", sessionKey, expireString, u.Model.ID)); err != nil {
		s.logger.Errorw("Couldn't set user session in Redis",
			"error", err,
		)
		return time.Time{}, "", err
	}

	return expiresT, sessionKey, nil
}

func (s *Service) RemoveUserSession(sessionToken string) error {
	ctx := context.Background()
	if err := s.server.RedisClient.Do(ctx, radix.Cmd(nil, "DEL", sessionToken)); err != nil {
		s.logger.Errorw("Couldn't remove user session in Redis",
			"error", err,
		)
		return err
	}
	return nil
}

func (s *Service) SetUserCookie(sessionToken string, w http.ResponseWriter, t time.Time) {
	authCookie := &http.Cookie{
		Name:     s.server.Config.AuthCookieName,
		Path:     "/",
		Value:    sessionToken,
		HttpOnly: true,
		Secure:   true,
		Expires:  t,
	}

	http.SetCookie(w, authCookie)
}

func (s *Service) RemoveUserCookie(w http.ResponseWriter) {
	authCookie := &http.Cookie{
		Name:     s.server.Config.AuthCookieName,
		Value:    "",
		HttpOnly: true,
		Secure:   true,
		Expires:  time.Unix(0, 0),
	}

	http.SetCookie(w, authCookie)
}
