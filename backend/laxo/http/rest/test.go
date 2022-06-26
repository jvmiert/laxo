package rest

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
	"laxo.vn/laxo/laxo"
	"laxo.vn/laxo/laxo/assets"
	"laxo.vn/laxo/laxo/lazada"
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
	_, err := h.service.shop.GetActiveShopByUserID(uID)
	if err != nil {
		h.server.Logger.Errorw("GetActiveShopByUserID returned error",
			"error", err,
		)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	test := `<div style="margin:0"><span style="font-family:none"><strong style="font-weight:bold;font-family:none">Test</strong></span></div><div style="margin:0"><span style="font-family:none"></span></div><div style="margin:0"><span style="font-family:none"><em>Test</em></span></div><div style="margin:0"><span style="font-family:none"></span></div><div style="margin:0;text-align:right;display:inline-block;width:100%"><u>Test</u></div><i><div style="margin:0;text-align:right;display:inline-block;width:100%">Test</div></i><div style="margin:0;text-align:right;display:inline-block;width:100%"> test</div><div style="margin:0"><span></span></div><div style="margin:0"><span></span></div><span><div style="margin:0;text-align:center;display:inline-block;width:100%">This is bold</div></span><div style="margin:0"><span></span></div><div style="margin:0"><span><u>This is underlined</u></span></div><div style="margin:0"><span></span></div><div style="margin:0"><span><em>This is italic</em></span></div><div style="margin:0"><span></span></div><div style="margin:0"><span></span></div><div style="margin:0"><h1><span>Heading1</span></h1></div><div style="margin:0"><h2><span>Heading2</span></h2></div><div style="margin:0"><h3><span>Heading3</span></h3></div><div style="margin:0"><span></span></div>`

	schema, _ := h.service.shop.HTMLToSlate(test, "01G1FZCVYH9J47DB2HZENSBC6E")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write([]byte(schema))
	return

	rows, _ := h.server.PglClient.Query(context.Background(), "select description from products_attribute_lazada")
	defer rows.Close()

	for rows.Next() {
		values, _ := rows.Values()

		if values[0] == nil {
			continue
		}

		h.server.Logger.Debugw("start schema",
			"data", values[0],
		)
		schema, _ := h.service.shop.HTMLToSlate(values[0].(string), "01G1FZCVYH9J47DB2HZENSBC6E")
		h.server.Logger.Debugw("finish schema",
			"schema", schema,
		)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write([]byte("test"))
}
