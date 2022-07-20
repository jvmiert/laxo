package rest

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
	"gopkg.in/guregu/null.v4"
	"laxo.vn/laxo/laxo"
	"laxo.vn/laxo/laxo/lazada"
	"laxo.vn/laxo/laxo/models"
	"laxo.vn/laxo/laxo/shop"
	temporal_client "laxo.vn/laxo/temporal/client"
)

type shopHandlerService struct {
	lazada   *lazada.Service
	shop     *shop.Service
	temporal *temporal_client.Client
}

type shopHandler struct {
	server  *laxo.Server
	service *shopHandlerService
}

func InitShopHandler(server *laxo.Server, shop *shop.Service, l *lazada.Service,
	r *mux.Router, n *negroni.Negroni, t *temporal_client.Client) {
	s := &shopHandlerService{
		lazada:   l,
		shop:     shop,
		temporal: t,
	}

	h := shopHandler{
		server:  server,
		service: s,
	}

	r.Handle("/product/{productID:[0123456789ABCDEFGHJKMNPQRSTVWXYZ]{26}}", n.With(
		negroni.WrapFunc(h.server.Middleware.AssureAuth(h.HandleProductDetails)),
	)).Methods("GET")

	r.Handle("/product/{productID:[0123456789ABCDEFGHJKMNPQRSTVWXYZ]{26}}", n.With(
		negroni.HandlerFunc(h.server.Middleware.AssureJSON),
		negroni.WrapFunc(h.server.Middleware.AssureAuth(h.HandlePostProductDetails)),
	)).Methods("POST")

	r.Handle("/product", n.With(
		negroni.WrapFunc(h.server.Middleware.AssureAuth(h.GetProduct)),
	)).Methods("GET")

	r.Handle("/product", n.With(
		negroni.HandlerFunc(h.server.Middleware.AssureJSON),
		negroni.WrapFunc(h.server.Middleware.AssureAuth(h.CreateProduct)),
	)).Methods("POST")

	r.Handle("/change-platform-sync/{productID:[0123456789ABCDEFGHJKMNPQRSTVWXYZ]{26}}", n.With(
		negroni.HandlerFunc(h.server.Middleware.AssureJSON),
		negroni.WrapFunc(h.server.Middleware.AssureAuth(h.HandleChangePlatformSync)),
	)).Methods("POST")

	r.Handle("/change-image-order/{productID:[0123456789ABCDEFGHJKMNPQRSTVWXYZ]{26}}", n.With(
		negroni.HandlerFunc(h.server.Middleware.AssureJSON),
		negroni.WrapFunc(h.server.Middleware.AssureAuth(h.HandleChangeImageOrder)),
	)).Methods("POST")

	r.Handle("/shop", n.With(
		negroni.WrapFunc(h.server.Middleware.AssureAuth(h.HandleGetMyShops)),
	)).Methods("GET")

	r.Handle("/shop", n.With(
		negroni.HandlerFunc(h.server.Middleware.AssureJSON),
		negroni.WrapFunc(h.server.Middleware.AssureAuth(h.HandleCreateShop)),
	)).Methods("POST")

	r.Handle("/oauth/verify", n.With(
		negroni.HandlerFunc(h.server.Middleware.AssureJSON),
		negroni.WrapFunc(h.server.Middleware.AssureAuth(h.HandleVerifyOAuth)),
	)).Methods("POST")

	r.Handle("/oauth/redirects", n.With(
		negroni.WrapFunc(h.server.Middleware.AssureAuth(h.HandleOAuthRedirects)),
	)).Methods("GET")

	r.Handle("/platform-sync", n.With(
		negroni.WrapFunc(h.server.Middleware.AssureAuth(h.HandlePlatformSync)),
	)).Methods("POST")

	r.Handle("/platforms/lazada", n.With(
		negroni.WrapFunc(h.server.Middleware.AssureAuth(h.HandleLazadaPlatformInfo)),
	)).Methods("GET")
}

