package rest

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
	"laxo.vn/laxo/laxo"
	"laxo.vn/laxo/laxo/assets"
	"laxo.vn/laxo/laxo/lazada"
	"laxo.vn/laxo/laxo/models"
	"laxo.vn/laxo/laxo/shop"
)

type testHandlerService struct {
	lazada *lazada.Service
	shop   *shop.Service
	assets *assets.Service
}

type testHandler struct {
	server  *laxo.Server
	service *testHandlerService
}

func InitTestHandler(server *laxo.Server, l *lazada.Service, p *shop.Service, a *assets.Service, r *mux.Router, n *negroni.Negroni) {
	s := &testHandlerService{
		lazada: l,
		shop:   p,
		assets: a,
	}

	h := testHandler{
		server:  server,
		service: s,
	}

	r.Handle("/test/test", n.With(
		negroni.WrapFunc(h.server.Middleware.AssureAuth(h.HandleTest)),
	)).Methods("GET")

	r.Handle("/test/upload", n.With(
		negroni.WrapFunc(h.server.Middleware.AssureAuth(h.HandleTestUpload)),
	)).Methods("GET")
}

func (h *testHandler) HandleTestUpload(w http.ResponseWriter, r *http.Request, uID string) {
	s, err := h.service.shop.GetActiveShopByUserID(uID)
	if err != nil {
		h.server.Logger.Errorw("GetActiveShopByUserID returned error",
			"error", err,
		)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	productID := "01G4CES9QT2EPWGJW17DKYXAFS"

	p, err := h.service.shop.GetProductDetailsByID(productID, s.Model.ID)
	if err != nil {
		h.server.Logger.Errorw("GetProductDetailsByID error",
			"error", err,
		)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	slateSchema := []models.Element{}

	err = json.Unmarshal([]byte(p.Model.DescriptionSlate.String), &slateSchema)
	if err != nil {
		h.server.Logger.Errorw("slate Unmarshal error",
			"error", err,
		)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	html, err := h.service.shop.SlateToHTML(slateSchema)
	if err != nil {
		h.server.Logger.Errorw("SlateToHTML error",
			"error", err,
		)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	parsedHTML, err := h.service.lazada.HandleDescriptionImages(html, s.Model.ID, s.Model.AssetsToken)
	if err != nil {
		h.server.Logger.Errorw("HandleDescriptionImages error",
			"error", err,
		)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	err = h.service.lazada.UpdateProductToLazada(p, parsedHTML)
	if err != nil {
		h.server.Logger.Errorw("UpdateProductToLazada error",
			"error", err,
		)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write([]byte(parsedHTML))
}

func (h *testHandler) HandleTest(w http.ResponseWriter, r *http.Request, uID string) {
	s, err := h.service.shop.GetActiveShopByUserID(uID)
	if err != nil {
		h.server.Logger.Errorw("GetActiveShopByUserID returned error",
			"error", err,
		)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	productID := "01G4CES9QT2EPWGJW17DKYXAFS"

	p, err := h.service.shop.GetProductDetailsByID(productID, s.Model.ID)
	if err != nil {
		h.server.Logger.Errorw("GetProductDetailsByID error",
			"error", err,
		)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	slateSchema := []models.Element{}

	err = json.Unmarshal([]byte(p.Model.DescriptionSlate.String), &slateSchema)
	if err != nil {
		h.server.Logger.Errorw("slate Unmarshal error",
			"error", err,
		)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	html, err := h.service.shop.SlateToHTML(slateSchema)
	if err != nil {
		h.server.Logger.Errorw("SlateToHTML error",
			"error", err,
		)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write([]byte(html))
}
