package store

import (
	"context"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"laxo.vn/laxo/laxo"
	"laxo.vn/laxo/laxo/sqlc"
)

type Store struct {
  lazadaStore
  notificationStore
  shopStore
  assetsStore
  userStore
  logger    *laxo.Logger
  pglClient *pgxpool.Pool
  queries   *sqlc.Queries
}

func NewStore(uri string, logger *laxo.Logger, assetsBasePath string) (*Store, error) {
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
  s.shopStore = newShopStore(&s)
  s.userStore = newUserStore(&s)

  a, err := newAssetsStore(&s, assetsBasePath)
  if err != nil {
    return nil, err
  }

  s.assetsStore = *a

  return &s, nil
}
