package rest

import (
	"errors"
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

  r.Handle("/asset/{assetID:[0123456789ABCDEFGHJKMNPQRSTVWXYZ]{26}}", n.With(
		negroni.WrapFunc(h.server.Middleware.AssureAuth(h.HandlePutAsset)),
	)).Methods("PUT")

	r.Handle("/asset/create", n.With(
		negroni.HandlerFunc(h.server.Middleware.AssureJSON),
		negroni.WrapFunc(h.server.Middleware.AssureAuth(h.HandleCreateAsset)),
	)).Methods("POST")


	r.Handle("/asset/assign-product", n.With(
		negroni.HandlerFunc(h.server.Middleware.AssureJSON),
		negroni.WrapFunc(h.server.Middleware.AssureAuth(h.HandleAssignProduct)),
	)).Methods("POST")
}

func (h *assetsHandler) HandleAssignProduct(w http.ResponseWriter, r *http.Request, uID string) {
  var a assets.AssignRequest

  if err := laxo.DecodeJSONBody(h.server.Logger, w, r, &a); err != nil {
    var mr *laxo.MalformedRequest
    if errors.As(err, &mr) {
      http.Error(w, mr.Msg, mr.Status)
    } else {
      http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
    }
    return
  }

  err := h.assets.ModifyAssetAssignment(&a)
  if err != nil {
    h.server.Logger.Errorw("ModifyAssetAssignment returned error",
      "error", err,
    )
    http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
    return
  }

  b, err := h.assets.ValidAssignReply()
  if err != nil {
    h.server.Logger.Errorw("ValidAssignReply returned error",
      "error", err,
    )
    http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
    return
  }

  w.Header().Set("Content-Type", "application/json; charset=utf-8")
  w.Write(b)
}

func (h *assetsHandler) HandleCreateAsset(w http.ResponseWriter, r *http.Request, uID string) {
  shop, err := h.shop.GetActiveShopByUserID(uID)
  if err != nil {
    h.server.Logger.Errorw("GetActiveShopByUserID returned error",
      "error", err,
    )
    http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
    return
  }

  var a assets.AssetRequest

  if err = laxo.DecodeJSONBody(h.server.Logger, w, r, &a); err != nil {
    var mr *laxo.MalformedRequest
    if errors.As(err, &mr) {
      http.Error(w, mr.Msg, mr.Status)
    } else {
      http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
    }
    return
  }

  err = h.assets.ValidateAssetExtension(a.OriginalName)
  if err != nil {
    h.server.Logger.Errorw("GetOrCreateAsset returned error",
      "error", err,
    )
    http.Error(w, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
    return
  }

  reply, err := h.assets.GetOrCreateAsset(a, shop.Model.ID)
  if err != nil {
    h.server.Logger.Errorw("GetOrCreateAsset returned error",
      "error", err,
    )
    http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
    return
  }

  b, err := reply.JSON()
  if err != nil {
    h.server.Logger.Errorw("AssetReply marshal error",
      "error", err,
    )
    http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
    return
  }

  w.Header().Set("Content-Type", "application/json; charset=utf-8")
  w.Write(b)
}

func (h *assetsHandler) HandlePutAsset(w http.ResponseWriter, r *http.Request, uID string) {
  shop, err := h.shop.GetActiveShopByUserID(uID)
  if err != nil {
    h.server.Logger.Errorw("GetActiveShopByUserID returned error",
      "error", err,
    )
    http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
    return
  }

  vars := mux.Vars(r)
  assetID := vars["assetID"]

  b, err := h.assets.ExtractImageFromRequest(w, r.Body)
  if err != nil {
    h.server.Logger.Errorw("ExtractImageFromRequest returned error",
      "error", err,
    )
    http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
    return
  }

  hex := h.assets.GetMurmurFromBytes(b)

  h.assets.ValidateAssetHash(hex, assetID)
  if err != nil {
    h.server.Logger.Errorw("ValidateAssetHash returned error",
      "error", err,
    )
    http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
    return
  }

  asset, err := h.assets.SaveAssetToDisk(b, assetID, shop.Model.AssetsToken)
  if err != nil {
    h.server.Logger.Errorw("SaveAssetToDisk returned error",
      "error", err,
    )
    http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
    return
  }

  js, err := h.assets.AssetJSON(asset)
  if err != nil {
    h.server.Logger.Errorw("asset marshal returned error",
      "error", err,
    )
    http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
    return
  }

  w.Header().Set("Content-Type", "application/json; charset=utf-8")
  w.Write(js)
}
