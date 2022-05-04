package processing

import (
	"context"
	"strconv"
	"time"

	"github.com/mediocregopher/radix/v4"
	"go.temporal.io/sdk/activity"
	"gopkg.in/guregu/null.v4"
	"laxo.vn/laxo/laxo/notification"
)

const (
  ActivityStateFetch = "fetch"
  ActivityStateSave = "save"
)

type LazadaFetchResult struct {
  DataKey         string
  TotalProducts   int
}

type LazadaSaveParam struct {
  UserID              string
  DataKey             string
  ProductIndex        int
  ProductTotal        int
}

type Activities struct {
  RedisClient radix.Client
  NotificationService notification.Service
}

func (a *Activities) FetchLazadaProductsFromAPI(ctx context.Context, shopID string, userID string) (LazadaFetchResult, error) {
  logger := activity.GetLogger(ctx)
  info := activity.GetInfo(ctx)
  workflowID := info.WorkflowExecution.ID

  notifyGroupParam := notification.NotificationGroupCreateParam{
    WorkflowID: null.NewString(workflowID, true),
    UserID: userID,
    EntityID: shopID,
    EntityType: notification.EntityTypeProductAdd,
    TotalMainSteps: null.NewInt(2, true),
    TotalSubSteps: null.NewInt(0, false),
  }

  notificationGroupID, err := a.NotificationService.CreateNotificationGroup(notifyGroupParam)
  if err != nil {
    return LazadaFetchResult{}, err
  }

  notifyParam := notification.NotificationCreateParam{
    GroupID: notificationGroupID,
    CurrentMainStep: null.NewInt(1, true),
    CurrentSubStep: null.NewInt(0, false),
    MainMessage: null.NewString(ActivityStateFetch, true),
    SubMessage: null.NewString("", false),
    ReadTime: null.NewTime(time.Time{}, false),
  }

  err = a.NotificationService.CreateNotification(notifyParam)
  if err != nil {
    return LazadaFetchResult{}, err
  }

  logger.Info("Retrieving Lazada products", "shopID", shopID)
  time.Sleep(5 * time.Second)

  var StreamID radix.StreamEntryID
  a.RedisClient.Do(context.Background(), radix.Cmd(&StreamID, "XADD", workflowID, "*", "state", "save", "complete", "-1", "total", "5"))

  updateParam := notification.NotificationGroupUpdateParam{
    TotalSubSteps: null.NewInt(5, true),
    ID: notificationGroupID,
  }

  err = a.NotificationService.UpdateNotificationGroup(updateParam)
  if err != nil {
    return LazadaFetchResult{}, err
  }

  notifyParam.MainMessage = null.NewString(ActivityStateSave, true)
  notifyParam.CurrentMainStep = null.NewInt(2, true)
  notifyParam.CurrentSubStep = null.NewInt(0, true)

  err = a.NotificationService.CreateNotification(notifyParam)
  if err != nil {
    return LazadaFetchResult{}, err
  }

  return LazadaFetchResult{DataKey: "dataKeyID", TotalProducts: 5}, nil
}

func (a *Activities) SaveLazadaProducts(ctx context.Context, param LazadaSaveParam) error {
  logger := activity.GetLogger(ctx)
  info := activity.GetInfo(ctx)
  workflowID := info.WorkflowExecution.ID

  logger.Info("Save Lazada products", "dataKeyID", param.DataKey, "index", param.ProductIndex)

  time.Sleep(5 * time.Second)

  strIndex := strconv.Itoa(param.ProductIndex + 1)
  strTotal := strconv.Itoa(param.ProductTotal)

  var StreamID radix.StreamEntryID
  a.RedisClient.Do(context.Background(), radix.Cmd(&StreamID, "XADD", workflowID, "*", "state", "save", "complete", strIndex, "total", strTotal))

  notificationGroupID, err := a.NotificationService.GetNotificationGroupIDByWorkflowID(workflowID, param.UserID)
  if err != nil {
    return err
  }

  notifyParam := notification.NotificationCreateParam{
    GroupID: notificationGroupID,
    CurrentMainStep: null.NewInt(2, true),
    CurrentSubStep: null.NewInt(int64(param.ProductIndex + 1), true),
    MainMessage: null.NewString(ActivityStateSave, true),
    SubMessage: null.NewString("", false),
    ReadTime: null.NewTime(time.Time{}, false),
  }

  err = a.NotificationService.CreateNotification(notifyParam)
  if err != nil {
    return err
  }

  return nil
}

func (a *Activities) ProcessLazadaProducts(ctx context.Context, userID string) error {
  logger := activity.GetLogger(ctx)
  info := activity.GetInfo(ctx)
  workflowID := info.WorkflowExecution.ID

  logger.Info("Processing Lazada products", "userID", userID)

  var StreamID radix.StreamEntryID
  a.RedisClient.Do(context.Background(), radix.Cmd(&StreamID, "XADD", workflowID, "*", "state", "process"))
  time.Sleep(5 * time.Second)
  a.RedisClient.Do(context.Background(), radix.Cmd(&StreamID, "XADD", workflowID, "*", "state", "complete"))

  // expire the key in a day
  a.RedisClient.Do(context.Background(), radix.Cmd("EXPIRE", workflowID, "86400"))
  return nil
}
