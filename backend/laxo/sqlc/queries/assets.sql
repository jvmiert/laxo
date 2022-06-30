-- name: CreateProductMedia :one
INSERT INTO products_media (
  product_id, asset_id,
  image_order, status
) VALUES (
  $1, $2, $3, $4
)
RETURNING *;

-- name: GetProductMedia :one
SELECT * FROM products_media
WHERE product_id = $1 AND asset_id = $2;

-- name: UpdateProductMedia :one
UPDATE products_media
SET
 image_order = coalesce(sqlc.narg('image_order'), image_order),
 status = coalesce(sqlc.narg('status'), status)
WHERE product_id = sqlc.arg('product_id') AND asset_id = sqlc.arg('asset_id')
RETURNING *;

-- name: DeleteProductMedia :exec
DELETE FROM products_media
WHERE product_id = $1 AND asset_id = $2;

-- name: CreateAsset :one
INSERT INTO assets (
  shop_id, murmur_hash, original_filename, extension,
  file_size, width, height
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
)
RETURNING *;

-- name: UpdateAsset :one
UPDATE assets
SET
 shop_id = coalesce(sqlc.narg('shop_id'), shop_id),
 murmur_hash = coalesce(sqlc.narg('murmur_hash'), murmur_hash),
 original_filename = coalesce(sqlc.narg('original_filename'), original_filename),
 extension = coalesce(sqlc.narg('extension'), extension),
 file_size = coalesce(sqlc.narg('file_size'), file_size),
 width = coalesce(sqlc.narg('width'), width),
 height = coalesce(sqlc.narg('height'), height)
WHERE id = sqlc.arg('id')
RETURNING *;

-- name: GetAssetByMurmur :one
SELECT * FROM assets
WHERE murmur_hash = $1 AND shop_id = $2
LIMIT 1;

-- name: GetAssetByOriginalName :one
SELECT * FROM assets
WHERE original_filename = $1 AND shop_id = $2
LIMIT 1;

-- name: GetAssetByID :one
SELECT * FROM assets
WHERE id = $1
LIMIT 1;

-- name: GetAllAssetsByShopID :many
SELECT
  c.count, p.*
FROM
(
  SELECT COUNT(*) AS COUNT
  FROM assets
  WHERE assets.shop_id = $1
) as c
LEFT JOIN (
  SELECT assets.*
  FROM assets
  ORDER BY assets.created
  LIMIT $2 OFFSET $3
) as p
ON true;
