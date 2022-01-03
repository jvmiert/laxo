package main

import (
  "context"
  "time"
  "net/http"
  "os"
  "os/signal"

  "laxo.vn/laxo/laxo"
)

func main() {
  var config laxo.Config
  logger := laxo.InitConfig(&config)

  if err := laxo.InitRedis(); err != nil {
    logger.Error("Failed to init Redis", "error", err)
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
      logger.Error("Failed to listen and service", "error", err)
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
