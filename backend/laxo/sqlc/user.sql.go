// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.13.0
// source: user.sql

package sqlc

import (
	"context"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users (
  email, password, fullname
) VALUES (
  $1, $2, $3
)
RETURNING id, password, email, created, fullname
`

type CreateUserParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Fullname string `json:"fullname"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRow(ctx, createUser, arg.Email, arg.Password, arg.Fullname)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Password,
		&i.Email,
		&i.Created,
		&i.Fullname,
	)
	return i, err
}

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT id, password, email, created, fullname FROM users
WHERE email = $1
LIMIT 1
`

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRow(ctx, getUserByEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Password,
		&i.Email,
		&i.Created,
		&i.Fullname,
	)
	return i, err
}

const getUserByID = `-- name: GetUserByID :one
SELECT id, password, email, created, fullname FROM users
WHERE id = $1
LIMIT 1
`

func (q *Queries) GetUserByID(ctx context.Context, id string) (User, error) {
	row := q.db.QueryRow(ctx, getUserByID, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Password,
		&i.Email,
		&i.Created,
		&i.Fullname,
	)
	return i, err
}
