-- name: CreateShop :one
INSERT INTO shops (
  shop_name
) VALUES (
  $1
)
RETURNING *;
