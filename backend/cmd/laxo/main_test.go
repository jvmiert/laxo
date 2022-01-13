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

func setupTest(t *testing.T) {
  os.Chdir("./../..")

  if err := godotenv.Load(".env"); err != nil {
    t.Error("Failed to load .env file")
  }

  var config laxo.Config
  _ = laxo.InitConfig(&config, true)

  if err := laxo.InitRedis(); err != nil {
    t.Fatalf("Failed to init Redis: %v", err)
    return
  }

  uri := os.Getenv("POSTGRESQL_TEST_URL")

  if err := laxo.InitDatabase(uri); err != nil {
    t.Fatalf("Failed to init Database: %v", err)
    return
  }
}

func TestRouteCreateUser(t *testing.T) {
  setupTest(t)

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
