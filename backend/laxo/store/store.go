package store

import (
	"context"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/jackc/pgx/v4/pgxpool"
	"laxo.vn/laxo/laxo/sqlc"
)

type Store struct {
  lazadaStore
  notificationStore
  productStore
  assetsStore
  logger    hclog.Logger
  pglClient *pgxpool.Pool
  queries   *sqlc.Queries
}

func NewStore(uri string, logger hclog.Logger, assetsBasePath string) (*Store, error) {
  config, err := pgxpool.ParseConfig(uri)
  if err != nil {
    return nil, err
  }

	config.MaxConns = 10
	config.MinConns = 5
	config.HealthCheckPeriod = 20 * time.Second

  dbpool, err := pgxpool.ConnectConfig(context.Background(), config)
  if err != nil {
    return nil, err
  }

  queries := sqlc.New(dbpool)

  s := Store{
    logger:    logger,
    pglClient: dbpool,
    queries: queries,
  }

  s.lazadaStore = newLazadaStore(&s)
  s.notificationStore = newNotificationStore(&s)
  s.productStore = newProductStore(&s)

  a, err := newAssetsStore(&s, assetsBasePath)
  if err != nil {
    return nil, err
  }

  s.assetsStore = *a

  return &s, nil
}
