package processing

import (
	"context"
	"strconv"
	"time"

	"github.com/mediocregopher/radix/v4"
	"go.temporal.io/sdk/activity"
	"laxo.vn/laxo/laxo/notification"
)

type LazadaFetchResult struct {
  DataKey         string
  TotalProducts   int
}

type LazadaSaveParam struct {
  DataKey         string
  ProductIndex    int
  ProductTotal    int
}

type Activities struct {
  RedisClient radix.Client
  NotificationService notification.Service
}

func (a *Activities) FetchLazadaProductsFromAPI(ctx context.Context, shopID string) (LazadaFetchResult, error) {
  logger := activity.GetLogger(ctx)
  info := activity.GetInfo(ctx)
  workFlowID := info.WorkflowExecution.ID

  logger.Info("Retrieving Lazada products", "shopID", shopID)
  time.Sleep(5 * time.Second)

  var StreamID radix.StreamEntryID
  a.RedisClient.Do(context.Background(), radix.Cmd(&StreamID, "XADD", workFlowID, "*", "state", "save", "complete", "-1", "total", "5"))
  return LazadaFetchResult{DataKey: "dataKeyID", TotalProducts: 5}, nil
}


func (a *Activities) SaveLazadaProducts(ctx context.Context, param LazadaSaveParam) error {
  logger := activity.GetLogger(ctx)
  info := activity.GetInfo(ctx)
  workFlowID := info.WorkflowExecution.ID

  logger.Info("Save Lazada products", "dataKeyID", param.DataKey, "index", param.ProductIndex)

  time.Sleep(5 * time.Second)


  strIndex := strconv.Itoa(param.ProductIndex + 1)
  strTotal := strconv.Itoa(param.ProductTotal)

  var StreamID radix.StreamEntryID
  a.RedisClient.Do(context.Background(), radix.Cmd(&StreamID, "XADD", workFlowID, "*", "state", "save", "complete", strIndex, "total", strTotal))
  return nil
}


func (a *Activities) ProcessLazadaProducts(ctx context.Context, shopID string) error {
  logger := activity.GetLogger(ctx)
  info := activity.GetInfo(ctx)
  workFlowID := info.WorkflowExecution.ID

  logger.Info("Processing Lazada products", "shopID", shopID)

  var StreamID radix.StreamEntryID
  a.RedisClient.Do(context.Background(), radix.Cmd(&StreamID, "XADD", workFlowID, "*", "state", "process"))
  time.Sleep(5 * time.Second)
  a.RedisClient.Do(context.Background(), radix.Cmd(&StreamID, "XADD", workFlowID, "*", "state", "complete"))

  // expire the key in a day
  a.RedisClient.Do(context.Background(), radix.Cmd("EXPIRE", workFlowID, "86400"))
  return nil
}
