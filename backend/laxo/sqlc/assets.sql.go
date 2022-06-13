// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.13.0
// source: assets.sql

package sqlc

import (
	"context"

	null "gopkg.in/guregu/null.v4"
)

const createAsset = `-- name: CreateAsset :one
INSERT INTO assets (
  shop_id, murmur_hash, original_filename, extension,
  file_size, width, height
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
)
RETURNING id, shop_id, murmur_hash, original_filename, extension, file_size, width, height
`

type CreateAssetParams struct {
	ShopID           string      `json:"shopID"`
	MurmurHash       string      `json:"murmurHash"`
	OriginalFilename null.String `json:"originalFilename"`
	Extension        null.String `json:"extension"`
	FileSize         null.Int    `json:"fileSize"`
	Width            null.Int    `json:"width"`
	Height           null.Int    `json:"height"`
}

func (q *Queries) CreateAsset(ctx context.Context, arg CreateAssetParams) (Asset, error) {
	row := q.db.QueryRow(ctx, createAsset,
		arg.ShopID,
		arg.MurmurHash,
		arg.OriginalFilename,
		arg.Extension,
		arg.FileSize,
		arg.Width,
		arg.Height,
	)
	var i Asset
	err := row.Scan(
		&i.ID,
		&i.ShopID,
		&i.MurmurHash,
		&i.OriginalFilename,
		&i.Extension,
		&i.FileSize,
		&i.Width,
		&i.Height,
	)
	return i, err
}

const createProductMedia = `-- name: CreateProductMedia :one
INSERT INTO products_media (
  product_id, asset_id,
  image_order, status
) VALUES (
  $1, $2, $3, $4
)
RETURNING product_id, asset_id, image_order, status
`

type CreateProductMediaParams struct {
	ProductID  string   `json:"productID"`
	AssetID    string   `json:"assetID"`
	ImageOrder null.Int `json:"imageOrder"`
	Status     string   `json:"status"`
}

func (q *Queries) CreateProductMedia(ctx context.Context, arg CreateProductMediaParams) (ProductsMedia, error) {
	row := q.db.QueryRow(ctx, createProductMedia,
		arg.ProductID,
		arg.AssetID,
		arg.ImageOrder,
		arg.Status,
	)
	var i ProductsMedia
	err := row.Scan(
		&i.ProductID,
		&i.AssetID,
		&i.ImageOrder,
		&i.Status,
	)
	return i, err
}

const deleteProductMedia = `-- name: DeleteProductMedia :exec
DELETE FROM products_media
WHERE product_id = $1 AND asset_id = $2
`

type DeleteProductMediaParams struct {
	ProductID string `json:"productID"`
	AssetID   string `json:"assetID"`
}

func (q *Queries) DeleteProductMedia(ctx context.Context, arg DeleteProductMediaParams) error {
	_, err := q.db.Exec(ctx, deleteProductMedia, arg.ProductID, arg.AssetID)
	return err
}

const getAssetByID = `-- name: GetAssetByID :one
SELECT id, shop_id, murmur_hash, original_filename, extension, file_size, width, height FROM assets
WHERE id = $1
LIMIT 1
`

func (q *Queries) GetAssetByID(ctx context.Context, id string) (Asset, error) {
	row := q.db.QueryRow(ctx, getAssetByID, id)
	var i Asset
	err := row.Scan(
		&i.ID,
		&i.ShopID,
		&i.MurmurHash,
		&i.OriginalFilename,
		&i.Extension,
		&i.FileSize,
		&i.Width,
		&i.Height,
	)
	return i, err
}

const getAssetByMurmur = `-- name: GetAssetByMurmur :one
SELECT id, shop_id, murmur_hash, original_filename, extension, file_size, width, height FROM assets
WHERE murmur_hash = $1 AND shop_id = $2
LIMIT 1
`

