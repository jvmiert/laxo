// Code generated by sqlc. DO NOT EDIT.
// source: shop.sql

package sqlc

import (
	"context"
	"database/sql"
)

const createShop = `-- name: CreateShop :one
INSERT INTO shops (
  shop_name, user_id
) VALUES (
  $1, $2
)
RETURNING id, user_id, shop_name, created, last_update
`

type CreateShopParams struct {
	ShopName string `json:"shopName"`
	UserID   string `json:"userID"`
}

func (q *Queries) CreateShop(ctx context.Context, arg CreateShopParams) (Shop, error) {
	row := q.db.QueryRow(ctx, createShop, arg.ShopName, arg.UserID)
	var i Shop
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.ShopName,
		&i.Created,
		&i.LastUpdate,
	)
	return i, err
}

const getShopByID = `-- name: GetShopByID :one
SELECT id, user_id, shop_name, created, last_update FROM shops
WHERE id = $1
LIMIT 1
`

func (q *Queries) GetShopByID(ctx context.Context, id string) (Shop, error) {
	row := q.db.QueryRow(ctx, getShopByID, id)
	var i Shop
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.ShopName,
		&i.Created,
		&i.LastUpdate,
	)
	return i, err
}

const getShopsByUserID = `-- name: GetShopsByUserID :many
SELECT id, user_id, shop_name, created, last_update FROM shops
WHERE user_id = $1
`

func (q *Queries) GetShopsByUserID(ctx context.Context, userID string) ([]Shop, error) {
	rows, err := q.db.Query(ctx, getShopsByUserID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Shop
	for rows.Next() {
		var i Shop
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.ShopName,
			&i.Created,
			&i.LastUpdate,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getShopsPlatformsByUserID = `-- name: GetShopsPlatformsByUserID :many
SELECT shops.id,
       shops.user_id,
       shops.shop_name,
       shops.created,
       shops_platforms.id as platform_id,
       shops_platforms.platform_name as platform_name,
       shops_platforms.created as platform_created
FROM shops
LEFT JOIN shops_platforms ON (shops.id = shops_platforms.shop_id)
WHERE user_id = $1
`

type GetShopsPlatformsByUserIDRow struct {
	ID              string         `json:"id"`
	UserID          string         `json:"userID"`
	ShopName        string         `json:"shopName"`
	Created         sql.NullTime   `json:"created"`
	PlatformID      sql.NullString `json:"platformID"`
	PlatformName    sql.NullString `json:"platformName"`
	PlatformCreated sql.NullTime   `json:"platformCreated"`
}

func (q *Queries) GetShopsPlatformsByUserID(ctx context.Context, userID string) ([]GetShopsPlatformsByUserIDRow, error) {
	rows, err := q.db.Query(ctx, getShopsPlatformsByUserID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetShopsPlatformsByUserIDRow
	for rows.Next() {
		var i GetShopsPlatformsByUserIDRow
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.ShopName,
			&i.Created,
			&i.PlatformID,
			&i.PlatformName,
			&i.PlatformCreated,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
