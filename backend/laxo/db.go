package laxo

import (
  "context"

  "github.com/jackc/pgx/v4/pgxpool"
  "laxo.vn/laxo/laxo/sqlc"
)

var PglClient *pgxpool.Pool
var Queries *sqlc.Queries

func InitDatabase(uri string) error {
  dbpool, err := pgxpool.Connect(context.Background(), uri)
  if err != nil {
    return err
  }

  PglClient = dbpool

  queries := sqlc.New(PglClient)
  Queries = queries

  return nil
}
