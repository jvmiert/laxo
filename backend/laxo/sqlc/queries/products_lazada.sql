-- name: GetLazadaProductByLazadaID :one
SELECT * FROM products_lazada
WHERE lazada_id = $1 AND shop_id = $2
LIMIT 1;

-- name: UpdateLazadaProduct :exec
UPDATE products_lazada SET
  lazada_id = $1, lazada_primary_category = $2, created = $3, updated = $4,
  status = $5, sub_status = $6
WHERE id = $7;

-- name: CreateLazadaProduct :one
INSERT INTO products_lazada (
  lazada_id, lazada_primary_category, created, updated,
  status, sub_status, shop_id
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
)
RETURNING *;
