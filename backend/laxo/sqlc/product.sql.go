// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.13.0
// source: product.sql

package sqlc

import (
	"context"

	"github.com/jackc/pgtype"
	null "gopkg.in/guregu/null.v4"
)

const checkProductOwner = `-- name: CheckProductOwner :one
SELECT products.id FROM products
WHERE products.id = $1 AND products.shop_id = $2
`

type CheckProductOwnerParams struct {
	ID     string `json:"id"`
	ShopID string `json:"shopID"`
}

func (q *Queries) CheckProductOwner(ctx context.Context, arg CheckProductOwnerParams) (string, error) {
	row := q.db.QueryRow(ctx, checkProductOwner, arg.ID, arg.ShopID)
	var id string
	err := row.Scan(&id)
	return id, err
}

const createProduct = `-- name: CreateProduct :one
INSERT INTO products (
  name, description, msku, selling_price, cost_price, shop_id,
  media_id, updated, description_slate
) VALUES (
  $1, $2, $3, $4, $5, $6, $7,
  $8, $9
)
RETURNING id, name, description, description_slate, msku, selling_price, cost_price, shop_id, media_id, created, updated
`

type CreateProductParams struct {
	Name             null.String    `json:"name"`
	Description      null.String    `json:"description"`
	Msku             null.String    `json:"msku"`
	SellingPrice     pgtype.Numeric `json:"sellingPrice"`
	CostPrice        pgtype.Numeric `json:"costPrice"`
	ShopID           string         `json:"shopID"`
	MediaID          null.String    `json:"mediaID"`
	Updated          null.Time      `json:"updated"`
	DescriptionSlate null.String    `json:"descriptionSlate"`
}

func (q *Queries) CreateProduct(ctx context.Context, arg CreateProductParams) (Product, error) {
	row := q.db.QueryRow(ctx, createProduct,
		arg.Name,
		arg.Description,
		arg.Msku,
		arg.SellingPrice,
		arg.CostPrice,
		arg.ShopID,
		arg.MediaID,
		arg.Updated,
		arg.DescriptionSlate,
	)
	var i Product
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Description,
		&i.DescriptionSlate,
		&i.Msku,
		&i.SellingPrice,
		&i.CostPrice,
		&i.ShopID,
		&i.MediaID,
		&i.Created,
		&i.Updated,
	)
	return i, err
}

const createProductPlatform = `-- name: CreateProductPlatform :one
INSERT INTO products_platform (
  product_id, products_lazada_id
) VALUES (
  $1, $2
)
RETURNING product_id, products_lazada_id, sync_lazada
`

type CreateProductPlatformParams struct {
	ProductID        string      `json:"productID"`
	ProductsLazadaID null.String `json:"productsLazadaID"`
}

func (q *Queries) CreateProductPlatform(ctx context.Context, arg CreateProductPlatformParams) (ProductsPlatform, error) {
	row := q.db.QueryRow(ctx, createProductPlatform, arg.ProductID, arg.ProductsLazadaID)
	var i ProductsPlatform
	err := row.Scan(&i.ProductID, &i.ProductsLazadaID, &i.SyncLazada)
	return i, err
}

const getProductAssetsByProductID = `-- name: GetProductAssetsByProductID :many
SELECT assets.id, assets.shop_id, assets.murmur_hash, assets.original_filename, assets.extension, assets.file_size, assets.width, assets.height, assets.created, products_media.image_order, products_media.status FROM products
LEFT JOIN products_media ON products_media.product_id = products.id
LEFT JOIN assets ON assets.id = products_media.asset_id
WHERE products.id = $1 AND products.shop_id = $2
ORDER BY products_media.image_order, products_media.product_id, assets.id
`

type GetProductAssetsByProductIDParams struct {
	ID     string `json:"id"`
	ShopID string `json:"shopID"`
}

type GetProductAssetsByProductIDRow struct {
	ID               null.String `json:"id"`
	ShopID           null.String `json:"shopID"`
	MurmurHash       null.String `json:"murmurHash"`
	OriginalFilename null.String `json:"originalFilename"`
	Extension        null.String `json:"extension"`
	FileSize         null.Int    `json:"fileSize"`
	Width            null.Int    `json:"width"`
	Height           null.Int    `json:"height"`
	Created          null.Time   `json:"created"`
	ImageOrder       null.Int    `json:"imageOrder"`
	Status           null.String `json:"status"`
}

