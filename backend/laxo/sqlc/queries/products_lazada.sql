-- name: CreateLazadaProduct :one
INSERT INTO products_lazada (
  lazada_id, lazada_primary_category, created, updated,
  status, sub_status, shop_id
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
)
RETURNING *;
