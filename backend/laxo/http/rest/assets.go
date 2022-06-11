package rest

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
	"laxo.vn/laxo/laxo"
	"laxo.vn/laxo/laxo/assets"
	"laxo.vn/laxo/laxo/shop"
)

type assetsHandler struct {
  server *laxo.Server
  shop   *shop.Service
  assets *assets.Service
}

func InitAssetsHandler(server *laxo.Server, shop *shop.Service, assets *assets.Service, r *mux.Router, n *negroni.Negroni) {
  h := assetsHandler{
    server: server,
		shop: shop,
		assets: assets,
	}

	r.Handle("/manage-assets", n.With(
		negroni.WrapFunc(h.server.Middleware.AssureAuth(h.HandlePutAsset)),
	)).Methods("POST")
}


func (h *assetsHandler) HandlePutAsset(w http.ResponseWriter, r *http.Request, uID string) {
  contentType := r.Header.Get("Content-type")

  // @TODO: use -> http.MaxBytesReader(w, r.Body, MaxSize)
   b, err := ioutil.ReadAll(r.Body)
   if err != nil {
    h.server.Logger.Errorw("ioutil.ReadAll returned error",
      "error", err,
    )
    http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
    return
   }

  h.server.Logger.Debugw("HandlePutAsset", "contentType", contentType, "length", len(b))

  w.Header().Set("Content-Type", "application/json; charset=utf-8")
  fmt.Fprintf(w, "Hello, testing: %s\n", "1, 2, 3")
}
