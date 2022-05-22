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
	"laxo.vn/laxo/laxo/shop"
	"laxo.vn/laxo/laxo/sqlc"
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

	r.Handle("/test/lazada", n.With(
		negroni.WrapFunc(h.server.Middleware.AssureAuth(h.TestLazada)),
	)).Methods("GET")
}

func (h *testHandler) TestLazada(w http.ResponseWriter, r *http.Request, uID string) {
  p, err := h.service.lazada.RetrieveProductFromRedis("product_lazada_test", 1)
  if err != nil {
    h.server.Logger.Errorw("RetrieveProductFromRedis error",
      "error", err,
    )
    http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
    return
  }

  pModel, pModelAttributes, pModelSKU, err := h.service.lazada.SaveOrUpdateLazadaProduct(p, "01G1FZCVYH9J47DB2HZENSBC6E")
  if err != nil {
    h.server.Logger.Errorw("SaveOrUpdateLazadaProduct error",
      "error", err,
    )
    http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
    return
  }

  numericPrice := pgtype.Numeric{}
  numericPrice.Set(pModelSKU.Price.String)

  sanitzedDescription := h.service.lazada.GetSantizedDescription(pModelAttributes.Description.String)

  product := &shop.Product{
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

  laxoP, err := h.service.shop.SaveOrUpdateProductToStore(product, "01G1FZCVYH9J47DB2HZENSBC6E", pModel.ID)
  if err != nil {
    h.server.Logger.Errorw("SaveOrUpdateProductToStore error",
      "error", err,
    )
    http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
    return
  }

  images, err := h.service.assets.ExtractImagesListFromProductResponse(p)
  if err != nil {
    h.server.Logger.Errorw("ExtractImagesListFromProductResponse error",
      "error", err,
    )
    http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
    return
  }

  err = h.service.assets.SaveProductImages(images, "01G1FZCVYH9J47DB2HZENSBC6E", laxoP.Model.ID)
  if err != nil {
    h.server.Logger.Errorw("SaveProductImages error",
      "error", err,
    )
    http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
    return
  }

  fmt.Fprintf(w, "Hello, your images are: %s\n", images)
}
