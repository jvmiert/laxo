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
	s, err := h.service.shop.GetActiveShopByUserID(uID)
	if err != nil {
		h.server.Logger.Errorw("GetActiveShopByUserID returned error",
			"error", err,
		)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	test := `<p style="margin: 0;padding: 8.0px 0;white-space: pre-wrap;"></p><div style="width: 100.0%;"><img src="https://vn-live-02.slatic.net/p/436f973ccf5d8ab22b6e0638ba4c16c2.png" style="width: 100.0%;display: block;"/></div><div style="width: 100.0%;"><img src="https://vn-live-02.slatic.net/p/5a878340f3d6fcbd2286aba8ce992d37.png" style="width: 100.0%;display: block;"/></div><div style="width: 100.0%;"><img src="https://vn-live-02.slatic.net/p/3430c6bf503cecd6874d03783441800e.jpg" style="width: 100.0%;display: block;"/></div><div style="width: 100.0%;margin: 0;padding: 8.0px 0;white-space: pre-wrap;"><div style="width: 100.0%;margin: 0;padding: 8.0px 0;white-space: pre-wrap;"><i></i></div></div>`

	h.service.assets.ExtractImagesFromDescription(test, s.Model.ID, s.Model.AssetsToken)

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
