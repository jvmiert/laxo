package main

import (
  "os"
  "bytes"
	"net/http"
  "encoding/json"
  "strings"
	"testing"
	"net/http/httptest"

  "laxo.vn/laxo/laxo"
  "github.com/joho/godotenv"
)

type TestState struct {
  CreateUserToken string
  Config *laxo.Config
}

var state TestState

func setupTest(t *testing.T) *laxo.Config {
  os.Chdir("./..")

  if err := godotenv.Load(".env"); err != nil {
    t.Error("Failed to load .env file")
  }

  _, config := laxo.InitConfig(true)

  redisURI := os.Getenv("REDIS_TEST_URL")

  if err := laxo.InitRedis(redisURI); err != nil {
    t.Fatalf("Failed to init Redis: %v", err)
    return nil
  }

  dbURI := os.Getenv("POSTGRESQL_TEST_URL")

  if err := laxo.InitDatabase(dbURI); err != nil {
    t.Fatalf("Failed to init Database: %v", err)
    return nil
  }

  state.Config = &config

  return &config
}

func TestRouteCreateUser(t *testing.T) {
  config := setupTest(t)

  // Testing empty post
  req, err := http.NewRequest("POST", "/api/user", nil)
  if err != nil {
    t.Fatal(err)
  }

  rr := httptest.NewRecorder()
  r := laxo.SetupRouter()

  r.ServeHTTP(rr, req)

  // Empty post should return StatusUnsupportedMediaType
  if status := rr.Code; status != http.StatusUnsupportedMediaType{
    t.Errorf("handler returned wrong status code: got %v want %v",
      status, http.StatusUnsupportedMediaType)
  }

  // Testing empty JSON post
  var jsonStr = []byte(`{}`)
  req, err = http.NewRequest("POST", "/api/user", bytes.NewBuffer(jsonStr))
  if err != nil {
    t.Fatal(err)
  }
  req.Header.Set("Content-Type", "application/json")

  rr = httptest.NewRecorder()
  r.ServeHTTP(rr, req)

  // Empty JSON post should return 422 validation error
  if status := rr.Code; status != http.StatusUnprocessableEntity{
    t.Errorf("handler returned wrong status code: got %v want %v",
      status, http.StatusUnprocessableEntity)
  }

  // Posting malformed JSON
  jsonStr = []byte(`{"malformed" "json"}`)
  req, err = http.NewRequest("POST", "/api/user", bytes.NewBuffer(jsonStr))
  if err != nil {
    t.Fatal(err)
  }
  req.Header.Set("Content-Type", "application/json")

  rr = httptest.NewRecorder()
  r.ServeHTTP(rr, req)

  if status := rr.Code; status != http.StatusBadRequest{
    t.Errorf("handler returned wrong status code: got %v want %v",
      status, http.StatusBadRequest)
  }

  // Posting an invalid email
  jsonStr = []byte(`{"email": "wrong email"}`)
  req, err = http.NewRequest("POST", "/api/user", bytes.NewBuffer(jsonStr))
  if err != nil {
    t.Fatal(err)
  }
  req.Header.Set("Content-Type", "application/json")

  rr = httptest.NewRecorder()
  r.ServeHTTP(rr, req)

  if status := rr.Code; status != http.StatusUnprocessableEntity{
    t.Errorf("handler returned wrong status code: got %v want %v",
      status, http.StatusUnprocessableEntity)
  }

  if !strings.Contains(rr.Body.String(), "must be a valid email") {
   t.Error("Email validation did not return correctly")
  }

  // Post an invalid email with vi locale
  jsonStr = []byte(`{"email": "wrong email"}`)
  req, err = http.NewRequest("POST", "/api/user", bytes.NewBuffer(jsonStr))
  if err != nil {
    t.Fatal(err)
  }
  req.Header.Set("Content-Type", "application/json")
  req.Header.Set("locale", "vi")

  rr = httptest.NewRecorder()
  r.ServeHTTP(rr, req)

  if status := rr.Code; status != http.StatusUnprocessableEntity{
    t.Errorf("handler returned wrong status code: got %v want %v",
      status, http.StatusUnprocessableEntity)
  }

  if !strings.Contains(rr.Body.String(), "Phải la một địa chỉ") {
   t.Error("Email validation did not return correctly")
  }

  // Posting a valid email without a valid password
  jsonStr = []byte(`{"email": "example@example.com", "password": "incorrect"}`)
  req, err = http.NewRequest("POST", "/api/user", bytes.NewBuffer(jsonStr))
  if err != nil {
    t.Fatal(err)
  }
  req.Header.Set("Content-Type", "application/json")

  rr = httptest.NewRecorder()
  r.ServeHTTP(rr, req)

  if status := rr.Code; status != http.StatusUnprocessableEntity{
    t.Errorf("handler returned wrong status code: got %v want %v",
      status, http.StatusUnprocessableEntity)
  }

  if !strings.Contains(rr.Body.String(), "must contain a digit") {
   t.Error("Password validation did not return correctly")
  }

  // Posting a valid user creation request
  jsonStr = []byte(`{"email": "example@example.com", "password": "correct123", "fullname": "first name last name"}`)
  req, err = http.NewRequest("POST", "/api/user", bytes.NewBuffer(jsonStr))
  if err != nil {
    t.Fatal(err)
  }
  req.Header.Set("Content-Type", "application/json")

  rr = httptest.NewRecorder()
  r.ServeHTTP(rr, req)

  if status := rr.Code; status != http.StatusOK{
    t.Errorf("handler returned wrong status code: got %v want %v",
      status, http.StatusOK)
  }

  if len(rr.Result().Cookies()) == 0 {
    t.Error("Cookie not present after user creation")
  }

  found := false
  for _, c := range rr.Result().Cookies() {
    if c.Name == config.AuthCookieName {
      state.CreateUserToken = c.Value
      found = true
    }
  }

  if !found {
    t.Error("Auth cookie not found after user creation")
  }

  // Posting the same user again should give an error
  jsonStr = []byte(`{"email": "example@example.com", "password": "correct123", "fullname": "first name last name"}`)
  req, err = http.NewRequest("POST", "/api/user", bytes.NewBuffer(jsonStr))
  if err != nil {
    t.Fatal(err)
  }
  req.Header.Set("Content-Type", "application/json")

  rr = httptest.NewRecorder()
  r.ServeHTTP(rr, req)

  if status := rr.Code; status != http.StatusUnprocessableEntity{
    t.Errorf("handler returned wrong status code: got %v want %v",
      status, http.StatusUnprocessableEntity)
  }

  if !strings.Contains(rr.Body.String(), "already exists") {
   t.Error("User duplication validation did not return correctly")
  }
}

