-- name: GetPlatformsByShopID :many
SELECT * FROM shops_platforms
WHERE shop_id = $1
ORDER BY shops_platforms.platform_name;


