package lazada

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/araddon/dateparse"
	"gopkg.in/guregu/null.v4"
	"laxo.vn/laxo/laxo"
	"laxo.vn/laxo/laxo/models"
)

var ErrPlatformFailed = errors.New("platform returned failed message")
var ErrProductsFailed = errors.New("get products returned failed message")
var ErrProductUpdateFailed = errors.New("product update returned failed")
var ErrProductsParseFailed = errors.New("couldn't parse products")

type APIResponse struct {
	Code      string `json:"code"`
	Type      string `json:"type"`
	Message   string `json:"message"`
	RequestID string `json:"request_id"`
}

type CountryUserInfo struct {
	Country   string `json:"country"`
	UserID    string `json:"user_id"`
	SellerID  string `json:"seller_id"`
	ShortCode string `json:"short_code"`
}

type AuthResponse struct {
	Code      string `json:"code"`
	Type      string `json:"type"`
	Message   string `json:"message"`
	RequestID string `json:"request_id"`

	AccessToken      string            `json:"access_token"`
	Country          string            `json:"country"`
	RefreshToken     string            `json:"refresh_token"`
	AccountPlatform  string            `json:"account_platform"`
	RefreshExpiresIn int               `json:"refresh_expires_in"`
	CountryUserInfo  []CountryUserInfo `json:"country_user_info"`
	ExpiresIn        int               `json:"expires_in"`
	Account          string            `json:"account"`

	DateAccessExpired  time.Time
	DateRefreshExpired time.Time
}

type ProductDescriptionChange struct {
	XMLName xml.Name `xml:"description"`
	Text    string   `xml:",cdata"`
}

type ProductsResponseAttributes struct {
	Name                 null.String `json:"name" xml:"name"`
	ShortDescription     null.String `json:"short_description" xml:"-"`
	Description          null.String `json:"description" xml:"-"`
	Brand                null.String `json:"brand" xml:"-"`
	Model                null.String `json:"model" xml:"-"`
	HeadphoneFeatures    null.String `json:"headphone_features" xml:"-"`
	Bluetooth            null.String `json:"bluetooth" xml:"-"`
	WarrantyType         null.String `json:"warranty_type" xml:"-"`
	Warranty             null.String `json:"warranty" xml:"-"`
	Hazmat               null.String `json:"Hazmat" xml:"-"`
	ExpireDate           null.String `json:"Expire_date" xml:"-"`
	BrandClassification  null.String `json:"brand_classification" xml:"-"`
	IngredientPreference null.String `json:"ingredient_preference" xml:"-"`
	LotNumber            null.String `json:"Lot_number" xml:"-"`
	UnitsHB              null.String `json:"units_hb" xml:"-"`
	FmltSkinCare         null.String `json:"fmlt_skin_care" xml:"-"`
	Quantitative         null.String `json:"Quantitative" xml:"-"`
	SkinCareByAge        null.String `json:"skin_care_by_age" xml:"-"`
	SkinBenefit          null.String `json:"skin_benefit" xml:"-"`
	SkinType             null.String `json:"skin_type" xml:"-"`
	UserManual           null.String `json:"User_Manual" xml:"-"`
	CountryOriginHB      null.String `json:"country_origin_hb" xml:"-"`
	ColorFamily          null.String `json:"color_family" xml:"-"`
	FragranceFamily      null.String `json:"fragrance_family" xml:"-"`
	Source               null.String `json:"source" xml:"-"`

	// This is used to create valid XML for updating the product details
	HTMLChangeDescription ProductDescriptionChange `json:"-"`
}

type ProductsResponseSuspendedSkus struct {
	RejectReason null.String `json:"rejectReason"`
	SellerSku    null.String `json:"SellerSku"`
	SkuID        null.Int    `json:"SkuId"`
}

