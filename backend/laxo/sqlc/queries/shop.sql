-- name: GetShopByID :one
SELECT * FROM shops
WHERE id = $1
LIMIT 1;

-- name: GetShopsByUserID :many
SELECT * FROM shops
WHERE user_id = $1;

-- name: CreateShop :one
INSERT INTO shops (
  shop_name, user_id
) VALUES (
  $1, $2
)
RETURNING *;
