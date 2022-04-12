-- name: CreateLazadaPlatform :one
INSERT INTO platform_lazada (
  shop_id, access_token, country, refresh_token, account_platform, account,
  user_id_vn, seller_id_vn, short_code_vn, refresh_expires_in,
  access_expires_in
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
)
RETURNING *;

-- name: GetLazadaPlatformByShopID :one
SELECT * FROM platform_lazada
WHERE shop_id = $1
LIMIT 1;

-- name: UpdateLazadaPlatform :exec
UPDATE platform_lazada SET shop_id = $1, access_token = $2, country = $3,
refresh_token = $4, account_platform = $5, account = $6, user_id_vn = $7,
seller_id_vn = $8, short_code_vn = $9, refresh_expires_in = $10,
access_expires_in = $11
WHERE id = $12;
