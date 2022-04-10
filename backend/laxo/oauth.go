package laxo

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/jackc/pgx/v4"
	"golang.org/x/text/message"
)

type ReturnRedirect struct {
  Platform     string     `json:"platform"`
  URL          string     `json:"url"`
}

type availablePlatforms map[string]struct{}

func (a availablePlatforms) Has(v string) bool {
  _, ok := a[v]
  return ok
}

type OAuthRedirectRequest struct {
	ShopID             string              `json:"shopID"`
	ReturnRedirects    []*ReturnRedirect   `json:"platforms"`
}


func (s *OAuthRedirectRequest) Validate(uID string, printer *message.Printer) error {
  dbShop, err := Queries.GetShopByID(context.Background(), s.ShopID)

  if err == pgx.ErrNoRows {
    return validation.Errors{
      "shopID": validation.NewError(
        "not_exists",
        printer.Sprintf("shop does not exist")),
    }
  } else if err != nil {
    Logger.Error("OauthRedirectRequest validation error", "error", err)
    return err
  }

  if dbShop.UserID != uID {
    return validation.Errors{
      "shopID": validation.NewError(
        "not_owned",
        printer.Sprintf("you don't own this shop")),
    }
  }

  return nil
}

func (s *OAuthRedirectRequest) JSON() ([]byte, error) {
  sort.Slice(s.ReturnRedirects, func(i, j int) bool { return s.ReturnRedirects[i].Platform < s.ReturnRedirects[j].Platform })
  bytes, err := json.Marshal(s)

  if err != nil {
    Logger.Error("Shop marshal error", "error", err)
    return bytes, err
  }

  return bytes, nil
}

func (s *OAuthRedirectRequest) GenerateRedirect() error {
  var aPlatforms = availablePlatforms{
    "lazada": struct{}{},
    "tiki": struct{}{},
    "shopee": struct{}{},
  }

  connectedPlatforms, err := Queries.GetPlatformsByShopID(context.Background(), s.ShopID)

  if err != nil {
    Logger.Error("OauthRedirectRequest redirect database error", "error", err)
    return err
  }

  for _, v := range connectedPlatforms {
    if aPlatforms.Has(v.PlatformName) {
      delete(aPlatforms, v.PlatformName);
    }
  }

  secretKey := os.Getenv("OAUTH_HMAC_SECRET")

  if secretKey == "" {
    Logger.Error("OAuth HMAC secret key is not set!")
    return errors.New("secret key not set")
  }

  for platform := range aPlatforms {
    var url strings.Builder

    hash := hmac.New(sha256.New, []byte(secretKey))
    io.WriteString(hash, s.ShopID+platform)

    if platform == "lazada" {
      clientID := os.Getenv("LAZADA_ID")

      if clientID == "" {
        Logger.Error("Lazada client id is not set!")
        return errors.New("lazada client id not set")
      }

      url.WriteString("https://auth.lazada.com/oauth/authorize?response_type=code&force_auth=true&redirect_uri=")
      url.WriteString(AppConfig.CallbackBasePath)
      url.WriteString("lazada&client_id=")
      url.WriteString(clientID)
      url.WriteString("&state=")
      url.WriteString(hex.EncodeToString(hash.Sum(nil)))
    }

    if platform == "tiki" {
      clientID := os.Getenv("TIKI_ID")

      if clientID == "" {
        Logger.Error("Tiki client id is not set!")
        return errors.New("tiki client id not set")
      }

      url.WriteString("https://api.tiki.vn/sc/oauth2/auth?response_type=code&redirect_uri=")
      url.WriteString(AppConfig.CallbackBasePath)
      url.WriteString("tiki&client_id=")
      url.WriteString(clientID)
      url.WriteString("&state=")
      url.WriteString(hex.EncodeToString(hash.Sum(nil)))
    }

    if platform == "shopee" {
      //clientID := os.Getenv("SHOPEE_ID")

      //if clientID == "" {
      //  Logger.Error("Shopee client id is not set!")
      //  return errors.New("shopee client id not set")
      //}

      // @TODO: https://open.shopee.com/documents/v2/OpenAPI%202.0%20Overview?module=87&type=2

      url.WriteString("NOT IMPLEMENTED YET")
    }

    s.ReturnRedirects = append(s.ReturnRedirects, &ReturnRedirect{Platform: platform,URL: fmt.Sprint(url.String())})
  }

  return nil
}

// @TODO: When retrieving the callback, match the hash with: http://www.inanzzz.com/index.php/post/g4nt/signing-messages-and-verifying-integrity-with-a-secret-using-hmac-in-golang

// @TODO: Handle this callback url for lazada: http://localhost:3000/setup-shop/callback/lazada?code=0_108382_cyxXWvT3otZCM41iODqfxx571782&state=9a430e6343286e862fa8921f9afb9bff22b6f94bf2f4cdd6aed019cd13dac51c
