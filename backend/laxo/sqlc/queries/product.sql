-- name: CreateProduct :one
INSERT INTO products (
  name, description, msku, selling_price, cost_price, shop_id,
  media_id, updated
) VALUES (
  $1, $2, $3, $4, $5, $6, $7,
  $8
)
RETURNING *;

-- name: CreateProductMedia :one
INSERT INTO products_media (
  product_id, original_filename,
  murmur_hash
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: CreateProductPlatform :one
INSERT INTO products_platform (
  product_id, products_lazada_id
) VALUES (
  $1, $2
)
RETURNING *;

-- name: UpdateProduct :one
UPDATE products
SET
 name = coalesce(sqlc.narg('name'), name),
 description = coalesce(sqlc.narg('description'), description),
 msku = coalesce(sqlc.narg('msku'), msku),
 selling_price = coalesce(sqlc.narg('selling_price'), selling_price),
 cost_price = coalesce(sqlc.narg('cost_price'), cost_price),
 shop_id = coalesce(sqlc.narg('shop_id'), shop_id),
 media_id = coalesce(sqlc.narg('media_id'), media_id),
 updated = coalesce(sqlc.narg('updated'), updated)
WHERE id = sqlc.arg('id')
RETURNING *;

-- name: UpdateProductMedia :one
UPDATE products_media
SET
 product_id = coalesce(sqlc.narg('product_id'), product_id),
 original_filename = coalesce(sqlc.narg('original_filename'), original_filename),
 murmur_hash = coalesce(sqlc.narg('murmur_hash'), murmur_hash)
WHERE id = sqlc.arg('id')
RETURNING *;

-- name: GetProductPlatformByProductID :one
SELECT * FROM products_platform
WHERE product_id = $1
LIMIT 1;

-- name: GetProductPlatformByLazadaID :one
SELECT * FROM products_platform
WHERE products_lazada_id = $1
LIMIT 1;

-- name: GetProductMediaByID :one
SELECT * FROM products_media
WHERE id = $1
LIMIT 1;

-- name: GetProductMediaByProductID :one
SELECT * FROM products_media
WHERE product_id = $1
LIMIT 1;

-- name: GetProductMediaByMurmur :one
SELECT * FROM products_media
WHERE murmur_hash = $1
LIMIT 1;

-- name: GetProductByID :one
SELECT * FROM products
WHERE id = $1
LIMIT 1;

-- name: GetProductByProductMSKU :one
SELECT * FROM products
WHERE msku = $1 AND shop_id = $2
LIMIT 1;
