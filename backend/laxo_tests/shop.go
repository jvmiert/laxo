package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"laxo.vn/laxo/laxo"
)

func testShopCreateFunc(t *testing.T, authCookieName string, authToken string) {
	table := []struct {
		id         string
		cookie     *http.Cookie
		statusCode int
		body       []byte
	}{
		{"Create shop unauth", nil, http.StatusUnauthorized, []byte(`{"foo": "ShopName"}`)},
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

			assert.Equal(t, w.Code, v.statusCode, "Status code should be equal")
		})
	}
}

func testShopGetFunc(t *testing.T, authCookieName string, authToken string) {
	t.Fatal("shop get function tests not implemented yet")
}