type ProductsResponseSkus struct {
	Status             null.String   `json:"Status" xml:"-"`
	Quantity           null.Int      `json:"quantity" xml:"-"`
	Images             []null.String `json:"Images" xml:"-"`
	MarketImages       []null.String `json:"marketImages" xml:"-"`
	SellerSku          null.String   `json:"SellerSku" xml:"SellerSku"`
	ShopSku            null.String   `json:"ShopSku" xml:"-"`
	PackageContent     null.String   `json:"package_content" xml:"-"`
	URL                null.String   `json:"Url" xml:"-"`
	PackageWidth       null.String   `json:"package_width" xml:"-"`
	ColorFamily        null.String   `json:"color_family" xml:"-"`
	PackageHeight      null.String   `json:"package_height" xml:"-"`
	SpecialPrice       json.Number   `json:"special_price" xml:"-"`
	Price              json.Number   `json:"price" xml:"-"`
	PackageLength      null.String   `json:"package_length" xml:"-"`
	PackageWeight      null.String   `json:"package_weight" xml:"-"`
	Available          null.Int      `json:"Available" xml:"-"`
	SkuID              null.Int      `json:"SkuId" xml:"SkuId"`
	SpecialToTimeRaw   null.String   `json:"special_to_time" xml:"-"`
	SpecialFromTimeRaw null.String   `json:"special_from_time" xml:"-"`
	SpecialFromDateRaw null.String   `json:"special_from_date" xml:"-"`
	SpecialToDateRaw   null.String   `json:"special_to_date" xml:"-"`
	SpecialToTime      time.Time     `xml:"-"`
	SpecialFromTime    time.Time     `xml:"-"`
	SpecialFromDate    time.Time     `xml:"-"`
	SpecialToDate      time.Time     `xml:"-"`
}

func (p *ProductsResponseSkus) ParseTime() error {
	adjustValues := [4]null.String{
		p.SpecialToTimeRaw,
		p.SpecialFromTimeRaw,
		p.SpecialToDateRaw,
		p.SpecialFromDateRaw,
	}

	parseValues := [4]time.Time{
		p.SpecialToTime,
		p.SpecialFromTime,
		p.SpecialToDate,
		p.SpecialFromDate,
	}

	for i, v := range adjustValues {
		if v.Valid {
			t, err := dateparse.ParseStrict(v.String)
			if err != nil {
				return err
			}

			parseValues[i] = t
		}
	}

	return nil
}

type Request struct {
	Product ProductsResponseProducts
}

type ProductsResponseProducts struct {
	XMLName         xml.Name                        `json:"-" xml:"Product"`
	Skus            []ProductsResponseSkus          `json:"skus" xml:"Skus>Sku"`
	ItemID          int64                           `json:"item_id" xml:"ItemId"`
	PrimaryCategory int64                           `json:"primary_category" xml:"-"`
	Attributes      ProductsResponseAttributes      `json:"attributes" xml:"Attributes"`
	CreatedTimeRaw  null.String                     `json:"created_time" xml:"-"`
	UpdatedTimeRaw  null.String                     `json:"updated_time" xml:"-"`
	Images          []null.String                   `json:"images" xml:"-"`
	MarketImages    []null.String                   `json:"marketImages" xml:"-"`
	Status          null.String                     `json:"status" xml:"-"`
	SubStatus       null.String                     `json:"subStatus" xml:"-"`
	SuspendedSkus   []ProductsResponseSuspendedSkus `json:"suspendedSkus" xml:"-"`
	Created         time.Time                       `xml:"-"`
	Updated         time.Time                       `xml:"-"`
}

func (p *ProductsResponseProducts) ParseTime() error {
	if !p.CreatedTimeRaw.Valid || !p.UpdatedTimeRaw.Valid {
		return ErrProductsParseFailed
	}

	i, err := strconv.ParseInt(p.CreatedTimeRaw.String, 10, 64)
	if err != nil {
		return err
	}

	p.Created = time.Unix(i/1000, 0)

	i, err = strconv.ParseInt(p.UpdatedTimeRaw.String, 10, 64)
	if err != nil {
		return err
	}

	p.Updated = time.Unix(i/1000, 0)
	return nil
}

type ProductResponseVariation struct {
	Name      string   `json:"name"`
	HasImage  bool     `json:"has_image"`
	Customize bool     `json:"customize"`
	Options   []string `json:"options"`
	Label     string   `json:"label"`
}

