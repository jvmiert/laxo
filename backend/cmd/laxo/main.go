package main

import (
  "context"
  "time"
  "net/http"
  "os"
  "os/signal"

  "laxo.vn/laxo/laxo"
  "github.com/joho/godotenv"
)

func main() {
  logger, config := laxo.InitConfig(false)

  if err := godotenv.Load(".env"); err != nil {
    logger.Error("Failed to load .env file")
  }

  redisURI := os.Getenv("REDIS_URL")

  if err := laxo.InitRedis(redisURI); err != nil {
    logger.Error("Failed to init Redis", "error", err)
    return
  }

  dbURI := os.Getenv("POSTGRESQL_URL")

  if err := laxo.InitDatabase(dbURI); err != nil {
    logger.Error("Failed to init Database", "uri", dbURI, "error", err)
    return
  }

  r := laxo.SetupRouter()

  logger.Info("Serving...", "port", config.Port)
  srv := &http.Server{
    Handler:      r,
    Addr:         "127.0.0.1:" + config.Port,
    WriteTimeout: 15 * time.Second,
    ReadTimeout:  15 * time.Second,
  }

  go func() {
    if err := srv.ListenAndServe(); err != nil {
      if err == http.ErrServerClosed {
        logger.Info("Closing server...")
      } else {
        logger.Error("Failed to listen and service", "error", err)
      }
    }
  }()

  c := make(chan os.Signal, 1)
  signal.Notify(c, os.Interrupt)

  <-c

  ctx, cancel := context.WithTimeout(context.Background(), 15 * time.Second)

  defer cancel()

  srv.Shutdown(ctx)

  logger.Info("Shutting down...")
  os.Exit(0)
 }