func (q *Queries) GetProductAssetsByProductID(ctx context.Context, arg GetProductAssetsByProductIDParams) ([]GetProductAssetsByProductIDRow, error) {
	rows, err := q.db.Query(ctx, getProductAssetsByProductID, arg.ID, arg.ShopID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetProductAssetsByProductIDRow
	for rows.Next() {
		var i GetProductAssetsByProductIDRow
		if err := rows.Scan(
			&i.ID,
			&i.ShopID,
			&i.MurmurHash,
			&i.OriginalFilename,
			&i.Extension,
			&i.FileSize,
			&i.Width,
			&i.Height,
			&i.Created,
			&i.ImageOrder,
			&i.Status,
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

const getProductByID = `-- name: GetProductByID :one
SELECT id, name, description, description_slate, msku, selling_price, cost_price, shop_id, media_id, created, updated FROM products
WHERE id = $1
LIMIT 1
`

func (q *Queries) GetProductByID(ctx context.Context, id string) (Product, error) {
	row := q.db.QueryRow(ctx, getProductByID, id)
	var i Product
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Description,
		&i.DescriptionSlate,
		&i.Msku,
		&i.SellingPrice,
		&i.CostPrice,
		&i.ShopID,
		&i.MediaID,
		&i.Created,
		&i.Updated,
	)
	return i, err
}

const getProductByProductMSKU = `-- name: GetProductByProductMSKU :one
SELECT id, name, description, description_slate, msku, selling_price, cost_price, shop_id, media_id, created, updated FROM products
WHERE msku = $1 AND shop_id = $2
LIMIT 1
`

type GetProductByProductMSKUParams struct {
	Msku   null.String `json:"msku"`
	ShopID string      `json:"shopID"`
}

func (q *Queries) GetProductByProductMSKU(ctx context.Context, arg GetProductByProductMSKUParams) (Product, error) {
	row := q.db.QueryRow(ctx, getProductByProductMSKU, arg.Msku, arg.ShopID)
	var i Product
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Description,
		&i.DescriptionSlate,
		&i.Msku,
		&i.SellingPrice,
		&i.CostPrice,
		&i.ShopID,
		&i.MediaID,
		&i.Created,
		&i.Updated,
	)
	return i, err
}

const getProductDetailsByID = `-- name: GetProductDetailsByID :one
SELECT products.id, products.name, products.description, products.description_slate, products.msku, products.selling_price, products.cost_price, products.shop_id, products.media_id, products.created, products.updated,
  products_lazada.lazada_id as lazada_id,
  products_sku_lazada.url as lazada_url,
  products_attribute_lazada.name as lazada_name,
  products_sku_lazada.sku_id as lazada_platform_sku,
  products_sku_lazada.seller_sku as lazada_seller_sku,
  products_lazada.status as lazada_status,
  products_platform.sync_lazada as lazada_sync_status
FROM products
LEFT JOIN products_media ON products_media.product_id = products.id
LEFT JOIN assets ON assets.id = products_media.asset_id
LEFT JOIN products_platform ON products_platform.product_id = products.id
LEFT JOIN products_lazada ON products_platform.products_lazada_id = products_lazada.id
LEFT JOIN products_sku_lazada ON products_sku_lazada.product_id = products_lazada.id
LEFT JOIN products_attribute_lazada ON products_attribute_lazada.product_id = products_lazada.id
WHERE products.id = $1 AND products.shop_id = $2
`

type GetProductDetailsByIDParams struct {
	ID     string `json:"id"`
	ShopID string `json:"shopID"`
}

type GetProductDetailsByIDRow struct {
	ID                string         `json:"id"`
	Name              null.String    `json:"name"`
	Description       null.String    `json:"description"`
	DescriptionSlate  null.String    `json:"descriptionSlate"`
	Msku              null.String    `json:"msku"`
	SellingPrice      pgtype.Numeric `json:"sellingPrice"`
	CostPrice         pgtype.Numeric `json:"costPrice"`
	ShopID            string         `json:"shopID"`
	MediaID           null.String    `json:"mediaID"`
	Created           null.Time      `json:"created"`
	Updated           null.Time      `json:"updated"`
	LazadaID          null.Int       `json:"lazadaID"`
	LazadaUrl         null.String    `json:"lazadaUrl"`
	LazadaName        null.String    `json:"lazadaName"`
	LazadaPlatformSku null.Int       `json:"lazadaPlatformSku"`
	LazadaSellerSku   null.String    `json:"lazadaSellerSku"`
	LazadaStatus      null.String    `json:"lazadaStatus"`
	LazadaSyncStatus  null.Bool      `json:"lazadaSyncStatus"`
}

func (q *Queries) GetProductDetailsByID(ctx context.Context, arg GetProductDetailsByIDParams) (GetProductDetailsByIDRow, error) {
	row := q.db.QueryRow(ctx, getProductDetailsByID, arg.ID, arg.ShopID)
	var i GetProductDetailsByIDRow
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Description,
		&i.DescriptionSlate,
		&i.Msku,
		&i.SellingPrice,
		&i.CostPrice,
		&i.ShopID,
		&i.MediaID,
		&i.Created,
		&i.Updated,
		&i.LazadaID,
		&i.LazadaUrl,
		&i.LazadaName,
		&i.LazadaPlatformSku,
		&i.LazadaSellerSku,
		&i.LazadaStatus,
		&i.LazadaSyncStatus,
	)
	return i, err
}

const getProductPlatformByLazadaID = `-- name: GetProductPlatformByLazadaID :one
SELECT product_id, products_lazada_id, sync_lazada FROM products_platform
WHERE products_lazada_id = $1
LIMIT 1
`

func (q *Queries) GetProductPlatformByLazadaID(ctx context.Context, productsLazadaID null.String) (ProductsPlatform, error) {
	row := q.db.QueryRow(ctx, getProductPlatformByLazadaID, productsLazadaID)
	var i ProductsPlatform
	err := row.Scan(&i.ProductID, &i.ProductsLazadaID, &i.SyncLazada)
	return i, err
}

const getProductPlatformByProductID = `-- name: GetProductPlatformByProductID :one
SELECT product_id, products_lazada_id, sync_lazada FROM products_platform
WHERE product_id = $1
LIMIT 1
`

func (q *Queries) GetProductPlatformByProductID(ctx context.Context, productID string) (ProductsPlatform, error) {
	row := q.db.QueryRow(ctx, getProductPlatformByProductID, productID)
	var i ProductsPlatform
	err := row.Scan(&i.ProductID, &i.ProductsLazadaID, &i.SyncLazada)
	return i, err
}

const getProductsByNameOrSKU = `-- name: GetProductsByNameOrSKU :many
SELECT
  c.count,
  COALESCE(p.id, ''), p.name, p.description, p.msku, p.selling_price, p.cost_price,
  COALESCE(p.shop_id, ''), p.media_id, p.created, p.updated, media_id_list,
  COALESCE(p.lazada_id, 0), lazada_url, lazada_name, lazada_platform_sku,
  COALESCE(lazada_seller_sku, ''), p.lazada_status
FROM
(
  SELECT COUNT(*) AS COUNT
  FROM products
  WHERE products.shop_id = $1 AND (products.name ILIKE $2 OR products.msku ILIKE $3)
) as c
LEFT JOIN (
  SELECT products.id, products.name, products.description, products.description_slate, products.msku, products.selling_price, products.cost_price, products.shop_id, products.media_id, products.created, products.updated,
    STRING_AGG(CONCAT(assets.id, assets.extension), ',' order by products_media.image_order) as media_id_list,
    products_lazada.lazada_id as lazada_id,
    products_sku_lazada.url as lazada_url,
    products_attribute_lazada.name as lazada_name,
    products_sku_lazada.sku_id as lazada_platform_sku,
    products_sku_lazada.seller_sku as lazada_seller_sku,
    products_lazada.status as lazada_status
  FROM products
  LEFT JOIN products_media ON products_media.product_id = products.id
  LEFT JOIN assets ON assets.id = products_media.asset_id
  LEFT JOIN products_platform ON products_platform.product_id = products.id
  LEFT JOIN products_lazada ON products_platform.products_lazada_id = products_lazada.id
  LEFT JOIN products_sku_lazada ON products_sku_lazada.product_id = products_lazada.id
  LEFT JOIN products_attribute_lazada ON products_attribute_lazada.product_id = products_lazada.id
  WHERE products.shop_id = $1 AND (products.name ILIKE $2 OR products.msku ILIKE $3)
  GROUP BY products.id, products_lazada.id, products_sku_lazada.id, products_attribute_lazada.id
  ORDER BY UPPER(products.name) COLLATE "vi-VN-x-icu"
  LIMIT $4 OFFSET $5
) as p
ON true
`

type GetProductsByNameOrSKUParams struct {
	ShopID string      `json:"shopID"`
	Name   null.String `json:"name"`
	Msku   null.String `json:"msku"`
	Limit  int32       `json:"limit"`
	Offset int32       `json:"offset"`
}

type GetProductsByNameOrSKURow struct {
	Count             int64          `json:"count"`
	ID                string         `json:"id"`
	Name              null.String    `json:"name"`
	Description       null.String    `json:"description"`
	Msku              null.String    `json:"msku"`
	SellingPrice      pgtype.Numeric `json:"sellingPrice"`
	CostPrice         pgtype.Numeric `json:"costPrice"`
	ShopID            string         `json:"shopID"`
	MediaID           null.String    `json:"mediaID"`
	Created           null.Time      `json:"created"`
	Updated           null.Time      `json:"updated"`
	MediaIDList       []byte         `json:"mediaIDList"`
	LazadaID          int64          `json:"lazadaID"`
	LazadaUrl         null.String    `json:"lazadaUrl"`
	LazadaName        null.String    `json:"lazadaName"`
	LazadaPlatformSku null.Int       `json:"lazadaPlatformSku"`
	LazadaSellerSku   string         `json:"lazadaSellerSku"`
	LazadaStatus      null.String    `json:"lazadaStatus"`
}

func (q *Queries) GetProductsByNameOrSKU(ctx context.Context, arg GetProductsByNameOrSKUParams) ([]GetProductsByNameOrSKURow, error) {
	rows, err := q.db.Query(ctx, getProductsByNameOrSKU,
		arg.ShopID,
		arg.Name,
		arg.Msku,
		arg.Limit,
		arg.Offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetProductsByNameOrSKURow
	for rows.Next() {
		var i GetProductsByNameOrSKURow
		if err := rows.Scan(
			&i.Count,
			&i.ID,
			&i.Name,
			&i.Description,
			&i.Msku,
			&i.SellingPrice,
			&i.CostPrice,
			&i.ShopID,
			&i.MediaID,
			&i.Created,
			&i.Updated,
			&i.MediaIDList,
			&i.LazadaID,
			&i.LazadaUrl,
			&i.LazadaName,
			&i.LazadaPlatformSku,
			&i.LazadaSellerSku,
			&i.LazadaStatus,
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

const getProductsByShopID = `-- name: GetProductsByShopID :many
SELECT
  c.count,
  COALESCE(p.id, ''), p.name, p.description, p.msku, p.selling_price, p.cost_price,
  COALESCE(p.shop_id, ''), p.media_id, p.created, p.updated, media_id_list,
  COALESCE(p.lazada_id, 0), lazada_url, lazada_name, lazada_platform_sku,
  COALESCE(lazada_seller_sku, ''), p.lazada_status
FROM
(
  SELECT COUNT(*) AS COUNT
  FROM products
  WHERE products.shop_id = $1
) as c
LEFT JOIN (
  SELECT products.id, products.name, products.description, products.description_slate, products.msku, products.selling_price, products.cost_price, products.shop_id, products.media_id, products.created, products.updated,
    STRING_AGG(CONCAT(assets.id, assets.extension), ',' order by products_media.image_order) as media_id_list,
    products_lazada.lazada_id as lazada_id,
    products_sku_lazada.url as lazada_url,
    products_attribute_lazada.name as lazada_name,
    products_sku_lazada.sku_id as lazada_platform_sku,
    products_sku_lazada.seller_sku as lazada_seller_sku,
    products_lazada.status as lazada_status
  FROM products
  LEFT JOIN products_media ON products_media.product_id = products.id
  LEFT JOIN assets ON assets.id = products_media.asset_id
  LEFT JOIN products_platform ON products_platform.product_id = products.id
  LEFT JOIN products_lazada ON products_platform.products_lazada_id = products_lazada.id
  LEFT JOIN products_sku_lazada ON products_sku_lazada.product_id = products_lazada.id
  LEFT JOIN products_attribute_lazada ON products_attribute_lazada.product_id = products_lazada.id
  WHERE products.shop_id = $1
  GROUP BY products.id, products_lazada.id, products_sku_lazada.id, products_attribute_lazada.id
  ORDER BY UPPER(products.name) COLLATE "vi-VN-x-icu"
  LIMIT $2 OFFSET $3
) as p
ON true
`

type GetProductsByShopIDParams struct {
	ShopID string `json:"shopID"`
	Limit  int32  `json:"limit"`
	Offset int32  `json:"offset"`
}

type GetProductsByShopIDRow struct {
	Count             int64          `json:"count"`
	ID                string         `json:"id"`
	Name              null.String    `json:"name"`
	Description       null.String    `json:"description"`
	Msku              null.String    `json:"msku"`
	SellingPrice      pgtype.Numeric `json:"sellingPrice"`
	CostPrice         pgtype.Numeric `json:"costPrice"`
	ShopID            string         `json:"shopID"`
	MediaID           null.String    `json:"mediaID"`
	Created           null.Time      `json:"created"`
	Updated           null.Time      `json:"updated"`
	MediaIDList       []byte         `json:"mediaIDList"`
	LazadaID          int64          `json:"lazadaID"`
	LazadaUrl         null.String    `json:"lazadaUrl"`
	LazadaName        null.String    `json:"lazadaName"`
	LazadaPlatformSku null.Int       `json:"lazadaPlatformSku"`
	LazadaSellerSku   string         `json:"lazadaSellerSku"`
	LazadaStatus      null.String    `json:"lazadaStatus"`
}

func (q *Queries) GetProductsByShopID(ctx context.Context, arg GetProductsByShopIDParams) ([]GetProductsByShopIDRow, error) {
	rows, err := q.db.Query(ctx, getProductsByShopID, arg.ShopID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetProductsByShopIDRow
	for rows.Next() {
		var i GetProductsByShopIDRow
		if err := rows.Scan(
			&i.Count,
			&i.ID,
			&i.Name,
			&i.Description,
			&i.Msku,
			&i.SellingPrice,
			&i.CostPrice,
			&i.ShopID,
			&i.MediaID,
			&i.Created,
			&i.Updated,
			&i.MediaIDList,
			&i.LazadaID,
			&i.LazadaUrl,
			&i.LazadaName,
			&i.LazadaPlatformSku,
			&i.LazadaSellerSku,
			&i.LazadaStatus,
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

const updateProduct = `-- name: UpdateProduct :one
UPDATE products
SET
 name = coalesce($1, name),
 description = coalesce($2, description),
 description_slate = coalesce($3, description),
 selling_price = coalesce($4, selling_price),
 cost_price = coalesce($5, cost_price),
 media_id = coalesce($6, media_id),
 updated = coalesce($7, updated)
WHERE id = $8 AND shop_id = $9
RETURNING id, name, description, description_slate, msku, selling_price, cost_price, shop_id, media_id, created, updated
`

type UpdateProductParams struct {
	Name             null.String    `json:"name"`
	Description      null.String    `json:"description"`
	DescriptionSlate null.String    `json:"descriptionSlate"`
	SellingPrice     pgtype.Numeric `json:"sellingPrice"`
	CostPrice        pgtype.Numeric `json:"costPrice"`
	MediaID          null.String    `json:"mediaID"`
	Updated          null.Time      `json:"updated"`
	ID               string         `json:"id"`
	ShopID           string         `json:"shopID"`
}

func (q *Queries) UpdateProduct(ctx context.Context, arg UpdateProductParams) (Product, error) {
	row := q.db.QueryRow(ctx, updateProduct,
		arg.Name,
		arg.Description,
		arg.DescriptionSlate,
		arg.SellingPrice,
		arg.CostPrice,
		arg.MediaID,
		arg.Updated,
		arg.ID,
		arg.ShopID,
	)
	var i Product
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Description,
		&i.DescriptionSlate,
		&i.Msku,
		&i.SellingPrice,
		&i.CostPrice,
		&i.ShopID,
		&i.MediaID,
		&i.Created,
		&i.Updated,
	)
	return i, err
}

const updateProductsPlatformSync = `-- name: UpdateProductsPlatformSync :one
UPDATE products_platform
SET
 sync_lazada = coalesce($1, sync_lazada)
WHERE product_id = $2
RETURNING product_id, products_lazada_id, sync_lazada
`

type UpdateProductsPlatformSyncParams struct {
	SyncLazada null.Bool `json:"syncLazada"`
	ProductID  string    `json:"productID"`
}

func (q *Queries) UpdateProductsPlatformSync(ctx context.Context, arg UpdateProductsPlatformSyncParams) (ProductsPlatform, error) {
	row := q.db.QueryRow(ctx, updateProductsPlatformSync, arg.SyncLazada, arg.ProductID)
	var i ProductsPlatform
	err := row.Scan(&i.ProductID, &i.ProductsLazadaID, &i.SyncLazada)
	return i, err
}
