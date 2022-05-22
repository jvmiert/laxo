package lazada

import (
	"context"
	"time"

	"github.com/mediocregopher/radix/v4"
	"go.temporal.io/sdk/activity"
	"gopkg.in/guregu/null.v4"
	"laxo.vn/laxo/laxo/notification"
)

const (
  ActivityStateFetch = "fetch"
  ActivityStateSave = "save"
  ActivityStateComplete = "complete"
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
    WorkflowID: null.StringFrom(workflowID),
    UserID: userID,
    EntityID: shopID,
    EntityType: notification.EntityTypeProductAdd,
    TotalMainSteps: null.IntFrom(2),
    TotalSubSteps: null.NewInt(0, false),
  }

  notificationGroupID, err := a.NotificationService.CreateNotificationGroup(notifyGroupParam)
  if err != nil {
    return LazadaFetchResult{}, err
  }

  notifyParam := notification.NotificationCreateParam{
    GroupID: notificationGroupID,
    CurrentMainStep: null.IntFrom(1),
    CurrentSubStep: null.NewInt(0, false),
    MainMessage: null.StringFrom(ActivityStateFetch),
    SubMessage: null.NewString("", false),
    ReadTime: null.NewTime(time.Time{}, false),
  }

  err = a.NotificationService.CreateNotification(notifyParam)
  if err != nil {
    return LazadaFetchResult{}, err
  }

  logger.Info("Retrieving Lazada products", "shopID", shopID)
  time.Sleep(5 * time.Second)

  updateParam := notification.NotificationGroupUpdateParam{
    TotalSubSteps: null.IntFrom(5),
    ID: notificationGroupID,
  }

  err = a.NotificationService.UpdateNotificationGroup(updateParam)
  if err != nil {
    return LazadaFetchResult{}, err
  }

  notifyParam.MainMessage = null.StringFrom(ActivityStateSave)
  notifyParam.CurrentMainStep = null.IntFrom(2)
  notifyParam.CurrentSubStep = null.IntFrom(0)

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

  notificationGroupID, err := a.NotificationService.GetNotificationGroupIDByWorkflowID(workflowID, param.UserID)
  if err != nil {
    return err
  }

  notifyParam := notification.NotificationCreateParam{
    GroupID: notificationGroupID,
    CurrentMainStep: null.IntFrom(2),
    CurrentSubStep: null.IntFrom(int64(param.ProductIndex + 1)),
    MainMessage: null.StringFrom(ActivityStateSave),
    SubMessage: null.NewString("", false),
    ReadTime: null.NewTime(time.Time{}, false),
  }

  err = a.NotificationService.CreateNotification(notifyParam)
  if err != nil {
    return err
  }

  return nil
}

func (a *Activities) CompleteLazadaProducts(ctx context.Context, userID string) error {
  logger := activity.GetLogger(ctx)
  info := activity.GetInfo(ctx)
  workflowID := info.WorkflowExecution.ID

  logger.Info("Complete Lazada products")

  notificationGroupID, err := a.NotificationService.GetNotificationGroupIDByWorkflowID(workflowID, userID)
  if err != nil {
    return err
  }

  notifyParam := notification.NotificationCreateParam{
    GroupID: notificationGroupID,
    CurrentMainStep: null.IntFrom(2),
    MainMessage: null.StringFrom(ActivityStateComplete),
  }

  err = a.NotificationService.CreateNotification(notifyParam)
  if err != nil {
    return err
  }
  return nil
}
