package laxo

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type ErrorReturn struct {
	ErrorDetails error `json:"errorDetails"`
	Error        bool  `json:"error"`
}

func ErrorJSONEncode(w http.ResponseWriter, error error, code int) {
	var returnError ErrorReturn

	/*
	     @TODO: In addition to the i18n error string, return a error
	            code that will be consistent between languages.

	  ozzoError := error.(validation.Errors)

	  Logger.Debug("test", ozzoError["email"].(validation.Error).Code())

	*/

	returnError.Error = true
	returnError.ErrorDetails = error

	b, err := json.Marshal(returnError)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	w.Write(b)
}

func GenerateRandomString(n int) ([]byte, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return nil, err
	}
	return b, nil
}

type MalformedRequest struct {
	Status int
	Msg    string
}

func (mr *MalformedRequest) Error() string {
	return mr.Msg
}

func DecodeJSONBody(logger *Logger, w http.ResponseWriter, r *http.Request, dst interface{}) error {
	// Don't accept JSON larger than 10MB
	r.Body = http.MaxBytesReader(w, r.Body, 10485760)

	dec := json.NewDecoder(r.Body)

	err := dec.Decode(&dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		case errors.As(err, &syntaxError):
			msg := fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset)
			logger.Errorw("Decode error",
				"error", err,
				"msg", msg,
			)
			return &MalformedRequest{Status: http.StatusBadRequest, Msg: msg}

		case errors.Is(err, io.ErrUnexpectedEOF):
			msg := "Request body contains badly-formed JSON"
			logger.Errorw("Decode error",
				"error", err,
				"msg", msg,
			)
			return &MalformedRequest{Status: http.StatusBadRequest, Msg: msg}

		case errors.As(err, &unmarshalTypeError):
			msg := fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
			logger.Errorw("Decode error",
				"error", err,
				"msg", msg,
			)
			return &MalformedRequest{Status: http.StatusBadRequest, Msg: msg}

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			msg := fmt.Sprintf("Request body contains unknown field %s", fieldName)
			logger.Errorw("Decode error",
				"error", err,
				"msg", msg,
			)
			return &MalformedRequest{Status: http.StatusBadRequest, Msg: msg}

		case errors.Is(err, io.EOF):
			msg := "Request body must not be empty"
			logger.Errorw("Decode error",
				"error", err,
				"msg", msg,
			)
			return &MalformedRequest{Status: http.StatusBadRequest, Msg: msg}

		case err.Error() == "http: request body too large":
			msg := "Request body must not be larger than 10MB"
			logger.Errorw("Decode error",
				"error", err,
				"msg", msg,
			)
			return &MalformedRequest{Status: http.StatusRequestEntityTooLarge, Msg: msg}

		default:
			logger.Errorw("Unknown decode error",
				"error", err,
			)
			return err
		}
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		msg := "Request body must only contain a single JSON object"
		logger.Errorw("Decode error",
			"error", err,
			"msg", msg,
		)
		return &MalformedRequest{Status: http.StatusBadRequest, Msg: msg}
	}

	return nil
}
