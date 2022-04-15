package processing

import (
	"context"

	"go.temporal.io/sdk/activity"
)

type Activities struct {
  RedisClient string
}

func (a *Activities) FetchLazadaProductsFromAPI(ctx context.Context, shopID string) (string, error) {
  logger := activity.GetLogger(ctx)
  logger.Info("We can access the Redis client fine", "RedisClient", a.RedisClient)
  logger.Info("Retrieving Lazada products", "shopID", shopID)
  return "dataKeyID", nil
}


func (a *Activities) SaveLazadaProducts(ctx context.Context, dataKeyID string) error {
  logger := activity.GetLogger(ctx)
  logger.Info("Save Lazada products", "dataKeyID", dataKeyID)
  return nil
}


func (a *Activities) ProcessLazadaProducts(ctx context.Context, shopID string) error {
  logger := activity.GetLogger(ctx)
  logger.Info("Processing Lazada products", "shopID", shopID)
  return nil
}