func (h *shopHandler) CreateProduct(w http.ResponseWriter, r *http.Request, uID string) {
	var p models.NewProductRequest

	if err := laxo.DecodeJSONBody(h.server.Logger, w, r, &p); err != nil {
		var mr *laxo.MalformedRequest
		if errors.As(err, &mr) {
			http.Error(w, mr.Msg, mr.Status)
		} else {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	s, err := h.service.shop.GetActiveShopByUserID(uID)
	if err != nil {
		h.server.Logger.Errorw("GetActiveShopByUserID returned error",
			"error", err,
		)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	err = h.service.shop.ValidateNewProductRequest(&p, s.Model.ID)
	if err != nil {
		laxo.OzzoErrorJSONEncode(w, err, http.StatusUnprocessableEntity, h.server.Logger)
		return
	}

	pModel, err := h.service.shop.CreateNewProduct(&p, s.Model.ID)
	if err != nil {
		h.server.Logger.Errorw("CreateNewProduct returned error",
			"error", err,
		)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	productDetails, err := h.service.shop.GetProductDetailsByID(pModel.ID, s.Model.ID)
	if err != nil {
		h.server.Logger.Errorw("GetProductDetailsByID returned error",
			"error", err,
		)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	b, err := productDetails.JSON()
	if err != nil {
		h.server.Logger.Errorw("product JSON marshall error",
			"error", err,
		)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(b)
}

func (h *shopHandler) HandleChangeImageOrder(w http.ResponseWriter, r *http.Request, uID string) {
	vars := mux.Vars(r)
	productID := vars["productID"]

	s, err := h.service.shop.GetActiveShopByUserID(uID)
	if err != nil {
		h.server.Logger.Errorw("GetActiveShopByUserID returned error",
			"error", err,
		)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	err = h.service.shop.ProductIsOwnedByStore(productID, s.Model.ID)
	if err != nil && err != shop.ErrProductNotOwned {
		h.server.Logger.Errorw("ProductIsOwnedByStore returned error",
			"error", err,
		)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if err == shop.ErrProductNotOwned {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

	var p models.ProductImageOrderRequest

	if err = laxo.DecodeJSONBody(h.server.Logger, w, r, &p); err != nil {
		var mr *laxo.MalformedRequest
		if errors.As(err, &mr) {
			http.Error(w, mr.Msg, mr.Status)
		} else {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	err = h.service.shop.ValidateProductImageOrderRequest(&p)
	if err != nil {
		laxo.OzzoErrorJSONEncode(w, err, http.StatusUnprocessableEntity, h.server.Logger)
		return
	}

	err = h.service.shop.UpdateProductImageOrderRequest(productID, &p)
	if err != nil {
		h.server.Logger.Errorw("UpdateProductImageOrderRequest returned error",
			"error", err,
		)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write([]byte("{\"success\":true}"))
}

func (h *shopHandler) HandleChangePlatformSync(w http.ResponseWriter, r *http.Request, uID string) {
	vars := mux.Vars(r)
	productID := vars["productID"]

	s, err := h.service.shop.GetActiveShopByUserID(uID)
	if err != nil {
		h.server.Logger.Errorw("GetActiveShopByUserID returned error",
			"error", err,
		)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	err = h.service.shop.ProductIsOwnedByStore(productID, s.Model.ID)
	if err != nil && err != shop.ErrProductNotOwned {
		h.server.Logger.Errorw("ProductIsOwnedByStore returned error",
			"error", err,
		)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if err == shop.ErrProductNotOwned {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

	var p models.ProductChangedSyncRequest

	if err = laxo.DecodeJSONBody(h.server.Logger, w, r, &p); err != nil {
		var mr *laxo.MalformedRequest
		if errors.As(err, &mr) {
			http.Error(w, mr.Msg, mr.Status)
		} else {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	err = h.service.shop.ChangedProductPlatformSync(&p, productID)
	if err != nil {
		h.server.Logger.Errorw("ChangedProductPlatformSync returned error",
			"error", err,
		)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write([]byte("{\"success\":true}"))
}

func (h *shopHandler) HandlePostProductDetails(w http.ResponseWriter, r *http.Request, uID string) {
	var p models.ProductDetailPostRequest

	if err := laxo.DecodeJSONBody(h.server.Logger, w, r, &p); err != nil {
		var mr *laxo.MalformedRequest
		if errors.As(err, &mr) {
			http.Error(w, mr.Msg, mr.Status)
		} else {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	printer := laxo.GetLocalePrinter(r)

	err := h.service.shop.ValidateProductDetails(&p, printer)
	if err != nil {
		laxo.OzzoErrorJSONEncode(w, err, http.StatusUnprocessableEntity, h.server.Logger)
		return
	}

	s, err := h.service.shop.GetActiveShopByUserID(uID)
	if err != nil {
		h.server.Logger.Errorw("GetActiveShopByUserID returned error",
			"error", err,
		)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	vars := mux.Vars(r)
	productID := vars["productID"]

	err = h.service.shop.UpdateProductFromRequest(&p, printer, s.Model.ID, productID)
	if err != nil {
		laxo.OzzoErrorJSONEncode(w, err, http.StatusUnprocessableEntity, h.server.Logger)
		return
	}

	//@TODO: handle product sync

	product, err := h.service.shop.GetProductDetailsByID(productID, s.Model.ID)
	if errors.Is(err, shop.ErrProductNotFound) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	if err != nil {
		h.server.Logger.Errorw("GetProductDetailsByID error",
			"error", err,
		)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	b, err := product.JSON()
	if err != nil {
		h.server.Logger.Errorw("product JSON marshall error",
			"error", err,
		)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(b)
}

func (h *shopHandler) HandleProductDetails(w http.ResponseWriter, r *http.Request, uID string) {
	s, err := h.service.shop.GetActiveShopByUserID(uID)
	if err != nil {
		h.server.Logger.Errorw("GetActiveShopByUserID returned error",
			"error", err,
		)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	vars := mux.Vars(r)
	productID := vars["productID"]

	product, err := h.service.shop.GetProductDetailsByID(productID, s.Model.ID)
	if errors.Is(err, shop.ErrProductNotFound) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	if err != nil {
		h.server.Logger.Errorw("GetProductDetailsByID error",
			"error", err,
		)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	b, err := product.JSON()
	if err != nil {
		h.server.Logger.Errorw("product JSON marshall error",
			"error", err,
		)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(b)
}

func (h *shopHandler) HandlePlatformSync(w http.ResponseWriter, r *http.Request, uID string) {
	shop, err := h.service.shop.GetActiveShopByUserID(uID)
	if err != nil {
		h.server.Logger.Errorw("GetActiveShopByUserID returned error",
			"error", err,
		)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	b, err := shop.JSON()
	if err != nil {
		h.server.Logger.Errorw("shop JSON returned error",
			"error", err,
		)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	_, err = h.service.temporal.StartLazadaPlatformSync(shop.Model.ID, uID, true)
	if err != nil {
		h.server.Logger.Errorw("StartLazadaPlatformSync returned error",
			"error", err,
		)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(b)
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
		laxo.OzzoErrorJSONEncode(w, err, http.StatusUnprocessableEntity, h.server.Logger)
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
		laxo.OzzoErrorJSONEncode(w, err, http.StatusUnprocessableEntity, h.server.Logger)
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

	name := r.URL.Query().Get("name")
	msku := r.URL.Query().Get("msku")

	nameParsed := null.StringFrom(name)
	mskuParsed := null.StringFrom(msku)

	if name == "" {
		nameParsed = null.NewString("", false)
	}

	if msku == "" {
		mskuParsed = null.NewString("", false)
	}

	var err error
	var products []models.Product
	var paginate models.Paginate

	if nameParsed.Valid || mskuParsed.Valid {
		products, paginate, err = h.service.shop.GetProductsByNameOrSKU(uID, nameParsed, mskuParsed, offset, limit)
	} else {
		products, paginate, err = h.service.shop.GetProductsByUserID(uID, offset, limit)
	}
	if err != nil {
		h.server.Logger.Errorw("GetProducts error",
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
		laxo.OzzoErrorJSONEncode(w, err, http.StatusUnprocessableEntity, h.server.Logger)
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
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
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
