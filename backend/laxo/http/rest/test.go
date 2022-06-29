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

	_, err = h.service.shop.GetProductDetailsByID("01G4CES9QT2EPWGJW17DKYXAFS", s.Model.ID)
	if err != nil {
		h.server.Logger.Errorw("GetProductDetailsByID error",
			"error", err,
		)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	test := `<p style="margin:0"><span>This is left aligned</span></p><p style="text-align:center;display:inline-block;width:100%;margin:0"><span>This is center aligned</span></p><p style="text-align:right;display:inline-block;width:100%;margin:0"><span>This is right aligned</span></p><p style="margin:0"><span></span></p><p style="margin:0"><strong style="font-weight:bold">This is bold</strong></p><p style="margin:0"><em>This is italic</em></p><p style="margin:0"><strong style="font-weight:bold"><u>This is underlined and bold</u></strong></p><p style="margin:0"><span></span></p><ul style="list-style:disc;margin-left:10px"><li><span>This is bullet point one</span></li><li><span style="text-align:right;display:inline-block;width:100%">This is bullet point two, right aligned and <strong style="font-weight:bold">bold</strong></span></li><li><span>This is bullet point three</span></li></ul><div style="margin:0"><span></span></div><ol style="list-style:decimal"><li><span style="text-align:center;display:inline-block;width:100%">This is a center aligned and <u>underlined</u> numbered bullet point one</span></li><li><span>This is a <strong style="font-weight:bold">bold</strong> numbered bullet point two</span></li><li><span>This is an <em>italic</em> numbered bullet point three</span></li></ol>`

	schema, _ := h.service.shop.HTMLToSlate(test, "01G1FZCVYH9J47DB2HZENSBC6E")

	resp := &[]models.Element{}
	_ = json.Unmarshal([]byte(schema), resp)

	html, _ := h.service.shop.SlateToHTML(*resp)
	h.server.Logger.Debugw("SlateToHTML", "html", html)

	//err = h.service.lazada.UpdateProductToLazada(product, html)
	//if err != nil {
	//	h.server.Logger.Errorw("UpdateProductToLazada error",
	//		"error", err,
	//	)
	//	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	//	return
	//}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write([]byte(schema))
}
