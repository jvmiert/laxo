package rest

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
	"laxo.vn/laxo/laxo"
	"laxo.vn/laxo/laxo/shop"
)

type shopHandler struct {
  server *laxo.Server
  service *shop.Service
}

func InitProductHandler(server *laxo.Server, s *shop.Service, r *mux.Router, n *negroni.Negroni) {
  h := shopHandler{
    server: server,
		service: s,
	}

	r.Handle("/product", n.With(
		negroni.WrapFunc(laxo.AssureAuth(h.GetProduct)),
	)).Methods("GET")

	r.Handle("/shop", n.With(
		negroni.WrapFunc(laxo.AssureAuth(h.HandleGetMyShops)),
	)).Methods("GET")

	r.Handle("/oauth/verify", n.With(
		negroni.HandlerFunc(laxo.AssureJSON),
		negroni.WrapFunc(laxo.AssureAuth(h.HandleVerifyOAuth)),
	)).Methods("POST")

	r.Handle("/oauth/redirects", n.With(
		negroni.WrapFunc(laxo.AssureAuth(h.HandleOAuthRedirects)),
	)).Methods("GET")
}

func (h *shopHandler) HandleOAuthRedirects(w http.ResponseWriter, r *http.Request, uID string) {
  shopID := r.URL.Query().Get("shopID")
  o := &shop.OAuthRedirectRequest{ShopID: shopID}

  printer := laxo.GetLocalePrinter(r)
  if err := h.service.ValidateOAuthRedirectRequest(o, uID, printer); err != nil {
    laxo.ErrorJSONEncode(w, err, http.StatusUnprocessableEntity)
    return
  }

  err := h.service.GenerateRedirect(o)

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
  if err := h.service.ValidateOAuthVerifyRequest(o, uID, printer); err != nil {
    laxo.ErrorJSONEncode(w, err, http.StatusUnprocessableEntity)
    return
  }

  fmt.Fprint(w, "OK")
}

func (h *shopHandler) GetProduct(w http.ResponseWriter, r *http.Request, uID string) {
  products, err := h.service.GetProductsByUserID(uID)
  if err != nil {
    h.server.Logger.Errorw("GetProductsByUserID error",
      "error", err,
    )
    laxo.ErrorJSONEncode(w, err, http.StatusUnprocessableEntity)
    return
  }

  js, err := h.service.GetProductListJSON(products)
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
  if err := h.service.ValidateNewShop(&s, printer); err != nil {
    laxo.ErrorJSONEncode(w, err, http.StatusUnprocessableEntity)
    return
  }

  if err := h.service.SaveNewShopToDB(&s, uID); err != nil {
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
  shops, err := h.service.RetrieveShopsPlatformsByUserID(uID)

  if err != nil {
    laxo.ErrorJSONEncode(w, err, http.StatusInternalServerError)
    return
  }

  js, err := h.service.GenerateShopPlatformList(shops)
  if err != nil {
    http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
    return
  }

  w.Header().Set("Content-Type", "application/json; charset=utf-8")
  w.Write(js)
}
