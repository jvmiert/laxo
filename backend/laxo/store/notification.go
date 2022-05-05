package store

import (
	"context"
	"errors"

	"gopkg.in/guregu/null.v4"
	"laxo.vn/laxo/laxo/notification"
	"laxo.vn/laxo/laxo/sqlc"
)

var ErrNotificationGroupCreateUpdate = errors.New("unable to create or update a notification group")

type ProductGroupParam struct {
  ID               null.String
  WorkflowID       null.String
  UserID           string
  EntityID         string
  EntityType       string
  TotalMainSteps   null.Int
  TotalSubSteps    null.Int
}

type notificationStore struct {
  *Store
}

func newNotificationStore(store *Store) notificationStore {
  return notificationStore{
    store,
  }
}

func (s *notificationStore) UpdateNotificationRedisID(notificationID string, redisID string) error {
  updateParam :=  sqlc.UpdateRedisIDByNotificationIDParams{
    RedisID: null.StringFrom(redisID),
    ID:      notificationID,
  }
  err := s.queries.UpdateRedisIDByNotificationID(
    context.Background(),
    updateParam,
  )

  return err
}

func (s *notificationStore) UpdateNotificationGroup(param *notification.NotificationGroupUpdateParam) error {
  updateParam := sqlc.UpdateNotificationGroupParams{
    UserIDDoUpdate: param.UserID.Valid,
    UserID: param.UserID.String,
    WorkflowIDDoUpdate: param.WorkflowID.Valid,
    WorkflowID: param.WorkflowID.String,
    EntityIDDoUpdate: param.EntityID.Valid,
    EntityID: param.EntityID.String,
    EntityTypeDoUpdate: param.EntityType.Valid,
    EntityType: param.EntityType.String,
    TotalMainStepsDoUpdate: param.TotalMainSteps.Valid,
    TotalMainSteps: param.TotalMainSteps.Int64,
    TotalSubStepsDoUpdate: param.TotalSubSteps.Valid,
    TotalSubSteps: param.TotalSubSteps.Int64,
    ID: param.ID,
  }

  _, err := s.queries.UpdateNotificationGroup(
    context.Background(),
    updateParam,
  )

  return err
}

func (s *notificationStore) CreateNotificationGroup(param *notification.NotificationGroupCreateParam) (string, error) {
    createParam := sqlc.CreateNotificationsGroupParams{
      UserID: param.UserID,
      WorkflowID: param.WorkflowID,
      EntityID: param.EntityID,
      EntityType: param.EntityType,
      TotalMainSteps: param.TotalMainSteps,
      TotalSubSteps: param.TotalSubSteps,
    }

  g, err := s.queries.CreateNotificationsGroup(
    context.Background(),
    createParam,
  )

  if err != nil {
    return "", err
  }

  return g.ID, nil
 }

func (s *notificationStore) GetNotificationGroupIDByWorkflowID(workflowID, userID string) (string, error) {
  param := sqlc.GetNotificationsGroupByWorkflowIDParams{
    WorkflowID: null.StringFrom(workflowID),
    UserID: userID,
  }
  n, err := s.queries.GetNotificationsGroupByWorkflowID(
    context.Background(),
    param,
  )
  if err != nil {
    return "", err
  }

  return n.ID, nil
}

func (s *notificationStore) SaveNotification(p *notification.NotificationCreateParam) (*notification.Notification, error) {
  createParam := sqlc.CreateNotificationParams {
    NotificationGroupID: p.GroupID,
    Read: p.ReadTime,
    CurrentMainStep: p.CurrentMainStep,
    CurrentSubStep: p.CurrentSubStep,
    MainMessage: p.MainMessage,
    SubMessage: p.SubMessage,
  }

  pModel, err := s.queries.CreateNotification(
    context.Background(),
    createParam,
  )
  if err != nil {
    return nil, err
  }

  gModel, err := s.queries.GetNotificationsGroupByID(
    context.Background(),
    pModel.NotificationGroupID,
  )
  if err != nil {
    return nil, err
  }

  return &notification.Notification{
    Model: &pModel,
    GroupModel: &gModel,
  }, nil
}

func (s *notificationStore) GetNotifications(userID string, offset, limit int32) ([]notification.Notification, error) {
  getParam := sqlc.GetNotificationsByUserIDParams {
    UserID: userID,
    Offset: offset,
    Limit: limit,
  }

  results, err := s.queries.GetNotificationsByUserID(
    context.Background(),
    getParam,
  )
  if err != nil {
    return nil, err
  }

  returnList := []notification.Notification{}

  for _, v := range results {
    returnList = append(
      returnList,
      notification.Notification{
        Model: &sqlc.Notification{
          ID: v.NotificationID,
          RedisID: v.NotificationRedisID,
          NotificationGroupID: v.ID,
          Created: v.NotificationCreated,
          Read: v.NotificationRead,
          CurrentMainStep: v.NotificationCurrentMainStep,
          CurrentSubStep: v.NotificationCurrentSubStep,
          MainMessage: v.NotificationMainMessage,
          SubMessage: v.NotificationSubMessage,
        },
        GroupModel: &sqlc.NotificationsGroup{
          ID: v.ID,
          UserID: v.UserID,
          WorkflowID: v.WorkflowID,
          EntityID: v.EntityID,
          EntityType: v.EntityType,
          TotalMainSteps: v.TotalMainSteps,
          TotalSubSteps: v.TotalSubSteps,
        },
      },
    )
  }

  return returnList, nil
}
