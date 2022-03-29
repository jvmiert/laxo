-- name: CreateShop :one
INSERT INTO shops (
  shop_name, user_id
) VALUES (
  $1, $2
)
RETURNING *;