func TestGetUser(t *testing.T) {
  req, err := http.NewRequest("GET", "/api/user", nil)
  if err != nil {
    t.Fatal(err)
  }

  rr := httptest.NewRecorder()
  r := laxo.SetupRouter()

  r.ServeHTTP(rr, req)

  // Route should return 401 without a valid cookie
  if status := rr.Code; status != http.StatusUnauthorized{
    t.Errorf("handler returned wrong status code: got %v want %v",
      status, http.StatusUnauthorized)
  }

  // Testing if we can retrieve the user info of the newly created account in
  // the previous test.
  req, err = http.NewRequest("GET", "/api/user", nil)
  if err != nil {
    t.Fatal(err)
  }

  req.AddCookie(&http.Cookie{Name: state.Config.AuthCookieName, Value: state.CreateUserToken})

  rr = httptest.NewRecorder()

  r.ServeHTTP(rr, req)

  // Route should return 200 with a valid cookie
  if status := rr.Code; status != http.StatusOK{
    t.Errorf("handler returned wrong status code: got %v want %v",
      status, http.StatusOK)
  }
}

func TestLogin(t *testing.T) {
  // Testing an incorrect password
  jsonStr := []byte(`{"email": "example@example.com", "password": "incorrect"}`)
  req, err := http.NewRequest("POST", "/api/login", bytes.NewBuffer(jsonStr))
  if err != nil {
    t.Fatal(err)
  }
  req.Header.Set("Content-Type", "application/json")

  rr := httptest.NewRecorder()
  r := laxo.SetupRouter()

  r.ServeHTTP(rr, req)

  // Route should return 401 with a wrong password
  if status := rr.Code; status != http.StatusUnauthorized{
    t.Errorf("handler returned wrong status code: got %v want %v",
      status, http.StatusUnauthorized)
  }

  // Route should return correct JSON error message
  var errorResponse map[string]interface{}

  err = json.Unmarshal(rr.Body.Bytes(), &errorResponse)

  if err != nil {
    t.Errorf("couldn't unmarshal error return: %v", err)
  }

  // The pw incorrect error
  var ErrPWIncorrect = "Password is incorrect"

  responseValue := errorResponse["errorDetails"].(map[string]interface{})["password"]
  if responseValue != ErrPWIncorrect {
    t.Errorf("Wrong error response: got %v want %v", responseValue, "Password is incorrect")
  }

  // Testing an incorrect password with vi locale
  jsonStr = []byte(`{"email": "example@example.com", "password": "incorrect"}`)
  req, err = http.NewRequest("POST", "/api/login", bytes.NewBuffer(jsonStr))
  if err != nil {
    t.Fatal(err)
  }
  req.Header.Set("Content-Type", "application/json")
  req.Header.Set("locale", "vi")

  rr = httptest.NewRecorder()

  r.ServeHTTP(rr, req)

  // Route should return 401 with a wrong password
  if status := rr.Code; status != http.StatusUnauthorized{
    t.Errorf("handler returned wrong status code: got %v want %v",
      status, http.StatusUnauthorized)
  }

  // Route should return correct JSON error message
  err = json.Unmarshal(rr.Body.Bytes(), &errorResponse)

  if err != nil {
    t.Errorf("couldn't unmarshal error return: %v", err)
  }

  // The pw incorrect error
  var ErrPWIncorrectVI = "Mật khẩu không đúng"

  responseValue = errorResponse["errorDetails"].(map[string]interface{})["password"]

  if responseValue != ErrPWIncorrectVI {
    t.Errorf("Wrong error response: got %v want %v", responseValue, "Mật khẩu không đúng")
  }

  // Testing a correct password
  jsonStr = []byte(`{"email": "example@example.com", "password": "correct123"}`)
  req, err = http.NewRequest("POST", "/api/login", bytes.NewBuffer(jsonStr))
  if err != nil {
    t.Fatal(err)
  }
  req.Header.Set("Content-Type", "application/json")

  rr = httptest.NewRecorder()

  r.ServeHTTP(rr, req)

  // Route should return 200 with a correct password
  if status := rr.Code; status != http.StatusOK{
    t.Errorf("handler returned wrong status code: got %v want %v",
      status, http.StatusOK)
  }

  if len(rr.Result().Cookies()) == 0 {
    t.Error("Cookie not present after login")
  }

  found := false
  for _, c := range rr.Result().Cookies() {
    if c.Name == state.Config.AuthCookieName {
      state.CreateUserToken = c.Value
      found = true
    }
  }

  if !found {
    t.Error("Auth cookie not found after login")
  }
}

