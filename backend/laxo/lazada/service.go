package lazada

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strconv"

	"github.com/jackc/pgx/v4"
	"github.com/mediocregopher/radix/v4"
	"laxo.vn/laxo/laxo"
	"laxo.vn/laxo/laxo/models"
	"laxo.vn/laxo/laxo/sqlc"
)

var ErrNoValidToken = errors.New("no valid token was found")
var ErrClientCredentialsNotSet = errors.New("the Lazada client ID or client secret was not set")
var ErrProductIndexNotFound = errors.New("the index returned empty")

const redisKeyPrefix = "product_lazada_"

type Store interface {
	SaveOrUpdateLazadaProduct(*ProductsResponseProducts, string) (*sqlc.ProductsLazada, *sqlc.ProductsAttributeLazada, *sqlc.ProductsSkuLazada, error)
	GetValidAccessTokenByShopID(string) (string, error)
	GetValidRefreshTokenByShopID(string) (string, error)
	SaveNewLazadaPlatform(string, *AuthResponse) (*sqlc.PlatformLazada, error)
	UpdateLazadaPlatform(string, *AuthResponse) error
	GetLazadaPlatformByShopID(string) (*sqlc.PlatformLazada, error)
}

type Service struct {
	store        Store
	logger       *laxo.Logger
	server       *laxo.Server
	clientID     string
	clientSecret string
}

func NewService(store Store, logger *laxo.Logger, server *laxo.Server, clientID, clientSecret string) Service {
	return Service{
		store:        store,
		logger:       logger,
		server:       server,
		clientID:     clientID,
		clientSecret: clientSecret,
	}
}

func (s *Service) ExtractImagesListFromProductResponse(p *ProductsResponseProducts) ([]string, error) {
	var imageList []string

	for _, v := range p.Images {
		if v.Valid {
			imageList = append(imageList, v.String)
		}
	}

	for _, v := range p.MarketImages {
		if v.Valid {
			imageList = append(imageList, v.String)
		}
	}

	if len(p.Skus) > 0 {
		for _, v := range p.Skus[0].Images {
			if v.Valid {
				imageList = append(imageList, v.String)
			}
		}

		for _, v := range p.Skus[0].MarketImages {
			if v.Valid {
				imageList = append(imageList, v.String)
			}
		}
	}

	return imageList, nil
}

func (s *Service) GetLazadaPlatformJSON(p *sqlc.PlatformLazada) ([]byte, error) {
	pReturn := PlatformLazadaReturn{
		ID:               p.ID,
		ShopID:           p.ShopID,
		Country:          p.Country,
		AccountPlatform:  p.AccountPlatform,
		Account:          p.Account,
		UserIDVn:         p.UserIDVn,
		SellerIDVn:       p.SellerIDVn,
		ShortCodeVn:      p.ShortCodeVn,
		RefreshExpiresIn: p.RefreshExpiresIn,
		AccessExpiresIn:  p.AccessExpiresIn,
		Created:          p.Created,
	}

	return json.Marshal(pReturn)
}

func (s *Service) GetLazadaPlatformByShopID(shopID string) (*sqlc.PlatformLazada, error) {
	return s.store.GetLazadaPlatformByShopID(shopID)
}

func (s *Service) NewLazadaClient(token string) (*LazadaClient, error) {
	if s.clientID == "" || s.clientSecret == "" {
		return nil, ErrClientCredentialsNotSet
	}

	client := NewClient(s.clientID, s.clientSecret, token, s.logger)

	return client, nil
}

func (s *Service) UpdateProductToLazada(p *models.ProductDetails, descriptionHTML string) error {
	token, err := s.GetValidTokenByShopID(p.Model.ShopID)
	if err != nil {
		return fmt.Errorf("GetValidTokenByShopID: %w", err)
	}

	client, err := s.NewLazadaClient(token)
	if err != nil {
		return fmt.Errorf("NewLazadaClient: %w", err)
	}

	err = client.UpdateProduct(p, descriptionHTML)
	if err != nil {
		return fmt.Errorf("client UpdateProduct: %w", err)
	}

	return nil
}

func (s *Service) RefreshAndGetTokenByShopID(shopID string) (string, error) {
	refreshToken, err := s.store.GetValidRefreshTokenByShopID(shopID)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", ErrNoValidToken
	} else if err != nil {
		return "", fmt.Errorf("GetValidRefreshTokenByShopID: %w", err)
	}

	client, err := s.NewLazadaClient("")
	if err != nil {
		return "", fmt.Errorf("NewLazadaClient: %w", err)
	}

	resp, err := client.Refresh(refreshToken)
	if err != nil {
		return "", fmt.Errorf("token refresh: %w", err)
	}

	lazInfo, err := s.store.GetLazadaPlatformByShopID(
		shopID,
	)
	if err != nil {
		return "", fmt.Errorf("GetLazadaPlatformByShopID: %w", err)
	}

	err = s.store.UpdateLazadaPlatform(lazInfo.ID, resp)
	if err != nil {
		return "", fmt.Errorf("UpdateLazadaPlatform: %w", err)
	}

	return resp.AccessToken, nil
}

