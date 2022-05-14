package rest

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
	"laxo.vn/laxo/laxo"
	"laxo.vn/laxo/laxo/lazada"
)

type testHandler struct {
  service *lazada.Service
}


func InitTestHandler(s *lazada.Service, r *mux.Router, n *negroni.Negroni) {
  h := testHandler{
		service: s,
	}

	r.Handle("/test/lazada", n.With(
		negroni.WrapFunc(laxo.AssureAuth(h.TestLazada)),
	)).Methods("GET")
}

func (h *testHandler) TestLazada(w http.ResponseWriter, r *http.Request, uID string) {
  key, err := h.service.RetrieveProductToRedis("product_lazada_test", 1)
  if err != nil {
    laxo.Logger.Error("fetch error", "error", err)
    http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
    return
  }
  fmt.Fprintf(w, "Hello, your redis key is: %s\n", key.Attributes.Name.String)
}
