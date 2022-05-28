package shop

import (
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
	"laxo.vn/laxo/laxo/lazada"
)

func (s *Service) GenerateRedirect(r *OAuthRedirectRequest) error {
  var aPlatforms = availablePlatforms{
    "lazada": struct{}{},
    "tiki": struct{}{},
    "shopee": struct{}{},
  }

  var cPlatforms []string

  connectedPlatforms, err := s.store.RetrieveShopsPlatformsByShopID(r.ShopID)

  if err != nil {
    s.server.Logger.Error("OauthRedirectRequest redirect database error",
      "error", err,
    )
    return err
  }

  for _, v := range connectedPlatforms {
    if aPlatforms.Has(v.PlatformName) {
      cPlatforms = append(cPlatforms, v.PlatformName);
    }
  }

  secretKey := os.Getenv("OAUTH_HMAC_SECRET")

  if secretKey == "" {
    s.server.Logger.Error("OAuth HMAC secret key is not set!")
    return errors.New("secret key not set")
  }

  for platform := range aPlatforms {
    var url strings.Builder

    hash := hmac.New(sha256.New, []byte(secretKey))
    io.WriteString(hash, r.ShopID+platform)

    if platform == "lazada" {
      clientID := os.Getenv("LAZADA_ID")

      if clientID == "" {
        s.server.Logger.Error("Lazada client id is not set!")
        return errors.New("lazada client id not set")
      }

      url.WriteString("https://auth.lazada.com/oauth/authorize?response_type=code&country=vn&force_auth=true&redirect_uri=")
      url.WriteString(s.server.Config.CallbackBasePath)
      url.WriteString("lazada&client_id=")
      url.WriteString(clientID)
      url.WriteString("&state=")
      url.WriteString(hex.EncodeToString(hash.Sum(nil)))
    }

    if platform == "tiki" {
      clientID := os.Getenv("TIKI_ID")

      if clientID == "" {
        s.server.Logger.Error("Tiki client id is not set!")
        return errors.New("tiki client id not set")
      }

      url.WriteString("https://api.tiki.vn/sc/oauth2/auth?response_type=code&redirect_uri=")
      url.WriteString(s.server.Config.CallbackBasePath)
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

    r.ReturnRedirects = append(r.ReturnRedirects, &ReturnRedirect{Platform: platform,URL: fmt.Sprint(url.String())})
    r.Connected = cPlatforms
  }

  return nil
}

func (s *Service) ValidateOAuthRedirectRequest(r *OAuthRedirectRequest, uID string, printer *message.Printer) error {
  dbShop, err := s.store.GetShopByID(r.ShopID)

  if err == pgx.ErrNoRows {
    return validation.Errors{
      "shopID": validation.NewError(
        "not_exists",
        printer.Sprintf("shop does not exist")),
    }
  } else if err != nil {
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

func (s *Service) ValidateOAuthVerifyRequest(o OAuthVerifyRequest, uID string, printer *message.Printer) error {
  if o.Code == "" {
    return validation.Errors{
      "code": validation.NewError(
        "code_missing",
        printer.Sprintf("code is required")),
    }
  }

  shop, err := s.GetActiveShopByUserID(uID)

  if err == ErrUserNoShops {
    return validation.Errors{
      "state": validation.NewError(
        "user_no_shops",
        printer.Sprintf("user does not have any shops")),
    }
  } else if err != nil {
    s.server.Logger.Errorw("OAuthVerifyRequest validation error",
      "error", err,
    )
    return validation.Errors{
      "state": validation.NewError(
        "general_failure",
        printer.Sprintf("something went wrong")),
    }
  }

  secretKey := os.Getenv("OAUTH_HMAC_SECRET")

  if secretKey == "" {
    s.server.Logger.Error("OAuth HMAC secret key is not set!")
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
        s.server.Logger.Error("Lazada client id or secret is not set!")
        return validation.Errors{
          "code": validation.NewError(
            "general_failure",
            printer.Sprintf("something went wrong")),
        }
      }

      client := lazada.NewClient(clientID, clientSecret, "", s.server.Logger)

      authResp, err := client.Auth(o.Code)

      if err != nil {
        s.server.Logger.Errorw("Lazada token request error",
          "error", err,
        )
        return validation.Errors{
          "code": validation.NewError(
            "general_failure",
            printer.Sprintf("something went wrong")),
        }
      }

      lazInfo, err := s.store.GetLazadaPlatformByShopID(
        shop.Model.ID,
      )

      if err == pgx.ErrNoRows {
        _, dbErr := s.store.SaveNewLazadaPlatform(shop.Model.ID, authResp)

        if dbErr != nil {
          return dbErr
        }
      } else if err != nil {
        s.server.Logger.Errorw("Platform retrieval error",
          "error", err,
        )
        return validation.Errors{
          "code": validation.NewError(
            "general_failure",
            printer.Sprintf("something went wrong")),
        }
      }

      err = s.store.UpdateLazadaPlatform(lazInfo.ID, authResp)
      if err != nil {
        s.server.Logger.Errorw("UpdateLazadaPlatform update error",
          "error", err,
        )
        return validation.Errors{
          "code": validation.NewError(
            "general_failure",
            printer.Sprintf("something went wrong")),
        }
      }

      _, err = s.store.RetrieveSpecificPlatformByShopID(shop.Model.ID, "lazada")
      if err == pgx.ErrNoRows {
        _, err = s.store.CreateShopsPlatforms(shop.Model.ID, "lazada")
        if err != nil {
          s.server.Logger.Errorw("Platform creation error",
            "error", err,
            )
          return validation.Errors{
            "code": validation.NewError(
              "general_failure",
              printer.Sprintf("something went wrong")),
          }
        }
      } else if err != nil {
        s.server.Logger.Errorw("RetrieveSpecificPlatformByShopID retrieval error",
          "error", err,
        )
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