func (s *Service) GetValidTokenByShopID(shopID string) (string, error) {
	token, err := s.store.GetValidAccessTokenByShopID(shopID)

	//Check for refresh token?
	if errors.Is(err, pgx.ErrNoRows) {
		token, err = s.RefreshAndGetTokenByShopID(shopID)
		if err != nil {
			return "", fmt.Errorf("RefreshAndGetTokenByShopID: %w", err)
		}
		return token, nil

	} else if err != nil {
		return "", fmt.Errorf("GetValidTokenByShopID: %w", err)
	}

	return token, nil
}

func (s *Service) RetrieveProductFromRedis(keyID string, index int) (*ProductsResponseProducts, error) {
	i := strconv.Itoa(index)
	ctx := context.Background()

	var rcv string
	mb := radix.Maybe{Rcv: &rcv}

	if err := s.server.RedisClient.Do(ctx, radix.Cmd(&mb, "HGET", keyID, i)); err != nil {
		return nil, fmt.Errorf("redis Do: %w", err)
	}

	if mb.Null || mb.Empty {
		return nil, ErrProductIndexNotFound
	}

	var response ProductsResponseProducts

	err := json.Unmarshal([]byte(*mb.Rcv.(*string)), &response)
	if err != nil {
		return nil, fmt.Errorf("product response unmarshal: %w", err)
	}

	return &response, nil
}

func (s *Service) ExpireRedisProducts(keyID string) error {
	ctx := context.Background()
	err := s.server.RedisClient.Do(ctx, radix.Cmd(nil, "DEL", keyID))

	if err != nil {
		s.server.Logger.Errorw("Error in auth handler function (Redis)",
			"error", err,
		)
		return err
	}

	return nil
}

func (s *Service) SaveProductToRedis(keyID string, index int, p ProductsResponseProducts) error {
	bytes, err := json.Marshal(p)
	if err != nil {
		return err
	}

	i := strconv.Itoa(index)
	m := map[string]string{i: string(bytes)}
	ctx := context.Background()
	if err := s.server.RedisClient.Do(ctx, radix.FlatCmd(nil, "HMSET", keyID, m)); err != nil {
		return err
	}
	return nil
}

func (s *Service) FetchProductsFromLazadaToRedis(shopID string) (string, int, error) {
	token := ""
	token, err := s.GetValidTokenByShopID(shopID)
	if err != nil {
		return "", 0, fmt.Errorf("GetValidTokenByShopID: %w", err)
	}

	client, err := s.NewLazadaClient(token)
	if err != nil {
		return "", 0, fmt.Errorf("NewLazadaClient: %w", err)
	}

	response, err := client.QueryProducts(QueryProductsParams{
		Limit:  50,
		Offset: 0,
	})
	if err != nil {
		return "", 0, fmt.Errorf("QueryProducts: %w", err)
	}

	keyID := redisKeyPrefix + shopID

	for i, product := range response.Data.Products {
		if err = s.SaveProductToRedis(keyID, i, product); err != nil {
			return "", 0, fmt.Errorf("SaveProductToRedis: %w", err)
		}
	}

	var fetchesRequired int
	if response.Data.TotalProducts > 50 {
		remainingProducts := (float64)(response.Data.TotalProducts - 50)
		fetchesRequired = (int)(math.Ceil(remainingProducts / 50))
	} else {
		return keyID, response.Data.TotalProducts, nil
	}

	for i := 1; i <= fetchesRequired; i++ {
		response, err = client.QueryProducts(QueryProductsParams{
			Limit:  50,
			Offset: 50 * i,
		})
		if err != nil {
			return "", 0, fmt.Errorf("QueryProducts: %w", err)
		}

		for j, product := range response.Data.Products {
			index := i*50 + j
			if err := s.SaveProductToRedis(keyID, index, product); err != nil {
				return "", 0, fmt.Errorf("SaveProductToRedis: %w", err)
			}
		}
	}

	return keyID, response.Data.TotalProducts, nil
}

func (s *Service) SaveOrUpdateLazadaProduct(p *ProductsResponseProducts, shopID string) (*sqlc.ProductsLazada, *sqlc.ProductsAttributeLazada, *sqlc.ProductsSkuLazada, error) {
	return s.store.SaveOrUpdateLazadaProduct(p, shopID)
}
