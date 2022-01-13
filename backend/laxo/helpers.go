package laxo

import (
  "net/http"
  "encoding/json"
  "errors"
  "strings"
  "fmt"
  "io"
  "crypto/rand"
)

func GenerateRandomString(n int) ([]byte, error) {
  b := make([]byte, n)
  if _, err := rand.Read(b); err != nil {
    return nil, err
  }
  return b, nil
}

type malformedRequest struct {
  status int
  msg    string
}

func (mr *malformedRequest) Error() string {
    return mr.msg
}

func decodeJSONBody(w http.ResponseWriter, r *http.Request, dst interface{}) error {
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
      Logger.Info("Decode error", "error", err, "msg", msg)
      return &malformedRequest{status: http.StatusBadRequest, msg: msg}

    case errors.Is(err, io.ErrUnexpectedEOF):
      msg := "Request body contains badly-formed JSON"
      Logger.Info("Decode error", "error", err, "msg", msg)
      return &malformedRequest{status: http.StatusBadRequest, msg: msg}

    case errors.As(err, &unmarshalTypeError):
      msg := fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
      Logger.Info("Decode error", "error", err, "msg", msg)
      return &malformedRequest{status: http.StatusBadRequest, msg: msg}

    case strings.HasPrefix(err.Error(), "json: unknown field "):
      fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
      msg := fmt.Sprintf("Request body contains unknown field %s", fieldName)
      Logger.Info("Decode error", "error", err, "msg", msg)
      return &malformedRequest{status: http.StatusBadRequest, msg: msg}

    case errors.Is(err, io.EOF):
      msg := "Request body must not be empty"
      Logger.Info("Decode error", "error", err, "msg", msg)
      return &malformedRequest{status: http.StatusBadRequest, msg: msg}

    case err.Error() == "http: request body too large":
      msg := "Request body must not be larger than 1MB"
      Logger.Info("Decode error", "error", err, "msg", msg)
      return &malformedRequest{status: http.StatusRequestEntityTooLarge, msg: msg}

    default:
      Logger.Error("Unknown decode error", "error", err)
      return err
    }
  }

	err = dec.Decode(&struct{}{})
  if err != io.EOF {
    msg := "Request body must only contain a single JSON object"
    Logger.Info("Decode error", "error", err, "msg", msg)
    return &malformedRequest{status: http.StatusBadRequest, msg: msg}
  }

  return nil
}
