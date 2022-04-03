-- name: GetShopByID :one
SELECT * FROM shops
WHERE id = $1
LIMIT 1;

-- name: GetShopsByUserID :many
SELECT * FROM shops
WHERE user_id = $1;

-- name: GetShopsPlatformsByUserID :many
SELECT shops.id,
       shops.user_id,
       shops.shop_name,
       shops.created,
       shops_platforms.id as platform_id,
       shops_platforms.platform_name as platform_name,
       shops_platforms.created as platform_created
FROM shops
LEFT JOIN shops_platforms ON (shops.id = shops_platforms.shop_id)
WHERE user_id = $1;

-- name: CreateShop :one
INSERT INTO shops (
  shop_name, user_id
) VALUES (
  $1, $2
)
RETURNING *;
