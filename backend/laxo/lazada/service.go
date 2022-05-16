package lazada

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"math"
	"os"
	"strconv"

	"github.com/hashicorp/go-hclog"
	"github.com/mediocregopher/radix/v4"
	"github.com/microcosm-cc/bluemonday"
	"laxo.vn/laxo/laxo/sqlc"
)

var ErrNoValidToken = errors.New("no valid token was found")
var ErrClientCredentialsNotSet = errors.New("the Lazada client ID or client secret was not set")
var ErrProductIndexNotFound = errors.New("the index returned empty")

const redisKeyPrefix = "product_lazada_"

type Store interface {
  SaveOrUpdateLazadaProduct(*ProductsResponseProducts, string) (*sqlc.ProductsLazada, *sqlc.ProductsAttributeLazada, *sqlc.ProductsSkuLazada, error)
  GetValidTokenByShopID(string) (string, error)
}

type Service struct {
  store Store
  logger hclog.Logger
  redisClient radix.Client
  clientID string
  clientSecret string
}

func NewService(store Store, logger hclog.Logger, redisClient radix.Client, clientID, clientSecret string) Service {
  return Service {
    store: store,
    logger: logger,
    redisClient: redisClient,
    clientID: clientID,
    clientSecret: clientSecret,
  }
}

func (s *Service) NewLazadaClient(token string) (*LazadaClient, error) {
  if s.clientID == "" || s.clientSecret == "" {
    return nil, ErrClientCredentialsNotSet
  }

  client := NewClient(s.clientID, s.clientSecret, token, s.logger)

  return client, nil
}

func (s *Service) GetSantizedDescription(d string) (string) {
  p := bluemonday.StrictPolicy()
  santized := p.Sanitize(d)

  return santized
}

func (s *Service) GetValidTokenByShopID(shopID string) (string, error) {
  token, err := s.store.GetValidTokenByShopID(shopID)
  if err != nil {
    return "", err
  }

  return token, nil
}

func (s *Service) RetrieveProductFromRedis(keyID string, index int) (*ProductsResponseProducts, error) {
  i := strconv.Itoa(index)
  ctx := context.Background()

  var rcv string
  mb := radix.Maybe{Rcv: &rcv}

  if err := s.redisClient.Do(ctx, radix.Cmd(&mb, "HGET", keyID, i)); err != nil {
    return nil, err
  }

  if mb.Null || mb.Empty {
    return nil, ErrProductIndexNotFound
  }

  var response ProductsResponseProducts

  err := json.Unmarshal([]byte(*mb.Rcv.(*string)), &response)
  if err != nil {
    return nil, err
  }

  return &response, nil
}

func (s *Service) SaveProductToRedis(keyID string, index int, p ProductsResponseProducts) error {
  bytes, err := json.Marshal(p)
  if err != nil {
    return err
  }

  i := strconv.Itoa(index)
  m := map[string]string{i: string(bytes)}
  ctx := context.Background()
  if err := s.redisClient.Do(ctx, radix.FlatCmd(nil, "HMSET", keyID, m)); err != nil {
    return err
  }
  return nil
}

func (s *Service) FetchProductsFromLazadaToRedis(shopID string) (string, error) {
  token := ""
  //token, err := s.GetValidTokenByShopID(shopID)
  //if err != nil {
  //  if err == pgx.ErrNoRows {
  //    return "", ErrNoValidToken
  //  }
  //  return "", err
  //}

  client, err := s.NewLazadaClient(token)
  if err != nil {
    return "", err
  }

  //response, err := client.QueryProducts(QueryProductsParams{
  //  Limit: 50,
  //  Offset: 0,
  //})
  //if err != nil {
  //  return "", err
  //}

  //-------------------------------------
  // for reading json file
  // so we don't have to bother Lazada
  //-------------------------------------
  jsonFile, err := os.Open(".\\laxo\\lazada\\test.json")
  if err != nil {
    return "", err
  }

  defer jsonFile.Close()

  byteValue, err := ioutil.ReadAll(jsonFile)
  if err != nil {
    return "", err
  }

  var response ProductsResponse

  json.Unmarshal(byteValue, &response)
  //-------------------------------------

  keyID := redisKeyPrefix + shopID

  for i, product := range response.Data.Products {
    if err := s.SaveProductToRedis(keyID, i, product); err != nil {
      return "", err
    }
  }

  var fetchesRequired int
  if response.Data.TotalProducts > 50 {
    remainingProducts := (float64)(response.Data.TotalProducts - 50)
    fetchesRequired = (int)(math.Ceil(remainingProducts / 50))
  } else {
    return keyID, nil
  }

  for i := 0; i < fetchesRequired; i++ {
    response, err := client.QueryProducts(QueryProductsParams{
      Limit: 50,
      Offset: 50 * i,
    })
    if err != nil {
      return "", err
    }

    for j, product := range response.Data.Products {
      index := (i + 1) * 50 + j
      if err := s.SaveProductToRedis(keyID, index, product); err != nil {
        return "", err
      }
    }
  }

  return "", nil
}

func (s *Service) SaveOrUpdateLazadaProduct(p *ProductsResponseProducts, shopID string) (*sqlc.ProductsLazada, *sqlc.ProductsAttributeLazada, *sqlc.ProductsSkuLazada, error) {
  return s.store.SaveOrUpdateLazadaProduct(p, shopID)
}
