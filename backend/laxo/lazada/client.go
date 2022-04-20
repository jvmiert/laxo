package lazada

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-hclog"
)

var ErrPlatformFailed = errors.New("platform returned failed message")
var ErrProductsFailed = errors.New("get products returned failed message")

type CountryUserInfo struct {
	Country     string `json:"country"`
	UserID      string `json:"user_id"`
	SellerID    string `json:"seller_id"`
	ShortCode   string `json:"short_code"`
}

type AuthResponse struct {
	Code               string            `json:"code"`
	Type               string            `json:"type"`
	Message            string            `json:"message"`
	RequestID          string            `json:"request_id"`


	AccessToken        string            `json:"access_token"`
	Country            string            `json:"country"`
	RefreshToken       string            `json:"refresh_token"`
	AccountPlatform    string            `json:"account_platform"`
	RefreshExpiresIn   int               `json:"refresh_expires_in"`
	CountryUserInfo    []CountryUserInfo `json:"country_user_info"`
	ExpiresIn          int               `json:"expires_in"`
	Account            string            `json:"account"`

  DateAccessExpired  time.Time
  DateRefreshExpired time.Time
}

type ProductsResponseAttributes struct {
  Name                 string `json:"name"`
  ShortDescription     string `json:"short_description"`
  Description          string `json:"description"`
  Brand                string `json:"brand"`
  Model                string `json:"model"`
  HeadphoneFeatures    string `json:"headphone_features"`
  Bluetooth            string `json:"bluetooth"`
  WarrantyType         string `json:"warranty_type"`
  Warranty             string `json:"warranty"`
  Hazmat               string `json:"Hazmat"`
  ExpireDate           string `json:"Expire_date"`
  BrandClassification  string `json:"brand_classification"`
  IngredientPreference string `json:"ingredient_preference"`
  LotNumber            string `json:"Lot_number"`
  UnitsHB              string `json:"units_hb"`
  FmltSkinCare         string `json:"fmlt_skin_care"`
  Quantitative         string `json:"Quantitative"`
  SkinCareByAge        string `json:"skin_care_by_age"`
  SkinBenefit          string `json:"skin_benefit"`
  SkinType             string `json:"skin_type"`
  UserManual           string `json:"User_Manual"`
  CountryOriginHB      string `json:"country_origin_hb"`
  ColorFamily          string `json:"color_family"`
  FragranceFamily      string `json:"fragrance_family"`
  Source               string `json:"source"`
}

type ProductsResponseSuspendedSkus struct {
  RejectReason    string   `json:"rejectReason"`
  SellerSku       string   `json:"SellerSku"`
  SkuID           int      `json:"SkuId"`
}

type ProductsResponseSkus struct {
  Status              string   `json:"Status"`
  Quantity            int      `json:"quantity"`
  Images              []string `json:"Images"`
  MarketImages        []string `json:"marketImages"`
  SellerSku           string   `json:"SellerSku"`
  ShopSku             string   `json:"ShopSku"`
  PackageContent      string   `json:"package_content"`
  URL                 string   `json:"Url"`
  PackageWidth        string   `json:"package_width"`
  SpecialToTime       string   `json:"special_to_time"`
  ColorFamily         string   `json:"color_family"`
  SpecialFromTime     string   `json:"special_from_time"`
  PackageHeight       string   `json:"package_height"`
  SpecialPrice        float64  `json:"special_price"`
  Price               float64  `json:"price"`
  PackageLength       string   `json:"package_length"`
  SpecialFromDate     string   `json:"special_from_date"`
  PackageWeight       string   `json:"package_weight"`
  Available           int      `json:"Available"`
  SkuID               int      `json:"SkuId"`
  SpecialToDate       string   `json:"special_to_date"`
}

type ProductsResponseProducts struct {
  Skus []ProductsResponseSkus                   `json:"skus"`
  ItemID          int                           `json:"item_id"`
  PrimaryCategory int                           `json:"primary_category"`
  Attributes      ProductsResponseAttributes    `json:"attributes"`
  CreatedTime     string                        `json:"created_time"`
  UpdatedTime     string                        `json:"updated_time"`
  Images          []string                      `json:"images"`
  MarketImages    []string                      `json:"marketImages"`
  Status          string                        `json:"status"`
  SubStatus       string                        `json:"subStatus"`
  SuspendedSkus   ProductsResponseSuspendedSkus `json:"suspendedSkus"`
}

type ProductsResponseData struct {
		TotalProducts int                        `json:"total_products"`
		Products      []ProductsResponseProducts `json:"products"`
}

type ProductsResponse struct {
  Data      ProductsResponseData  `json:"data"`
	Code      string                `json:"code"`
	RequestID string                `json:"request_id"`
}

type QueryProductsParams struct {
  Filter          string
  UpdateBefore    time.Time
  CreatedBefore   time.Time
  Offset          int
  CreatedAfter    time.Time
  UpdatedAfter    time.Time
  Limit           int
}

type LazadaClient struct {
	APIKey    string
	APISecret string
	Region    string

  Logger    hclog.Logger

	Method     string
	SysParams  map[string]string
	APIParams  map[string]string
	FileParams map[string][]byte
}

func (lc *LazadaClient) QueryProducts(params QueryProductsParams) (*ProductsResponse, error) {
	var req *http.Request
	var err error

  lc.AddAPIParam("filter", "live")
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

  lc.Logger.Debug("Making product query", "url", fullURL)

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

  resp := &ProductsResponse{}
  err = json.Unmarshal(respBody, resp)

  if err != nil {
    return nil, err
  }

  if resp.Code != "0" {
    return nil, ErrProductsFailed
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

func NewClient(id string, secret string, access string, logger hclog.Logger) *LazadaClient {
  client := &LazadaClient{
		APIKey:    id,
		APISecret: secret,
    Logger: logger,
		SysParams: map[string]string{
			"app_key":       id,
			"sign_method":   "sha256",
			"timestamp":     fmt.Sprintf("%d000", time.Now().Unix()),
		},
		APIParams:  map[string]string{},
		FileParams: map[string][]byte{},
	}

  if access != "" {
    client.SysParams["access_token"] = access
  }

  return client
}
