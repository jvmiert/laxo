package rest

import (
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
  products, err := h.service.GetProductsByUserID(uID)
  if err != nil {
    laxo.Logger.Error("GetProductsByUserID error", "error", err)
    laxo.ErrorJSONEncode(w, err, http.StatusUnprocessableEntity)
    return
  }

  js, err := h.service.GetProductListJSON(products)
  if err != nil {
    laxo.Logger.Error("GetProductListJSON error", "error", err)
    http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
    return
  }

  w.Header().Set("Content-Type", "application/json; charset=utf-8")
  w.Write(js)
}
