package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"github.com/joho/godotenv"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"laxo.vn/laxo/laxo"
	"laxo.vn/laxo/laxo/assets"
	"laxo.vn/laxo/laxo/http/rest"
	"laxo.vn/laxo/laxo/lazada"
	"laxo.vn/laxo/laxo/notification"
	laxo_proto "laxo.vn/laxo/laxo/proto"
	laxo_proto_gen "laxo.vn/laxo/laxo/proto/gen"
	"laxo.vn/laxo/laxo/shop"
	"laxo.vn/laxo/laxo/store"
	"laxo.vn/laxo/laxo/user"
	temporal_client "laxo.vn/laxo/temporal/client"
)

func main() {
  logger := laxo.NewLogger()
  defer logger.Zap.Sync()

  config, err := laxo.InitConfig()
  if err != nil {
    logger.Errorw("Could not init config",
      "error", err,
    )
  }

  if err = godotenv.Load(".env"); err != nil {
    logger.Errorw("Failed to load .env file",
      "error", err,
    )
  }

  server, err := laxo.NewServer(logger, config)
  if err != nil {
    logger.Errorw("Failed to get server struct",
      "error", err,
    )
  }

  server.InitMiddleware()

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

  temporalClient, err := temporal_client.NewClient(logger)
  if err != nil {
    logger.Errorw("Failed to create Temporal client",
      "error", err,
    )
  }

  assetsBasePath := os.Getenv("ASSETS_BASE_PATH")

  store, err := store.NewStore(dbURI, logger, assetsBasePath)
  if err != nil {
    logger.Error("Failed to create new store", "error", err)
    return
  }

  notificationService := notification.NewService(store, logger, server)
  rest.InitNotificationHandler(server, &notificationService, server.Router, server.Negroni)

  lazadaID := os.Getenv("LAZADA_ID")
  lazadaSecret := os.Getenv("LAZADA_SECRET")
  lazadaService := lazada.NewService(store, logger, server, lazadaID, lazadaSecret)

  shopService := shop.NewService(store, logger, server)
  rest.InitShopHandler(server, &shopService, &lazadaService, server.Router,
                       server.Negroni, temporalClient)

  userService := user.NewService(store, logger, server)
  rest.InitUserHandler(server, &userService, server.Router, server.Negroni)

  assetsService := assets.NewService(store, logger, server)

  rest.InitTestHandler(server, &lazadaService, &shopService, &assetsService, server.Router, server.Negroni)

  ctx := context.Background()
  ctx, cancel := context.WithCancel(ctx)
  defer cancel()

  interrupt := make(chan os.Signal, 1)
  signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
  defer signal.Stop(interrupt)

  g, ctx := errgroup.WithContext(ctx)

  var httpServer *http.Server

  logger.Info("Serving http...", "port", config.Port)

  g.Go(func() error {
    httpServer = &http.Server{
      Handler:      server.Router,
      Addr:         "127.0.0.1:" + config.Port,
      WriteTimeout: 15 * time.Second,
      ReadTimeout:  15 * time.Second,
    }

    if errServer := httpServer.ListenAndServe(); err != http.ErrServerClosed {
      return errServer
    }

    return nil
  })

  var grpcHttpServer *http.Server

  logger.Info("Serving GRPC...", "port", "8081")

  g.Go(func() error {
    protoServer, errGRPC := laxo_proto.NewServer(
      &notificationService,
      logger,
      redisURI,
      ctx,
      server,
    )
    if errGRPC != nil {
      logger.Error("GRPC Redis error", "error", errGRPC)
      return errGRPC
    }

    protoMiddleware := laxo_proto.NewGRPCMiddleware(server)

    opts := []grpc.ServerOption{
      grpc.StreamInterceptor(grpc_auth.StreamServerInterceptor(protoMiddleware.StreamAuthFunc)),
    }
    grpcServer := grpc.NewServer(opts...)

    laxo_proto_gen.RegisterUserServiceServer(grpcServer, protoServer)

    option := []grpcweb.Option{
      grpcweb.WithWebsockets(true),
      grpcweb.WithOriginFunc(func(origin string) bool {
        // Allow all origins, DO NOT do this in production
        return true
      }),
    }
    wrappedServer := grpcweb.WrapServer(
      grpcServer,
      option...,
    )

    handler := func(resp http.ResponseWriter, req *http.Request) {
      wrappedServer.ServeHTTP(resp, req)
    }

    grpcHttpServer = &http.Server{
      Addr:    "127.0.0.1:8081",
      Handler: http.HandlerFunc(handler),
    }

    if errGRPC = grpcHttpServer.ListenAndServe(); err != http.ErrServerClosed {
      return errGRPC
    }

    return nil
  })

  select {
  case <-interrupt:
    break
  case <-ctx.Done():
    break
  }

  logger.Info("Received shutdown signal...")

  cancel()

  shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 15*time.Second)
  defer shutdownCancel()

  if temporalClient != nil {
    temporalClient.Close()
  }

  _ = httpServer.Shutdown(shutdownCtx)
  _ = grpcHttpServer.Shutdown(shutdownCtx)

  err = g.Wait()
  if err != nil && err != http.ErrServerClosed {
    logger.Error("Server returning an error", "error", err)
    os.Exit(2)
  }
 }
