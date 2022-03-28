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
	laxo_proto "laxo.vn/laxo/laxo/proto"
	laxo_proto_gen "laxo.vn/laxo/laxo/proto/gen"
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
      Handler:      r,
      Addr:         "127.0.0.1:" + config.Port,
      WriteTimeout: 15 * time.Second,
      ReadTimeout:  15 * time.Second,
    }

    if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
      return err
    }

    return nil
  })


  var grpcHttpServer *http.Server

  logger.Info("Serving GRPC...", "port", "8081")

  g.Go(func() error {
    opts := []grpc.ServerOption{
      grpc.StreamInterceptor(grpc_auth.StreamServerInterceptor(laxo_proto.StreamAuthFunc)),
    }
    grpcServer := grpc.NewServer(opts...)
    protoServer := &laxo_proto.ProtoServer{}
    laxo_proto_gen.RegisterProductServiceServer(grpcServer, protoServer)

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

    if err := grpcHttpServer.ListenAndServe(); err != http.ErrServerClosed {
      return err
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

  _ = httpServer.Shutdown(shutdownCtx)
  _ = grpcHttpServer.Shutdown(shutdownCtx)

  err := g.Wait()
  if err != nil {
    logger.Error("Server returning an error", "error", err)
    os.Exit(2)
  }
 }