type GetAssetByMurmurParams struct {
	MurmurHash string `json:"murmurHash"`
	ShopID     string `json:"shopID"`
}

func (q *Queries) GetAssetByMurmur(ctx context.Context, arg GetAssetByMurmurParams) (Asset, error) {
	row := q.db.QueryRow(ctx, getAssetByMurmur, arg.MurmurHash, arg.ShopID)
	var i Asset
	err := row.Scan(
		&i.ID,
		&i.ShopID,
		&i.MurmurHash,
		&i.OriginalFilename,
		&i.Extension,
		&i.FileSize,
		&i.Width,
		&i.Height,
	)
	return i, err
}

const getProductMedia = `-- name: GetProductMedia :one
SELECT product_id, asset_id, image_order, status FROM products_media
WHERE product_id = $1 AND asset_id = $2
`

type GetProductMediaParams struct {
	ProductID string `json:"productID"`
	AssetID   string `json:"assetID"`
}

func (q *Queries) GetProductMedia(ctx context.Context, arg GetProductMediaParams) (ProductsMedia, error) {
	row := q.db.QueryRow(ctx, getProductMedia, arg.ProductID, arg.AssetID)
	var i ProductsMedia
	err := row.Scan(
		&i.ProductID,
		&i.AssetID,
		&i.ImageOrder,
		&i.Status,
	)
	return i, err
}

const updateAsset = `-- name: UpdateAsset :one
UPDATE assets
SET
 shop_id = coalesce($1, shop_id),
 murmur_hash = coalesce($2, murmur_hash),
 original_filename = coalesce($3, original_filename),
 extension = coalesce($4, extension),
 file_size = coalesce($5, file_size),
 width = coalesce($6, width),
 height = coalesce($7, height)
WHERE id = $8
RETURNING id, shop_id, murmur_hash, original_filename, extension, file_size, width, height
`

type UpdateAssetParams struct {
	ShopID           null.String `json:"shopID"`
	MurmurHash       null.String `json:"murmurHash"`
	OriginalFilename null.String `json:"originalFilename"`
	Extension        null.String `json:"extension"`
	FileSize         null.Int    `json:"fileSize"`
	Width            null.Int    `json:"width"`
	Height           null.Int    `json:"height"`
	ID               string      `json:"id"`
}

func (q *Queries) UpdateAsset(ctx context.Context, arg UpdateAssetParams) (Asset, error) {
	row := q.db.QueryRow(ctx, updateAsset,
		arg.ShopID,
		arg.MurmurHash,
		arg.OriginalFilename,
		arg.Extension,
		arg.FileSize,
		arg.Width,
		arg.Height,
		arg.ID,
	)
	var i Asset
	err := row.Scan(
		&i.ID,
		&i.ShopID,
		&i.MurmurHash,
		&i.OriginalFilename,
		&i.Extension,
		&i.FileSize,
		&i.Width,
		&i.Height,
	)
	return i, err
}

const updateProductMedia = `-- name: UpdateProductMedia :one
UPDATE products_media
SET
 image_order = coalesce($1, image_order),
 status = coalesce($2, status)
WHERE product_id = $3 AND asset_id = $4
RETURNING product_id, asset_id, image_order, status
`

type UpdateProductMediaParams struct {
	ImageOrder null.Int    `json:"imageOrder"`
	Status     null.String `json:"status"`
	ProductID  string      `json:"productID"`
	AssetID    string      `json:"assetID"`
}

func (q *Queries) UpdateProductMedia(ctx context.Context, arg UpdateProductMediaParams) (ProductsMedia, error) {
	row := q.db.QueryRow(ctx, updateProductMedia,
		arg.ImageOrder,
		arg.Status,
		arg.ProductID,
		arg.AssetID,
	)
	var i ProductsMedia
	err := row.Scan(
		&i.ProductID,
		&i.AssetID,
		&i.ImageOrder,
		&i.Status,
	)
	return i, err
}
