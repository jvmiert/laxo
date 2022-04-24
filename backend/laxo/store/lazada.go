package store

import (
	"context"

	"github.com/jackc/pgx/v4"
	"laxo.vn/laxo/laxo/lazada"
	"laxo.vn/laxo/laxo/sqlc"
)

type lazadaStore struct {
  *Store
}

func newLazadaStore(store *Store) lazadaStore {
  return lazadaStore{
    store,
  }
}

func (s *lazadaStore) SaveOrUpdateProduct(p lazada.ProductsResponseProducts, shopID string) error {
  qParam := sqlc.GetLazadaProductByLazadaIDParams{
    LazadaID: p.ItemID,
    ShopID: shopID,
  }

  pModel, err := s.queries.GetLazadaProductByLazadaID(context.Background(), qParam)

  if err != pgx.ErrNoRows && err != nil {
    return err
  }

  if pModel.ID == "" {
    params := sqlc.CreateLazadaProductParams{
      LazadaID: p.ItemID,
      LazadaPrimaryCategory: p.PrimaryCategory,
      Created: p.Created,
      Updated: p.Updated,
      Status: p.Status.NullString,
      SubStatus: p.SubStatus.NullString,
      ShopID: shopID,
    }
    _, err = s.queries.CreateLazadaProduct(
      context.Background(),
      params,
    )

    if err != nil {
      return err
    }

    return nil
  }

  params := sqlc.UpdateLazadaProductParams {
    LazadaID: p.ItemID,
    LazadaPrimaryCategory: p.PrimaryCategory,
    Created: p.Created,
    Updated: p.Updated,
    Status: p.Status.NullString,
    SubStatus: p.SubStatus.NullString,
    ID: pModel.ID,
  }

  err = s.queries.UpdateLazadaProduct(
    context.Background(),
    params,
  )

  if err != nil {
    return err
  }

  return nil
}
