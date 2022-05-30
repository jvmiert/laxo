// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.13.0
// source: platform_lazada.sql

package sqlc

import (
	"context"

	null "gopkg.in/guregu/null.v4"
)

const createLazadaPlatform = `-- name: CreateLazadaPlatform :one
INSERT INTO platform_lazada (
  shop_id, access_token, country, refresh_token, account_platform, account,
  user_id_vn, seller_id_vn, short_code_vn, refresh_expires_in,
  access_expires_in
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
)
RETURNING id, shop_id, access_token, country, refresh_token, account_platform, account, user_id_vn, seller_id_vn, short_code_vn, refresh_expires_in, access_expires_in, created
`

type CreateLazadaPlatformParams struct {
	ShopID           string    `json:"shopID"`
	AccessToken      string    `json:"accessToken"`
	Country          string    `json:"country"`
	RefreshToken     string    `json:"refreshToken"`
	AccountPlatform  string    `json:"accountPlatform"`
	Account          string    `json:"account"`
	UserIDVn         string    `json:"userIDVn"`
	SellerIDVn       string    `json:"sellerIDVn"`
	ShortCodeVn      string    `json:"shortCodeVn"`
	RefreshExpiresIn null.Time `json:"refreshExpiresIn"`
	AccessExpiresIn  null.Time `json:"accessExpiresIn"`
}

func (q *Queries) CreateLazadaPlatform(ctx context.Context, arg CreateLazadaPlatformParams) (PlatformLazada, error) {
	row := q.db.QueryRow(ctx, createLazadaPlatform,
		arg.ShopID,
		arg.AccessToken,
		arg.Country,
		arg.RefreshToken,
		arg.AccountPlatform,
		arg.Account,
		arg.UserIDVn,
		arg.SellerIDVn,
		arg.ShortCodeVn,
		arg.RefreshExpiresIn,
		arg.AccessExpiresIn,
	)
	var i PlatformLazada
	err := row.Scan(
		&i.ID,
		&i.ShopID,
		&i.AccessToken,
		&i.Country,
		&i.RefreshToken,
		&i.AccountPlatform,
		&i.Account,
		&i.UserIDVn,
		&i.SellerIDVn,
		&i.ShortCodeVn,
		&i.RefreshExpiresIn,
		&i.AccessExpiresIn,
		&i.Created,
	)
	return i, err
}

const getLazadaPlatformByShopID = `-- name: GetLazadaPlatformByShopID :one
SELECT id, shop_id, access_token, country, refresh_token, account_platform, account, user_id_vn, seller_id_vn, short_code_vn, refresh_expires_in, access_expires_in, created FROM platform_lazada
WHERE shop_id = $1
LIMIT 1
`

func (q *Queries) GetLazadaPlatformByShopID(ctx context.Context, shopID string) (PlatformLazada, error) {
	row := q.db.QueryRow(ctx, getLazadaPlatformByShopID, shopID)
	var i PlatformLazada
	err := row.Scan(
		&i.ID,
		&i.ShopID,
		&i.AccessToken,
		&i.Country,
		&i.RefreshToken,
		&i.AccountPlatform,
		&i.Account,
		&i.UserIDVn,
		&i.SellerIDVn,
		&i.ShortCodeVn,
		&i.RefreshExpiresIn,
		&i.AccessExpiresIn,
		&i.Created,
	)
	return i, err
}

const getValidAccessTokenByShopID = `-- name: GetValidAccessTokenByShopID :one
SELECT access_token FROM platform_lazada
WHERE shop_id = $1 AND access_expires_in > NOW()
LIMIT 1
`

func (q *Queries) GetValidAccessTokenByShopID(ctx context.Context, shopID string) (string, error) {
	row := q.db.QueryRow(ctx, getValidAccessTokenByShopID, shopID)
	var access_token string
	err := row.Scan(&access_token)
	return access_token, err
}

const getValidRefreshTokenByShopID = `-- name: GetValidRefreshTokenByShopID :one
SELECT refresh_token FROM platform_lazada
WHERE shop_id = $1 AND refresh_expires_in > NOW()
LIMIT 1
`

func (q *Queries) GetValidRefreshTokenByShopID(ctx context.Context, shopID string) (string, error) {
	row := q.db.QueryRow(ctx, getValidRefreshTokenByShopID, shopID)
	var refresh_token string
	err := row.Scan(&refresh_token)
	return refresh_token, err
}

const updateLazadaPlatform = `-- name: UpdateLazadaPlatform :exec
UPDATE platform_lazada SET access_token = $1, refresh_token = $2,
refresh_expires_in = $3, access_expires_in = $4
WHERE id = $5
`

type UpdateLazadaPlatformParams struct {
	AccessToken      string    `json:"accessToken"`
	RefreshToken     string    `json:"refreshToken"`
	RefreshExpiresIn null.Time `json:"refreshExpiresIn"`
	AccessExpiresIn  null.Time `json:"accessExpiresIn"`
	ID               string    `json:"id"`
}

func (q *Queries) UpdateLazadaPlatform(ctx context.Context, arg UpdateLazadaPlatformParams) error {
	_, err := q.db.Exec(ctx, updateLazadaPlatform,
		arg.AccessToken,
		arg.RefreshToken,
		arg.RefreshExpiresIn,
		arg.AccessExpiresIn,
		arg.ID,
	)
	return err
}
