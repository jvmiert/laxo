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

	"database/sql"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/jackc/pgx/v4"
	"golang.org/x/text/message"
	"laxo.vn/laxo/laxo/lazada"
	"laxo.vn/laxo/laxo/sqlc"
)

type OAuthVerifyRequest struct {
  Platform     string     `json:"platform"`
  Code         string     `json:"code"`
  State        string     `json:"state"`
}

func (o *OAuthVerifyRequest) Verify(uID string, printer *message.Printer) error {
  if o.Code == "" {
    return validation.Errors{
      "code": validation.NewError(
        "code_missing",
        printer.Sprintf("code is required")),
    }
  }

  shop, err := GetActiveShopByUserID(uID)

  if err == ErrUserNoShops {
    return validation.Errors{
      "state": validation.NewError(
        "user_no_shops",
        printer.Sprintf("user does not have any shops")),
    }
  } else if err != nil {
    Logger.Error("OAuthVerifyRequest validation error", "error", err)
    return validation.Errors{
      "state": validation.NewError(
        "general_failure",
        printer.Sprintf("something went wrong")),
    }
  }

  secretKey := os.Getenv("OAUTH_HMAC_SECRET")

  if secretKey == "" {
    Logger.Error("OAuth HMAC secret key is not set!")
    return errors.New("secret key not set")
  }

  if o.State == "" && o.Platform != "shopee" {
    return validation.Errors{
      "state": validation.NewError(
        "invalid_state",
        printer.Sprintf("supplied state is invalid")),
    }
  }

  if o.State != "" {
    recHash, err := hex.DecodeString(o.State)

    if err != nil {
      return validation.Errors{
        "state": validation.NewError(
          "invalid_state",
          printer.Sprintf("supplied state is invalid")),
      }
    }

    hash := hmac.New(sha256.New, []byte(secretKey))
    io.WriteString(hash, shop.Model.ID+o.Platform)

    if ok := hmac.Equal(recHash, hash.Sum(nil)); !ok {
      return validation.Errors{
        "state": validation.NewError(
          "invalid_state",
          printer.Sprintf("supplied state is invalid")),
      }
    }

    // validate the actual code
    if o.Platform == "lazada" {
      clientID := os.Getenv("LAZADA_ID")
      clientSecret := os.Getenv("LAZADA_SECRET")

      if clientID == "" || clientSecret == "" {
        Logger.Error("Lazada client id or secret is not set!")
        return validation.Errors{
          "code": validation.NewError(
            "general_failure",
            printer.Sprintf("something went wrong")),
        }
      }

      client := lazada.NewClient(clientID, clientSecret, "", Logger)

      authResp, err := client.Auth(o.Code)

      if err != nil {
        Logger.Error("Lazada token request error", "error", err)
        return validation.Errors{
          "code": validation.NewError(
            "general_failure",
            printer.Sprintf("something went wrong")),
        }
      }

      lazInfo, err := Queries.GetLazadaPlatformByShopID(
        context.Background(),
        shop.Model.ID,
      )

      if err == pgx.ErrNoRows {
        _, dbErr := Queries.CreateLazadaPlatform(
          context.Background(),
          sqlc.CreateLazadaPlatformParams{
            ShopID: shop.Model.ID,
            AccessToken: authResp.AccessToken,
            Country: authResp.Country,
            RefreshToken: authResp.RefreshToken,
            AccountPlatform: authResp.AccountPlatform,
            Account: authResp.Account,
            UserIDVn: authResp.CountryUserInfo[0].UserID,
            SellerIDVn: authResp.CountryUserInfo[0].SellerID,
            ShortCodeVn: authResp.CountryUserInfo[0].ShortCode,
            RefreshExpiresIn: sql.NullTime{Time: authResp.DateRefreshExpired, Valid: true},
            AccessExpiresIn: sql.NullTime{Time: authResp.DateAccessExpired, Valid: true},
          },
        )

        if dbErr != nil {
          return dbErr
        }
      } else if err != nil {
        Logger.Error("Platform retrieval error", err)
        return validation.Errors{
          "code": validation.NewError(
            "general_failure",
            printer.Sprintf("something went wrong")),
        }
      }

      err = Queries.UpdateLazadaPlatform(
        context.Background(),
        sqlc.UpdateLazadaPlatformParams{
          AccessToken: authResp.AccessToken,
          RefreshToken: authResp.RefreshToken,
          RefreshExpiresIn: sql.NullTime{Time: authResp.DateRefreshExpired, Valid: true},
          AccessExpiresIn: sql.NullTime{Time: authResp.DateAccessExpired, Valid: true},
          ID: lazInfo.ID,
        },
      )
      if err != nil {
        Logger.Error("Platform update error", err)
        return validation.Errors{
          "code": validation.NewError(
            "general_failure",
            printer.Sprintf("something went wrong")),
        }
      }
    }
  }


  return nil
}

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

