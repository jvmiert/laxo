package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"laxo.vn/laxo/laxo"
)

func testShopFunc(t *testing.T, authCookieName string, authToken string) {
  table := []struct {
    id         string
    cookie     *http.Cookie
    statusCode int
    body       []byte
  }{
      {"Create shop invalid", &http.Cookie{Name: authCookieName, Value: authToken}, http.StatusUnprocessableEntity, []byte(`{"foo": "ShopName"}`)},
      {"Create shop valid", &http.Cookie{Name: authCookieName, Value: authToken}, http.StatusOK, []byte(`{"shopName": "ShopName"}`)},
    }

  router := laxo.SetupRouter(true)

  for _, v := range table {
    t.Run(v.id, func(t *testing.T) {
      w := httptest.NewRecorder()
      r := httptest.NewRequest(http.MethodPost, "/api/shop", bytes.NewBuffer(v.body))
      r.Header.Set("Content-Type", "application/json")

      if v.cookie != nil {
        r.AddCookie(v.cookie)
      }

      router.ServeHTTP(w, r)

      if w.Code != v.statusCode {
        t.Errorf("Expected status code: %d, but got: %d", v.statusCode, w.Code)
      }
    })
  }
}
