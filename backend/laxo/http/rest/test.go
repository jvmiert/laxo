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

	product, err := h.service.shop.GetProductDetailsByID("01G4CES9QT2EPWGJW17DKYXAFS", s.Model.ID)
	if err != nil {
		h.server.Logger.Errorw("GetProductDetailsByID error",
			"error", err,
		)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	test := `<div style="margin:0"><span style="font-family:none"><strong style="font-weight:bold;font-family:none">Test</strong></span></div><div style="margin:0"><span style="font-family:none"></span></div><div style="margin:0"><span style="font-family:none"><em>Test</em></span></div><div style="margin:0"><span style="font-family:none"></span></div><div style="margin:0;text-align:right;display:inline-block;width:100%"><u>Test</u></div><i><div style="margin:0;text-align:right;display:inline-block;width:100%">Test</div></i><div style="margin:0;text-align:right;display:inline-block;width:100%"> test</div><div style="margin:0"><span></span></div><div style="margin:0"><span></span></div><span><div style="margin:0;text-align:center;display:inline-block;width:100%">This is bold</div></span><div style="margin:0"><span></span></div><div style="margin:0"><span><strong><u>This is underlined</u></strong></span></div><div style="margin:0"><span></span></div><div style="margin:0"><span><em>This is italic</em></span></div><div style="margin:0"><span></span></div><div style="margin:0"><span></span></div><div style="margin:0"><h1><span>Heading1</span></h1></div><div style="margin:0"><h2><span>Heading2</span></h2></div><div style="margin:0"><h3><span>Heading3</span></h3></div><div style="margin:0"><span></span></div>`

	schema, _ := h.service.shop.HTMLToSlate(test, "01G1FZCVYH9J47DB2HZENSBC6E")

	resp := &[]models.Element{}
	_ = json.Unmarshal([]byte(schema), resp)

	html, _ := h.service.shop.SlateToHTML(*resp)
	h.server.Logger.Debugw("SlateToHTML", "html", html)

	err = h.service.lazada.UpdateProductToLazada(product, html)
	if err != nil {
		h.server.Logger.Errorw("UpdateProductToLazada error",
			"error", err,
		)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write([]byte(schema))
}
