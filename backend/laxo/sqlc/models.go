// Code generated by sqlc. DO NOT EDIT.

package sqlc

import (
	"database/sql"
)

type Shop struct {
	ID         string       `json:"id"`
	UserID     string       `json:"userID"`
	ShopName   string       `json:"shopName"`
	Created    sql.NullTime `json:"created"`
	LastUpdate sql.NullTime `json:"lastUpdate"`
}

type User struct {
	ID       string       `json:"id"`
	Password string       `json:"password"`
	Email    string       `json:"email"`
	Created  sql.NullTime `json:"created"`
	Fullname string       `json:"fullname"`
}
