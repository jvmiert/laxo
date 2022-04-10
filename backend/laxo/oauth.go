package laxo

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/jackc/pgx/v4"
	"golang.org/x/text/message"
)

type OAuthRedirectRequest struct {
	ShopID      string   `json:"shopID"`
	Platform    string   `json:"platform"`
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

  if s.Platform != "lazada" && s.Platform != "shopee" && s.Platform != "tiki" {
    return validation.Errors{
      "platform": validation.NewError(
        "unknown_platform",
        printer.Sprintf("you supplied an unknown platform")),
    }
  }

  return nil
}

func (s *OAuthRedirectRequest) GenerateRedirect() (string, error) {
  var url strings.Builder

  secretKey := os.Getenv("OAUTH_HMAC_SECRET")

  if secretKey == "" {
    Logger.Error("OAuth HMAC secret key is not set!")
    return "", errors.New("secret key not set")
  }

  hash := hmac.New(sha256.New, []byte(secretKey))
  io.WriteString(hash, s.ShopID+s.Platform)

  if s.Platform == "lazada" {
    clientID := os.Getenv("LAZADA_ID")

    if clientID == "" {
      Logger.Error("Lazada client id is not set!")
      return "", errors.New("lazada client id not set")
    }

    url.WriteString("https://auth.lazada.com/oauth/authorize?response_type=code&force_auth=true&redirect_uri=")
    url.WriteString(AppConfig.CallbackBasePath)
    url.WriteString("lazada&client_id=")
    url.WriteString(clientID)
    url.WriteString("&state=")
    url.WriteString(hex.EncodeToString(hash.Sum(nil)))
  }

  if s.Platform == "tiki" {
    clientID := os.Getenv("TIKI_ID")

    if clientID == "" {
      Logger.Error("Tiki client id is not set!")
      return "", errors.New("tiki client id not set")
    }

    url.WriteString("https://api.tiki.vn/sc/oauth2/auth?response_type=code&redirect_uri=")
    url.WriteString(AppConfig.CallbackBasePath)
    url.WriteString("tiki&client_id=")
    url.WriteString(clientID)
    url.WriteString("&state=")
    url.WriteString(hex.EncodeToString(hash.Sum(nil)))
  }

  if s.Platform == "shopee" {
    clientID := os.Getenv("SHOPEE_ID")

    if clientID == "" {
      Logger.Error("Shopee client id is not set!")
      return "", errors.New("shopee client id not set")
    }

    // @TODO: https://open.shopee.com/documents/v2/OpenAPI%202.0%20Overview?module=87&type=2

    url.WriteString("NOT IMPLEMENTED YET")
  }

  return fmt.Sprint(url.String()), nil
}

// @TODO: When retrieving the callback, match the hash with: http://www.inanzzz.com/index.php/post/g4nt/signing-messages-and-verifying-integrity-with-a-secret-using-hmac-in-golang