type ProductResponseProducts struct {
	ProductsResponseProducts
	Variations []ProductResponseVariation `json:"varation"`
}

type ProductsResponseData struct {
	TotalProducts int                        `json:"total_products"`
	Products      []ProductsResponseProducts `json:"products"`
}

type ProductsResponse struct {
	Data      ProductsResponseData `json:"data"`
	Code      string               `json:"code"`
	RequestID string               `json:"request_id"`
	RawData   []byte               `json:"laxo_raw"`
}

type ProductResponse struct {
	Data      ProductResponseProducts `json:"data"`
	Code      string                  `json:"code"`
	RequestID string                  `json:"request_id"`
	RawData   []byte                  `json:"laxo_raw"`
}

type QueryProductsParams struct {
	Filter        string
	UpdateBefore  time.Time
	CreatedBefore time.Time
	Offset        int
	CreatedAfter  time.Time
	UpdatedAfter  time.Time
	Limit         int
}

type LazadaClient struct {
	APIKey    string
	APISecret string
	Region    string

	Logger *laxo.Logger

	Method     string
	SysParams  map[string]string
	APIParams  map[string]string
	FileParams map[string][]byte
}

func (lc *LazadaClient) UpdateProduct(p *models.ProductDetails, lazadaHTML string) error {
	var lazadaPlatform models.ProductPlatformInformation
	for _, v := range p.Platforms {
		if v.PlatformName == "lazada" {
			lazadaPlatform = v
		}
	}

	if lazadaPlatform.ID == "" {
		return errors.New("no lazada platform information supplied")
	}

	lazadaID, err := strconv.ParseInt(lazadaPlatform.ID, 10, 64)
	if err != nil {
		return fmt.Errorf("ParseInt: %w", err)
	}

	lazHTMLDescription := ProductDescriptionChange{
		Text: lazadaHTML,
	}

	lazadaAttributes := ProductsResponseAttributes{
		Name:                  p.Model.Name,
		HTMLChangeDescription: lazHTMLDescription,
	}

	SkuID, err := strconv.ParseInt(lazadaPlatform.PlatformSKU, 10, 64)
	if err != nil {
		return fmt.Errorf("strconv ParseInt: %w", err)
	}

	lazadaSKU := ProductsResponseSkus{
		SellerSku: null.StringFrom(lazadaPlatform.SellerSKU),
		SkuID:     null.IntFrom(SkuID),
	}

	lazadaProduct := ProductsResponseProducts{
		ItemID:     lazadaID,
		Attributes: lazadaAttributes,
		Skus:       []ProductsResponseSkus{lazadaSKU},
	}
	request := Request{
		Product: lazadaProduct,
	}

	out, err := xml.MarshalIndent(request, " ", "  ")
	if err != nil {
		return fmt.Errorf("xml MarshalIndent: %w", err)
	}
	xml := xml.Header + string(out)

	var req *http.Request

	lc.AddAPIParam("payload", xml)

	// add query params
	values := url.Values{}
	for key, val := range lc.SysParams {
		values.Add(key, val)
	}

	for key, val := range lc.APIParams {
		values.Add(key, val)
	}

	apiPath := "/product/update"
	apiServerURL := "https://api.lazada.vn/rest"

	values.Add("sign", lc.sign(apiPath))
	fullURL := fmt.Sprintf("%s%s?%s", apiServerURL, apiPath, values.Encode())

	lc.Logger.Debugw("UpdateProduct", "fullURL", fullURL, "xml", xml)

	req, err = http.NewRequest("POST", fullURL, nil)
	if err != nil {
		return fmt.Errorf("NewRequest: %w", err)
	}

	httpResp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("get http DefaultClient: %w", err)
	}

	defer httpResp.Body.Close()

	respBody, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		return fmt.Errorf("ioutil.ReadAll: %w", err)
	}

	lc.Logger.Debugw("lazada response", "resp", string(respBody))

	resp := &APIResponse{}
	err = json.Unmarshal(respBody, resp)
	if err != nil {
		return fmt.Errorf("response Unmarshal: %w", err)
	}

	if resp.Code != "0" {
		lc.Logger.Errorw("lazada returned failure", "body", respBody)
		return ErrProductUpdateFailed
	}

	return nil
}

