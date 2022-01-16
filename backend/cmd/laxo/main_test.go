package main

import (
  "os"
  "bytes"
	"net/http"
  "strings"
	"testing"
	"net/http/httptest"

  "laxo.vn/laxo/laxo"
  "github.com/joho/godotenv"
)

func setupTest(t *testing.T) *laxo.Config {
  os.Chdir("./../..")

  if err := godotenv.Load(".env"); err != nil {
    t.Error("Failed to load .env file")
  }

  _, config := laxo.InitConfig(true)

  if err := laxo.InitRedis(); err != nil {
    t.Fatalf("Failed to init Redis: %v", err)
    return nil
  }

  uri := os.Getenv("POSTGRESQL_TEST_URL")

  if err := laxo.InitDatabase(uri); err != nil {
    t.Fatalf("Failed to init Database: %v", err)
    return nil
  }

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
  jsonStr = []byte(`{"email": "example@example.com", "password": "incorrect123"}`)
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
      found = true
    }
  }

  if !found {
    t.Error("Auth cookie not found after user creation")
  }

  // Posting the same user again should give an error
  jsonStr = []byte(`{"email": "example@example.com", "password": "incorrect123"}`)
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
}

