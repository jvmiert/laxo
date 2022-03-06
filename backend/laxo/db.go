package laxo

import (
  "context"
  "time"

  //"github.com/sirupsen/logrus"
  //"github.com/jackc/pgx/v4/log/logrusadapter"
  "github.com/jackc/pgx/v4/pgxpool"
  //"github.com/jackc/pgx/v4"
  "laxo.vn/laxo/laxo/sqlc"
)

var PglClient *pgxpool.Pool
var Queries *sqlc.Queries

func InitDatabase(uri string) error {
  Logger.Debug("Connecting to Postgres", "uri", uri)

  config, err := pgxpool.ParseConfig(uri)
  if err != nil {
    return err
  }

	config.MaxConns = 10
	config.MinConns = 5
	config.HealthCheckPeriod = 20 * time.Second

//  config.ConnConfig.LogLevel = pgx.LogLevelTrace
//	config.ConnConfig.Logger = logrusadapter.NewLogger(logrus.New())

  dbpool, err := pgxpool.ConnectConfig(context.Background(), config)
  if err != nil {
    return err
  }

  PglClient = dbpool

  queries := sqlc.New(PglClient)
  Queries = queries

  return nil
}