func (lc *LazadaClient) QueryProduct(itemID null.Int, sellerSKU null.String) (*ProductResponse, error) {
	var req *http.Request
	var err error

	if itemID.Valid {
		lc.AddAPIParam("item_id", strconv.Itoa(int(itemID.Int64)))
	}

	if sellerSKU.Valid {
		lc.AddAPIParam("seller_sku", sellerSKU.String)
	}

	// add query params
	values := url.Values{}
	for key, val := range lc.SysParams {
		values.Add(key, val)
	}

	for key, val := range lc.APIParams {
		values.Add(key, val)
	}

	apiPath := "/product/item/get"
	apiServerURL := "https://api.lazada.vn/rest"

	values.Add("sign", lc.sign(apiPath))
	fullURL := fmt.Sprintf("%s%s?%s", apiServerURL, apiPath, values.Encode())

	req, err = http.NewRequest("GET", fullURL, nil)

	if err != nil {
		return nil, err
	}

	httpResp, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, err
	}

	defer httpResp.Body.Close()

	respBody, err := ioutil.ReadAll(httpResp.Body)

	if err != nil {
		return nil, err
	}

	resp := &ProductResponse{}
	err = json.Unmarshal(respBody, resp)

	if err != nil {
		return nil, err
	}

	resp.RawData = respBody

	if resp.Code != "0" {
		return nil, ErrProductsFailed
	}

	//@TODO: parse time

	return resp, err
}

func (lc *LazadaClient) QueryProducts(params QueryProductsParams) (*ProductsResponse, error) {
	var req *http.Request
	var err error

	lc.AddAPIParam("filter", "all")
	lc.AddAPIParam("limit", strconv.Itoa(params.Limit))
	lc.AddAPIParam("offset", strconv.Itoa(params.Offset))

	// add query params
	values := url.Values{}
	for key, val := range lc.SysParams {
		values.Add(key, val)
	}

	for key, val := range lc.APIParams {
		values.Add(key, val)
	}

	apiPath := "/products/get"
	apiServerURL := "https://api.lazada.vn/rest"

	values.Add("sign", lc.sign(apiPath))
	fullURL := fmt.Sprintf("%s%s?%s", apiServerURL, apiPath, values.Encode())

	//lc.Logger.Debugw("Making Lazada query", "url", fullURL)

	req, err = http.NewRequest("GET", fullURL, nil)

	if err != nil {
		return nil, fmt.Errorf("NewRequest: %w", err)
	}

	httpResp, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, fmt.Errorf("http Do: %w", err)
	}

	defer httpResp.Body.Close()

	respBody, err := ioutil.ReadAll(httpResp.Body)

	if err != nil {
		return nil, fmt.Errorf("ioutil ReadAll: %w", err)
	}

	resp := &ProductsResponse{}
	err = json.Unmarshal(respBody, resp)
	if err != nil {
		var unmarshalTypeError *json.UnmarshalTypeError

		if errors.As(err, &unmarshalTypeError) {
			return nil, fmt.Errorf("ProductsResponse Unmarshal: %w", unmarshalTypeError)
		}

		return nil, fmt.Errorf("ProductsResponse Unmarshal: %w", err)
	}

	resp.RawData = respBody

	if resp.Code != "0" {
		return nil, fmt.Errorf("lazada API response: %w", ErrProductsFailed)
	}

	for i := range resp.Data.Products {
		// parse the product update/created times
		p := &resp.Data.Products[i]
		if err = p.ParseTime(); err != nil {
			return nil, fmt.Errorf("product ParseTime: %w", err)
		}

		// parse sku special times
		for j := range resp.Data.Products[i].Skus {
			s := &resp.Data.Products[i].Skus[j]
			if err = s.ParseTime(); err != nil {
				return nil, fmt.Errorf("product SKU ParseTime: %w", err)
			}
		}
	}

	return resp, err
}

func (lc *LazadaClient) AddAPIParam(key string, val string) *LazadaClient {
	lc.APIParams[key] = val
	return lc
}

