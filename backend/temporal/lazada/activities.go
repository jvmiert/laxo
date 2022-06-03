package lazada

import (
	"context"
	"time"

	"github.com/mediocregopher/radix/v4"
	"go.temporal.io/sdk/activity"
	"gopkg.in/guregu/null.v4"
	"laxo.vn/laxo/laxo/assets"
	"laxo.vn/laxo/laxo/lazada"
	"laxo.vn/laxo/laxo/notification"
	"laxo.vn/laxo/laxo/shop"
)

const (
  ActivityStateFetch = "fetch"
  ActivityStateSave = "save"
  ActivityStateComplete = "complete"
)

type LazadaFetchResult struct {
  DataKey         string
  TotalProducts   int64
}

type LazadaSaveParam struct {
  UserID              string
  ShopID              string
  DataKey             string
  ProductIndex        int64
  ProductTotal        int64
}

type Activities struct {
  redisClient         *radix.Client
  notificationService *notification.Service
  lazadaService       *lazada.Service
  shopService         *shop.Service
  assetsService       *assets.Service
}

func NewActivities(r *radix.Client, n *notification.Service, l *lazada.Service, s *shop.Service, a *assets.Service) *Activities {
  return &Activities{
    redisClient: r,
    notificationService: n,
    lazadaService: l,
    shopService: s,
    assetsService: a,
  }
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
    PlatformName: "lazada",
  }

  notificationGroupID, err := a.notificationService.CreateNotificationGroup(notifyGroupParam)
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
    Error: null.BoolFrom(false),
  }

  err = a.notificationService.CreateNotification(notifyParam)
  if err != nil {
    return LazadaFetchResult{}, err
  }

  logger.Info("Retrieving Lazada products", "shopID", shopID)
  key, total, err := a.lazadaService.FetchProductsFromLazadaToRedis(shopID)
  if err != nil {
    logger.Error("FetchProductsFromLazadaToRedis error", "error", err)
    a.notificationService.ErrorNotification(notificationGroupID)
    return LazadaFetchResult{}, err
  }

  updateParam := notification.NotificationGroupUpdateParam{
    TotalSubSteps: null.IntFrom(int64(total)),
    ID: notificationGroupID,
  }

  err = a.notificationService.UpdateNotificationGroup(updateParam)
  if err != nil {
    return LazadaFetchResult{}, err
  }

  notifyParam.MainMessage = null.StringFrom(ActivityStateSave)
  notifyParam.CurrentMainStep = null.IntFrom(2)
  notifyParam.CurrentSubStep = null.IntFrom(0)

  err = a.notificationService.CreateNotification(notifyParam)
  if err != nil {
    return LazadaFetchResult{}, err
  }

  return LazadaFetchResult{DataKey: key, TotalProducts: int64(total)}, nil
}

func (a *Activities) SaveLazadaProducts(ctx context.Context, param LazadaSaveParam) error {
  logger := activity.GetLogger(ctx)
  info := activity.GetInfo(ctx)
  workflowID := info.WorkflowExecution.ID

  notificationGroupID, err := a.notificationService.GetNotificationGroupIDByWorkflowID(workflowID, param.UserID)
  if err != nil {
    return err
  }

  logger.Info("Save Lazada products", "dataKeyID", param.DataKey, "index", param.ProductIndex)

  p, err := a.lazadaService.RetrieveProductFromRedis(param.DataKey, int(param.ProductIndex))
  if err != nil {
    logger.Error("RetrieveProductFromRedis error",
      "error", err,
    )
    a.notificationService.ErrorNotification(notificationGroupID)
    return err
  }

  shop, err := a.shopService.GetShopByID(param.ShopID)
  if err != nil {
    logger.Error("GetLaxoProductFromLazadaData error",
      "error", err,
    )
    a.notificationService.ErrorNotification(notificationGroupID)
    return err
  }

  pModel, pModelAttributes, pModelSKU, err := a.lazadaService.SaveOrUpdateLazadaProduct(p, param.ShopID)
  if err != nil {
    logger.Error("SaveOrUpdateLazadaProduct error",
      "error", err,
    )
    a.notificationService.ErrorNotification(notificationGroupID)
    return err
  }

  product, err := a.shopService.GetLaxoProductFromLazadaData(pModel, pModelAttributes, pModelSKU)
  if err != nil {
    logger.Error("GetLaxoProductFromLazadaData error",
      "error", err,
    )
    a.notificationService.ErrorNotification(notificationGroupID)
    return err
  }

  laxoP, err := a.shopService.SaveOrUpdateProductToStore(product, param.ShopID, pModel.ID)
  if err != nil {
    logger.Error("SaveOrUpdateProductToStore error",
      "error", err,
    )
    a.notificationService.ErrorNotification(notificationGroupID)
    return err
  }

  images, err := a.assetsService.ExtractImagesListFromProductResponse(p)
  if err != nil {
    logger.Error("ExtractImagesListFromProductResponse error",
      "error", err,
    )
    a.notificationService.ErrorNotification(notificationGroupID)
    return err
  }

  err = a.assetsService.SaveProductImages(images, param.ShopID, laxoP.Model.ID, shop.Model.AssetsToken)
  if err != nil {
    logger.Error("SaveProductImages error",
      "error", err,
    )
    a.notificationService.ErrorNotification(notificationGroupID)
    return err
  }

  notifyParam := notification.NotificationCreateParam{
    GroupID: notificationGroupID,
    CurrentMainStep: null.IntFrom(2),
    CurrentSubStep: null.IntFrom(int64(param.ProductIndex + 1)),
    MainMessage: null.StringFrom(ActivityStateSave),
    SubMessage: null.NewString("", false),
    ReadTime: null.NewTime(time.Time{}, false),
    Error: null.BoolFrom(false),
  }

  err = a.notificationService.CreateNotification(notifyParam)
  if err != nil {
    return err
  }

  return nil
}

func (a *Activities) CompleteLazadaProducts(ctx context.Context, userID string, keyID string) error {
  logger := activity.GetLogger(ctx)
  info := activity.GetInfo(ctx)
  workflowID := info.WorkflowExecution.ID

  logger.Info("Complete Lazada products")

  notificationGroupID, err := a.notificationService.GetNotificationGroupIDByWorkflowID(workflowID, userID)
  if err != nil {
    return err
  }

  err = a.lazadaService.ExpireRedisProducts(keyID)
  if err != nil {
    return err
  }

  notifyParam := notification.NotificationCreateParam{
    GroupID: notificationGroupID,
    CurrentMainStep: null.IntFrom(2),
    MainMessage: null.StringFrom(ActivityStateComplete),
    Error: null.BoolFrom(false),
  }

  err = a.notificationService.CreateNotification(notifyParam)
  if err != nil {
    return err
  }
  return nil
}
