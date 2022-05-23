-- name: GetPlatformsByShopID :many
SELECT * FROM shops_platforms
WHERE shop_id = $1
ORDER BY shops_platforms.platform_name;

-- name: GetSpecificPlatformByShopID :one
SELECT * FROM shops_platforms
WHERE shop_id = $1 AND platform_name = $2
LIMIT 1;

-- name: CreatePlatform :one
INSERT INTO shops_platforms (
  shop_id, platform_name
) VALUES (
  $1, $2
)
RETURNING *;

