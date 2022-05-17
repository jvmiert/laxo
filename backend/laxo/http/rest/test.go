package rest

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jackc/pgtype"
	"github.com/urfave/negroni"
	"gopkg.in/guregu/null.v4"
	"laxo.vn/laxo/laxo"
	"laxo.vn/laxo/laxo/assets"
	"laxo.vn/laxo/laxo/lazada"
	"laxo.vn/laxo/laxo/product"
	"laxo.vn/laxo/laxo/sqlc"
)

type testHandlerService struct {
  lazada *lazada.Service
  product *product.Service
  assets  *assets.Service
}

type testHandler struct {
  service *testHandlerService
}


func InitTestHandler(l *lazada.Service, p *product.Service, a *assets.Service, r *mux.Router, n *negroni.Negroni) {
  s := &testHandlerService{
    lazada: l,
    product: p,
    assets: a,
  }

  h := testHandler{
		service: s,
	}

	r.Handle("/test/lazada", n.With(
		negroni.WrapFunc(laxo.AssureAuth(h.TestLazada)),
	)).Methods("GET")
}

func (h *testHandler) TestLazada(w http.ResponseWriter, r *http.Request, uID string) {
  p, err := h.service.lazada.RetrieveProductFromRedis("product_lazada_test", 1)
  if err != nil {
    laxo.Logger.Error("RetrieveProductFromRedis error", "error", err)
    http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
    return
  }

  pModel, pModelAttributes, pModelSKU, err := h.service.lazada.SaveOrUpdateLazadaProduct(p, "01G1FZCVYH9J47DB2HZENSBC6E")
  if err != nil {
    laxo.Logger.Error("SaveOrUpdateLazadaProduct error", "error", err)
    http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
    return
  }

  numericPrice := pgtype.Numeric{}
  numericPrice.Set(pModelSKU.Price.String)

  sanitzedDescription := h.service.lazada.GetSantizedDescription(pModelAttributes.Description.String)

  product := &product.Product{
    Model: &sqlc.Product{
      ID: "",
      Name: pModelAttributes.Name,
      Description : null.StringFrom(sanitzedDescription),
      Msku: null.StringFrom(pModelSKU.SellerSku),
      SellingPrice: numericPrice,
      CostPrice: pgtype.Numeric{Status: pgtype.Null},
      ShopID: pModel.ShopID,
    },
  }

  laxoP, err := h.service.product.SaveOrUpdateProductToStore(product, "01G1FZCVYH9J47DB2HZENSBC6E", pModel.ID)
  if err != nil {
    laxo.Logger.Error("SaveOrUpdateProductToStore error", "error", err)
    http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
    return
  }

  images, err := h.service.assets.ExtractImagesListFromProductResponse(p)
  if err != nil {
    laxo.Logger.Error("ExtractImagesListFromProductResponse error", "error", err)
    http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
    return
  }

  err = h.service.assets.SaveProductImages(images, "01G1FZCVYH9J47DB2HZENSBC6E", laxoP.Model.ID)
  if err != nil {
    laxo.Logger.Error("SaveProductImages error", "error", err)
    http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
    return
  }

  fmt.Fprintf(w, "Hello, your images are: %s\n", images)
}
