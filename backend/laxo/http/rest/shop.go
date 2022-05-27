package rest

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
	"laxo.vn/laxo/laxo"
	"laxo.vn/laxo/laxo/lazada"
	"laxo.vn/laxo/laxo/shop"
	temporal_client "laxo.vn/laxo/temporal/client"
)

type shopHandlerService struct {
  lazada    *lazada.Service
  shop      *shop.Service
  temporal  *temporal_client.Client
}

type shopHandler struct {
  server *laxo.Server
  service *shopHandlerService
}

func InitShopHandler(server *laxo.Server, shop *shop.Service, l *lazada.Service,
                     r *mux.Router, n *negroni.Negroni, t *temporal_client.Client) {
  s := &shopHandlerService{
    lazada: l,
    shop: shop,
    temporal: t,
  }

  h := shopHandler{
    server: server,
		service: s,
	}

	r.Handle("/product", n.With(
		negroni.WrapFunc(h.server.Middleware.AssureAuth(h.GetProduct)),
	)).Methods("GET")

	r.Handle("/shop", n.With(
		negroni.WrapFunc(h.server.Middleware.AssureAuth(h.HandleGetMyShops)),
	)).Methods("GET")

	r.Handle("/oauth/verify", n.With(
		negroni.HandlerFunc(h.server.Middleware.AssureJSON),
		negroni.WrapFunc(h.server.Middleware.AssureAuth(h.HandleVerifyOAuth)),
	)).Methods("POST")

	r.Handle("/oauth/redirects", n.With(
		negroni.WrapFunc(h.server.Middleware.AssureAuth(h.HandleOAuthRedirects)),
	)).Methods("GET")

	r.Handle("/platforms/lazada", n.With(
		negroni.WrapFunc(h.server.Middleware.AssureAuth(h.HandleLazadaPlatformInfo)),
	)).Methods("GET")
}


func (h *shopHandler) HandleLazadaPlatformInfo(w http.ResponseWriter, r *http.Request, uID string) {
  shop, err := h.service.shop.GetActiveShopByUserID(uID)
  if err != nil {
    h.server.Logger.Errorw("GetActiveShopByUserID returned error",
      "error", err,
    )
    http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
    return
  }

  lazInfo, err := h.service.lazada.GetLazadaPlatformByShopID(shop.Model.ID)
  if err != nil {
    //@TODO: Handle empty return by returning 404 instead of  500
    h.server.Logger.Errorw("GetLazadaPlatformByShopID returned error",
      "error", err,
    )
    http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
    return
  }

  js, err := h.service.lazada.GetLazadaPlatformJSON(lazInfo)
  if err != nil {
    h.server.Logger.Errorw("GetLazadaPlatformJSON returned error",
      "error", err,
    )
    http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
    return
  }

  w.Header().Set("Content-Type", "application/json; charset=utf-8")
  w.Write(js)
}

func (h *shopHandler) HandleOAuthRedirects(w http.ResponseWriter, r *http.Request, uID string) {
  shopID := r.URL.Query().Get("shopID")
  o := &shop.OAuthRedirectRequest{ShopID: shopID}

  printer := laxo.GetLocalePrinter(r)
  if err := h.service.shop.ValidateOAuthRedirectRequest(o, uID, printer); err != nil {
    laxo.ErrorJSONEncode(w, err, http.StatusUnprocessableEntity)
    return
  }

  err := h.service.shop.GenerateRedirect(o)

  if err != nil {
    http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
    return
  }

  js, err := o.JSON()

  if err != nil {
    http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
    return
  }

  w.Header().Set("Content-Type", "application/json; charset=utf-8")
  w.Write(js)
}

func (h *shopHandler) HandleVerifyOAuth(w http.ResponseWriter, r *http.Request, uID string) {
  var o shop.OAuthVerifyRequest

  if err := laxo.DecodeJSONBody(h.server.Logger, w, r, &o); err != nil {
    var mr *laxo.MalformedRequest
    if errors.As(err, &mr) {
      http.Error(w, mr.Msg, mr.Status)
    } else {
      http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
    }
    return
  }

  printer := laxo.GetLocalePrinter(r)
  if err := h.service.shop.ValidateOAuthVerifyRequest(o, uID, printer); err != nil {
    laxo.ErrorJSONEncode(w, err, http.StatusUnprocessableEntity)
    return
  }

  fmt.Fprint(w, "OK")
}

func (h *shopHandler) GetProduct(w http.ResponseWriter, r *http.Request, uID string) {
  offset := r.URL.Query().Get("offset")
  limit := r.URL.Query().Get("limit")

  if offset == "" {
    offset = "0"
  }

  if limit == "" {
    limit = "50"
  }

  products, paginate, err := h.service.shop.GetProductsByUserID(uID, offset, limit)
  if err != nil {
    h.server.Logger.Errorw("GetProductsByUserID error",
      "error", err,
    )
    http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
    return
  }

  js, err := h.service.shop.GetProductListJSON(products, &paginate)
  if err != nil {
    h.server.Logger.Errorw("GetProductListJSON error",
      "error", err,
    )
    http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
    return
  }

  w.Header().Set("Content-Type", "application/json; charset=utf-8")
  w.Write(js)
}

func (h *shopHandler) HandleCreateShop(w http.ResponseWriter, r *http.Request, uID string) {
  var s shop.Shop

  if err := laxo.DecodeJSONBody(h.server.Logger, w, r, &s.Model); err != nil {
    var mr *laxo.MalformedRequest
    if errors.As(err, &mr) {
      http.Error(w, mr.Msg, mr.Status)
    } else {
      http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
    }
    return
  }

  printer := laxo.GetLocalePrinter(r)
  if err := h.service.shop.ValidateNewShop(&s, printer); err != nil {
    laxo.ErrorJSONEncode(w, err, http.StatusUnprocessableEntity)
    return
  }

  if err := h.service.shop.SaveNewShopToDB(&s, uID); err != nil {
    http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
    return
  }

  js, err := s.JSON()

  if err != nil {
    http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
    return
  }

  w.Header().Set("Content-Type", "application/json; charset=utf-8")
  w.Write(js)
}

func (h *shopHandler) HandleGetMyShops(w http.ResponseWriter, r *http.Request, uID string) {
  shops, err := h.service.shop.RetrieveShopsPlatformsByUserID(uID)

  if err != nil {
    laxo.ErrorJSONEncode(w, err, http.StatusInternalServerError)
    return
  }

  js, err := h.service.shop.GenerateShopPlatformList(shops)
  if err != nil {
    http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
    return
  }

  w.Header().Set("Content-Type", "application/json; charset=utf-8")
  w.Write(js)
}
