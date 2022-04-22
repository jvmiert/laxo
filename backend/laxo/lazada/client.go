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

	"github.com/araddon/dateparse"
	"github.com/hashicorp/go-hclog"
	"gopkg.in/guregu/null.v4"
)

var ErrPlatformFailed = errors.New("platform returned failed message")
var ErrProductsFailed = errors.New("get products returned failed message")
var ErrProductsParseFailed = errors.New("couldn't parse products")

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
  Name                 null.String `json:"name"`
  ShortDescription     null.String `json:"short_description"`
  Description          null.String `json:"description"`
  Brand                null.String `json:"brand"`
  Model                null.String `json:"model"`
  HeadphoneFeatures    null.String `json:"headphone_features"`
  Bluetooth            null.String `json:"bluetooth"`
  WarrantyType         null.String `json:"warranty_type"`
  Warranty             null.String `json:"warranty"`
  Hazmat               null.String `json:"Hazmat"`
  ExpireDate           null.String `json:"Expire_date"`
  BrandClassification  null.String `json:"brand_classification"`
  IngredientPreference null.String `json:"ingredient_preference"`
  LotNumber            null.String `json:"Lot_number"`
  UnitsHB              null.String `json:"units_hb"`
  FmltSkinCare         null.String `json:"fmlt_skin_care"`
  Quantitative         null.String `json:"Quantitative"`
  SkinCareByAge        null.String `json:"skin_care_by_age"`
  SkinBenefit          null.String `json:"skin_benefit"`
  SkinType             null.String `json:"skin_type"`
  UserManual           null.String `json:"User_Manual"`
  CountryOriginHB      null.String `json:"country_origin_hb"`
  ColorFamily          null.String `json:"color_family"`
  FragranceFamily      null.String `json:"fragrance_family"`
  Source               null.String `json:"source"`
}

type ProductsResponseSuspendedSkus struct {
  RejectReason    null.String   `json:"rejectReason"`
  SellerSku       null.String   `json:"SellerSku"`
  SkuID           null.Int      `json:"SkuId"`
}

type ProductsResponseSkus struct {
  Status               null.String   `json:"Status"`
  Quantity             null.Int      `json:"quantity"`
  Images               []null.String `json:"Images"`
  MarketImages         []null.String `json:"marketImages"`
  SellerSku            null.String   `json:"SellerSku"`
  ShopSku              null.String   `json:"ShopSku"`
  PackageContent       null.String   `json:"package_content"`
  URL                  null.String   `json:"Url"`
  PackageWidth         null.String   `json:"package_width"`
  ColorFamily          null.String   `json:"color_family"`
  PackageHeight        null.String   `json:"package_height"`
  SpecialPrice         null.Float    `json:"special_price"`
  Price                null.Float    `json:"price"`
  PackageLength        null.String   `json:"package_length"`
  PackageWeight        null.String   `json:"package_weight"`
  Available            null.Int      `json:"Available"`
  SkuID                null.Int      `json:"SkuId"`
  SpecialToTimeRaw     null.String   `json:"special_to_time"`
  SpecialFromTimeRaw   null.String   `json:"special_from_time"`
  SpecialFromDateRaw   null.String   `json:"special_from_date"`
  SpecialToDateRaw     null.String   `json:"special_to_date"`
  SpecialToTime        time.Time
  SpecialFromTime      time.Time
  SpecialFromDate      time.Time
  SpecialToDate        time.Time
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

type ProductsResponseProducts struct {
  Skus []ProductsResponseSkus                        `json:"skus"`
  ItemID          null.Int                           `json:"item_id"`
  PrimaryCategory null.Int                           `json:"primary_category"`
  Attributes      ProductsResponseAttributes         `json:"attributes"`
  CreatedTimeRaw  null.String                        `json:"created_time"`
  UpdatedTimeRaw  null.String                        `json:"updated_time"`
  Images          []null.String                      `json:"images"`
  MarketImages    []null.String                      `json:"marketImages"`
  Status          null.String                        `json:"status"`
  SubStatus       null.String                        `json:"subStatus"`
  SuspendedSkus   ProductsResponseSuspendedSkus      `json:"suspendedSkus"`
  CreatedTime     time.Time
  UpdatedTime     time.Time
}

func (p *ProductsResponseProducts) ParseTime() error {
  if !p.CreatedTimeRaw.Valid || !p.UpdatedTimeRaw.Valid {
    return ErrProductsParseFailed
  }

  i, err := strconv.ParseInt(p.CreatedTimeRaw.String, 10, 64)
  if err != nil {
    return err
  }

  p.CreatedTime = time.Unix(i/1000, 0)


  i, err = strconv.ParseInt(p.UpdatedTimeRaw.String, 10, 64)
  if err != nil {
    return err
  }

  p.UpdatedTime = time.Unix(i/1000, 0)
  return nil
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

  for i := range resp.Data.Products {
    // parse the product update/created times
    p := &resp.Data.Products[i]
    if err = p.ParseTime(); err != nil {
      return nil, err
    }

    // parse sku special times
    for j := range resp.Data.Products[i].Skus {
      s := &resp.Data.Products[i].Skus[j]
      if err = s.ParseTime(); err != nil {
        return nil, err
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
