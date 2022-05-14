package rest

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
	"laxo.vn/laxo/laxo"
	"laxo.vn/laxo/laxo/product"
)

type productHandler struct {
  service *product.Service
}


func InitProductHandler(s *product.Service, r *mux.Router, n *negroni.Negroni) {
  h := productHandler{
		service: s,
	}

	r.Handle("/product", n.With(
		negroni.WrapFunc(laxo.AssureAuth(h.GetProduct)),
	)).Methods("GET")
}

func (h *productHandler) GetProduct(w http.ResponseWriter, r *http.Request, uID string) {
  fmt.Fprintf(w, "Hello, your uID is: %s\n", uID)
}
