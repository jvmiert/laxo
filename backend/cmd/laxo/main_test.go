package main

import (
  "os"
	"net/http"
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
func TestRouteUser(t *testing.T) {
  setupTest(t)

  // Testing empty post
  req, err := http.NewRequest("POST", "/api/user", nil)
  if err != nil {
    t.Fatal(err)
  }

  rr := httptest.NewRecorder()
  r := laxo.SetupRouter()

  r.ServeHTTP(rr, req)

  if status := rr.Code; status != http.StatusUnsupportedMediaType{
    t.Errorf("handler returned wrong status code: got %v want %v",
      status, http.StatusOK)
  }
}
