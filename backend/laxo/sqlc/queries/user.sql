-- name: GetUser :one
SELECT * FROM users
WHERE user_id = $1
LIMIT 1;

-- name: CreateUser :one
INSERT INTO users (
  username, password, email
) VALUES (
  $1, $2, $3
)
RETURNING *;
