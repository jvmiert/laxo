package store

import (
	"context"
	"strings"

	"laxo.vn/laxo/laxo/shop"
	"laxo.vn/laxo/laxo/sqlc"
	"laxo.vn/laxo/laxo/user"
)

type userStore struct {
	*Store
}

func newUserStore(store *Store) userStore {
	return userStore{
		store,
	}
}

func (s *userStore) RetrieveShopsByUserID(userID string) ([]shop.Shop, error) {
	shops, err := s.queries.GetShopsByUserID(
		context.Background(),
		userID,
	)
	if err != nil {
		return nil, err
	}

	var sReturn []shop.Shop

	for _, s := range shops {
		sReturn = append(sReturn, shop.Shop{Model: &s})
	}

	return sReturn, nil
}

func (s *userStore) RetrieveUserFromDBbyEmail(email string) (*user.User, error) {
	lowerEmail := strings.ToLower(email)
	uModel, err := s.queries.GetUserByEmail(
		context.Background(),
		strings.TrimSpace(lowerEmail),
	)

	var u user.User
	u.Model = &uModel

	return &u, err
}

func (s *userStore) RetrieveUserFromDBbyID(uID string) (*user.User, error) {
	uModel, err := s.queries.GetUserByID(
		context.Background(),
		uID,
	)

	var u user.User
	u.Model = &uModel

	return &u, err
}

func (s *userStore) SaveNewUserToDB(param sqlc.CreateUserParams) (*user.User, error) {
	savedUser, err := s.queries.CreateUser(
		context.Background(),
		param,
	)

	var u user.User
	u.Model = &savedUser

	return &u, err
}
