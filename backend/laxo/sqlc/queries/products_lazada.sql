-- name: GetLazadaProductAttributeByProductID :one
SELECT * FROM products_attribute_lazada
WHERE product_id = $1
LIMIT 1;

-- name: UpdateLazadaProductAttribute :one
UPDATE products_attribute_lazada SET
  name = $1, short_description = $2, description = $3, brand = $4, model = $5,
  headphone_features = $6, bluetooth = $7, warranty_type = $8, warranty = $9,
  hazmat = $10, expire_date = $11, brand_classification = $12,
  ingredient_preference = $13, lot_number = $14, units_hb = $15,
  fmlt_skincare = $16, quantitative = $17, skincare_by_age = $18,
  skin_benefit = $19, skin_type = $20, user_manual = $21,
  country_origin_hb = $22, color_family = $23, fragrance_family = $24,
  source = $25
WHERE id = $26
RETURNING *;

-- name: CreateLazadaProductAttribute :one
INSERT INTO products_attribute_lazada (
  name, short_description, description, brand, model, headphone_features,
  bluetooth, warranty_type, warranty, hazmat, expire_date,
  brand_classification, ingredient_preference, lot_number, units_hb,
  fmlt_skincare, quantitative, skincare_by_age, skin_benefit, skin_type,
  user_manual, country_origin_hb, color_family, fragrance_family,
  source, product_id
) VALUES (
  $1, $2, $3, $4, $5, $6, $7,
  $8, $9, $10, $11, $12, $13, $14,
  $15, $16, $17, $18, $19, $20, $21,
  $22, $23, $24, $25, $26
)
RETURNING *;

-- name: GetLazadaProductSKUByProductID :one
SELECT * FROM products_sku_lazada
WHERE product_id = $1
LIMIT 1;

-- name: CreateLazadaProductSKU :one
INSERT INTO products_sku_lazada (
  status, quantity, seller_sku, shop_sku, sku_id, url,
  price, available, package_content, package_width, package_weight,
  package_length, package_height, special_price, special_to_time,
  special_from_time, special_from_date, special_to_date, product_id, shop_id
) VALUES (
  $1, $2, $3, $4, $5, $6, $7,
  $8, $9, $10, $11, $12, $13, $14,
  $15, $16, $17, $18, $19, $20
)
RETURNING *;

-- name: UpdateLazadaProductSKU :one
UPDATE products_sku_lazada SET
  status = $1, quantity = $2, seller_sku = $3, shop_sku = $4, sku_id = $5, url = $6,
  price = $7, available = $8, package_content = $9, package_width = $10, package_weight = $11,
  package_length = $12, package_height = $13, special_price = $14, special_to_time = $15,
  special_from_time = $16, special_from_date = $17, special_to_date = $18, product_id = $19, shop_id = $20
WHERE id = $21
RETURNING *;

-- name: GetLazadaProductByLazadaID :one
SELECT * FROM products_lazada
WHERE lazada_id = $1 AND shop_id = $2
LIMIT 1;

-- name: UpdateLazadaProduct :one
UPDATE products_lazada SET
  lazada_id = $1, lazada_primary_category = $2, created = $3, updated = $4,
  status = $5, sub_status = $6
WHERE id = $7
RETURNING *;

-- name: CreateLazadaProduct :one
INSERT INTO products_lazada (
  lazada_id, lazada_primary_category, created, updated,
  status, sub_status, shop_id
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
)
RETURNING *;
