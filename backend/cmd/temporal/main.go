package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/mediocregopher/radix/v4"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"laxo.vn/laxo/laxo"
	"laxo.vn/laxo/laxo/notification"
	"laxo.vn/laxo/laxo/store"
	"laxo.vn/laxo/processing"
)

func main() {
  log.Println("Starting workers...")

  if err := godotenv.Load(".env"); err != nil {
		log.Fatalln("Failed to load .env file")
  }

  c, err := client.NewClient(client.Options{
		HostPort: client.DefaultHostPort,
	})

	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

  redisURI := os.Getenv("REDIS_URL")

  if redisURI == "" {
		log.Fatalln("Redis URL not set in env")
  }

  client, err := (radix.PoolConfig{}).New(context.Background(), "tcp", redisURI)

  if err != nil {
		log.Fatalln("Unable to connect to Redis", err)
  }

  logger, _ := laxo.InitConfig(false)

  dbURI := os.Getenv("POSTGRESQL_URL")
  store, err := store.NewStore(dbURI, logger)

  if err != nil {
    log.Fatalln("Failed to create new store",  err)
    return
  }

  notificationService := notification.NewService(store, logger, client)

	w := worker.New(c, "product", worker.Options{})

	w.RegisterWorkflow(processing.ProcessLazadaProducts)

  activities := &processing.Activities{RedisClient: client, NotificationService: notificationService}
	w.RegisterActivity(activities)

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}
}
