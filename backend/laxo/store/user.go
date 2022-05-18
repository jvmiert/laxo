package store

import (
	"context"

	"laxo.vn/laxo/laxo"
)

type userStore struct {
	*Store
}

func newUserStore(store *Store) userStore {
	return userStore{
		store,
	}
}

func (s *userStore) RetrieveShopsByUserID(userID string) ([]laxo.Shop, error) {
  shops, err := s.queries.GetShopsByUserID(
    context.Background(),
    userID,
  )
  if err != nil {
    return nil, err
  }

  var sReturn []laxo.Shop

  for _, s := range shops {
    sReturn = append(sReturn, laxo.Shop{Model: &s})
  }

  return sReturn, nil
}
