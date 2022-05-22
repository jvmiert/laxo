package laxo

import (
	"context"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"laxo.vn/laxo/laxo/sqlc"
)

func (s *Server) InitDatabase(uri string) error {
  s.Logger.Infow("Connecting to Postgres",
    "uri", uri,
  )

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

  s.PglClient = dbpool
  s.Queries = sqlc.New(s.PglClient)

  return nil
}
