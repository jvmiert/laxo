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
	"strings"
	"time"

	"github.com/hashicorp/go-hclog"
)

var ErrPlatformFailed = errors.New("platform returned failed message")

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

func NewClient(id string, secret string, logger hclog.Logger) *LazadaClient {
	return &LazadaClient{
		APIKey:    id,
		APISecret: secret,
    Logger: logger,
		SysParams: map[string]string{
			"app_key":     id,
			"sign_method": "sha256",
			"timestamp":   fmt.Sprintf("%d000", time.Now().Unix()),
		},
		APIParams:  map[string]string{},
		FileParams: map[string][]byte{},
	}
}
