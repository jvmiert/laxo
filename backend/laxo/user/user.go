package user

import (
	"encoding/json"

	"gopkg.in/guregu/null.v4"
	"laxo.vn/laxo/laxo/sqlc"
)

type UserReturn struct {
	ID       string    `json:"id"`
	Email    string    `json:"email"`
	Created  null.Time `json:"created"`
	Fullname string    `json:"fullname"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type User struct {
	Model      *sqlc.User
	SessionKey string
}

func (u *User) JSON() ([]byte, error) {
	var ur UserReturn
	ur.ID = u.Model.ID
	ur.Email = u.Model.Email
	ur.Created = u.Model.Created
	ur.Fullname = u.Model.Fullname

	bytes, err := json.Marshal(ur)

	if err != nil {
		return bytes, err
	}

	return bytes, nil
}