func (lc *LazadaClient) sign(url string) string {
	keys := []string{}
	union := map[string]string{}
	for key, val := range lc.SysParams {
		union[key] = val
		keys = append(keys, key)
	}
	for key, val := range lc.APIParams {
		union[key] = val
		keys = append(keys, key)
	}

	// sort sys params and api params by key
	sort.Strings(keys)

	var message bytes.Buffer
	message.WriteString(url)
	for _, key := range keys {
		message.WriteString(fmt.Sprintf("%s%s", key, union[key]))
	}

	hash := hmac.New(sha256.New, []byte(lc.APISecret))
	hash.Write(message.Bytes())
	return strings.ToUpper(hex.EncodeToString(hash.Sum(nil)))
}

func (lc *LazadaClient) Refresh(refresh string) (*AuthResponse, error) {
	var req *http.Request
	var err error

	lc.AddAPIParam("refresh_token", refresh)

	// add query params
	values := url.Values{}
	for key, val := range lc.SysParams {
		values.Add(key, val)
	}

	for key, val := range lc.APIParams {
		values.Add(key, val)
	}

	apiPath := "/auth/token/refresh"
	apiServerURL := "https://auth.lazada.com/rest"

	values.Add("sign", lc.sign(apiPath))
	fullURL := fmt.Sprintf("%s%s?%s", apiServerURL, apiPath, values.Encode())

	req, err = http.NewRequest("GET", fullURL, nil)

	if err != nil {
		return nil, err
	}

	httpResp, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, err
	}

	defer httpResp.Body.Close()

	t, err := http.ParseTime(httpResp.Header.Get("Date"))

	if err != nil {
		return nil, err
	}

	respBody, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		return nil, err
	}

	resp := &AuthResponse{}
	err = json.Unmarshal(respBody, resp)

	if err != nil {
		return nil, err
	}

	if resp.Code != "0" {
		return nil, ErrPlatformFailed
	}

	resp.DateRefreshExpired = t.Add(time.Second * time.Duration(resp.RefreshExpiresIn))
	resp.DateAccessExpired = t.Add(time.Second * time.Duration(resp.ExpiresIn))

	return resp, err
}

func (lc *LazadaClient) Auth(code string) (*AuthResponse, error) {
	var req *http.Request
	var err error

	lc.AddAPIParam("code", code)

	// add query params
	values := url.Values{}
	for key, val := range lc.SysParams {
		values.Add(key, val)
	}

	for key, val := range lc.APIParams {
		values.Add(key, val)
	}

	apiPath := "/auth/token/create"
	apiServerURL := "https://auth.lazada.com/rest"

	values.Add("sign", lc.sign(apiPath))
	fullURL := fmt.Sprintf("%s%s?%s", apiServerURL, apiPath, values.Encode())

	req, err = http.NewRequest("GET", fullURL, nil)

	if err != nil {
		return nil, err
	}

	httpResp, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, err
	}

	defer httpResp.Body.Close()

	t, err := http.ParseTime(httpResp.Header.Get("Date"))

	if err != nil {
		return nil, err
	}

	respBody, err := ioutil.ReadAll(httpResp.Body)

	if err != nil {
		return nil, err
	}

	resp := &AuthResponse{}
	err = json.Unmarshal(respBody, resp)

	if err != nil {
		return nil, err
	}

	if resp.Code != "0" {
		return nil, ErrPlatformFailed
	}

	resp.DateRefreshExpired = t.Add(time.Second * time.Duration(resp.RefreshExpiresIn))
	resp.DateAccessExpired = t.Add(time.Second * time.Duration(resp.ExpiresIn))

	return resp, err
}

func NewClient(id string, secret string, access string, logger *laxo.Logger) *LazadaClient {
	client := &LazadaClient{
		APIKey:    id,
		APISecret: secret,
		Logger:    logger,
		SysParams: map[string]string{
			"app_key":     id,
			"sign_method": "sha256",
			"timestamp":   fmt.Sprintf("%d000", time.Now().Unix()),
		},
		APIParams:  map[string]string{},
		FileParams: map[string][]byte{},
	}

	if access != "" {
		client.SysParams["access_token"] = access
	}

	return client
}
