package main

import (
	"os"

	"github.com/joho/godotenv"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"laxo.vn/laxo/laxo"
	"laxo.vn/laxo/laxo/assets"
	lazada_service "laxo.vn/laxo/laxo/lazada"
	"laxo.vn/laxo/laxo/notification"
	"laxo.vn/laxo/laxo/shop"
	"laxo.vn/laxo/laxo/store"
	"laxo.vn/laxo/temporal/lazada"
)

func main() {
	logger := laxo.NewLogger()
	defer logger.Zap.Sync()

	logger.Info("Starting workers...")

	config, err := laxo.InitConfig()
	if err != nil {
		logger.Fatalw("Could not init config",
			"error", err,
		)
		return
	}

	server, err := laxo.NewServer(logger, config)
	if err != nil {
		logger.Fatalw("Failed to get server struct",
			"error", err,
		)
		return
	}

	if err = godotenv.Load(".env"); err != nil {
		logger.Fatal("Failed to load .env file")
	}

	c, err := client.NewClient(client.Options{
		HostPort: client.DefaultHostPort,
	})

	if err != nil {
		logger.Fatalw("Unable to create client",
			"error", err,
		)
	}
	defer c.Close()

	redisURI := os.Getenv("REDIS_URL")

	if err = server.InitRedis(redisURI); err != nil {
		logger.Errorw("Failed to init Redis",
			"error", err,
		)
		return
	}

	dbURI := os.Getenv("POSTGRESQL_URL")

	if err = server.InitDatabase(dbURI); err != nil {
		logger.Errorw("Failed to init Database",
			"error", err,
			"uri", dbURI,
		)
		return
	}

	assetsBasePath := os.Getenv("ASSETS_BASE_PATH")

	store, err := store.NewStore(dbURI, logger, assetsBasePath)
	if err != nil {
		logger.Fatalw("Failed to create new store",
			"error", err,
		)
		return
	}

	notificationService := notification.NewService(store, logger, server)
	shopService := shop.NewService(store, logger, server)
	assetsService := assets.NewService(store, logger, server)

	lazadaID := os.Getenv("LAZADA_ID")
	lazadaSecret := os.Getenv("LAZADA_SECRET")
	lazadaService := lazada_service.NewService(store, logger, server, lazadaID, lazadaSecret)

	w := worker.New(c, "product", worker.Options{})

	w.RegisterWorkflow(lazada.SyncLazadaPlatform)

	activities := lazada.NewActivities(&server.RedisClient, &notificationService, &lazadaService, &shopService, &assetsService)
	w.RegisterActivity(activities)

	err = w.Run(worker.InterruptCh())
	if err != nil {
		logger.Fatalw("Unable to start worker",
			"error", err,
		)
	}
}
