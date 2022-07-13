-- name: GetProductsByNameOrSKU :many
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
  SELECT products.*,
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
ON true;

-- name: GetProductsByShopID :many
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
  SELECT products.*,
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
ON true;

-- name: GetProductDetailsByID :one
SELECT products.*,
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
WHERE products.id = $1 AND products.shop_id = $2;

-- name: GetProductAssetsByProductID :many
SELECT assets.*, products_media.image_order, products_media.status FROM products
LEFT JOIN products_media ON products_media.product_id = products.id
LEFT JOIN assets ON assets.id = products_media.asset_id
WHERE products.id = $1 AND products.shop_id = $2
ORDER BY products_media.image_order, products_media.product_id, assets.id;

-- name: CreateProduct :one
INSERT INTO products (
  name, description, msku, selling_price, cost_price, shop_id,
  media_id, updated
) VALUES (
  $1, $2, $3, $4, $5, $6, $7,
  $8
)
RETURNING *;

-- name: CheckProductOwner :one
SELECT products.id FROM products
WHERE products.id = $1 AND products.shop_id = $2;

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
 description_slate = coalesce(sqlc.narg('description_slate'), description),
 selling_price = coalesce(sqlc.narg('selling_price'), selling_price),
 cost_price = coalesce(sqlc.narg('cost_price'), cost_price),
 media_id = coalesce(sqlc.narg('media_id'), media_id),
 updated = coalesce(sqlc.narg('updated'), updated)
WHERE id = sqlc.arg('id') AND shop_id = sqlc.arg('shop_id')
RETURNING *;


-- name: GetProductPlatformByProductID :one
SELECT * FROM products_platform
WHERE product_id = $1
LIMIT 1;

-- name: GetProductPlatformByLazadaID :one
SELECT * FROM products_platform
WHERE products_lazada_id = $1
LIMIT 1;

-- name: GetProductByID :one
SELECT * FROM products
WHERE id = $1
LIMIT 1;

-- name: GetProductByProductMSKU :one
SELECT * FROM products
WHERE msku = $1 AND shop_id = $2
LIMIT 1;

-- name: UpdateProductsPlatformSync :one
UPDATE products_platform
SET
 sync_lazada = coalesce(sqlc.narg('sync_lazada'), sync_lazada)
WHERE product_id = sqlc.arg('product_id')
RETURNING *;
