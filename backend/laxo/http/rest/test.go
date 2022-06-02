package rest

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
	"laxo.vn/laxo/laxo"
	"laxo.vn/laxo/laxo/assets"
	"laxo.vn/laxo/laxo/lazada"
	"laxo.vn/laxo/laxo/shop"
)

type testHandlerService struct {
  lazada *lazada.Service
  shop   *shop.Service
  assets *assets.Service
}

type testHandler struct {
  server *laxo.Server
  service *testHandlerService
}


func InitTestHandler(server *laxo.Server, l *lazada.Service, p *shop.Service, a *assets.Service, r *mux.Router, n *negroni.Negroni) {
  s := &testHandlerService{
    lazada: l,
    shop: p,
    assets: a,
  }

  h := testHandler{
    server: server,
		service: s,
	}

	r.Handle("/test/test", n.With(
		negroni.WrapFunc(h.server.Middleware.AssureAuth(h.HandleTest)),
	)).Methods("GET")
}

func (h *testHandler) HandleTest(w http.ResponseWriter, r *http.Request, uID string) {
  w.Header().Set("Content-Type", "application/json; charset=utf-8")
  fmt.Fprintf(w, "Hello, testing: %s\n", "1, 2, 3")
}
