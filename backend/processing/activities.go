package processing

import (
	"context"
	"time"

	"go.temporal.io/sdk/activity"
)

type LazadaFetchResult struct {
  DataKey         string
  TotalProducts   int
}

type LazadaSaveParam struct {
  DataKey         string
  ProductIndex    int
}

type Activities struct {
  RedisClient string
}

func (a *Activities) FetchLazadaProductsFromAPI(ctx context.Context, shopID string) (LazadaFetchResult, error) {
  logger := activity.GetLogger(ctx)
  logger.Info("We can access the Redis client fine", "RedisClient", a.RedisClient)
  logger.Info("Retrieving Lazada products", "shopID", shopID)
  time.Sleep(5 * time.Second)
  return LazadaFetchResult{DataKey: "dataKeyID", TotalProducts: 20}, nil
}


func (a *Activities) SaveLazadaProducts(ctx context.Context, param LazadaSaveParam) error {
  logger := activity.GetLogger(ctx)
  logger.Info("Save Lazada products", "dataKeyID", param.DataKey, "index", param.ProductIndex)

  time.Sleep(5 * time.Second)

  return nil
}


func (a *Activities) ProcessLazadaProducts(ctx context.Context, shopID string) error {
  logger := activity.GetLogger(ctx)
  logger.Info("Processing Lazada products", "shopID", shopID)
  time.Sleep(5 * time.Second)
  return nil
}
