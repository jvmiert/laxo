package store

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v4"
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

func (s *notificationStore) GetOrCreateProductGroup(param ProductGroupParam) (*sqlc.NotificationsGroup, error) {
  if !param.ID.Valid && !param.WorkflowID.Valid {
    return nil, ErrNotificationGroupCreateUpdate
  }

  var g sqlc.NotificationsGroup
  var err error

  if param.ID.Valid {
    g, err = s.queries.GetNotificationsGroupByID(
      context.Background(),
      param.ID.String,
    )

    if err != pgx.ErrNoRows && err != nil {
      return nil, err
    }
  }

  if param.WorkflowID.Valid {
    getParam := sqlc.GetNotificationsGroupByWorkflowIDParams{
      WorkflowID: param.WorkflowID.String,
      UserID: param.UserID,
    }

    g, err = s.queries.GetNotificationsGroupByWorkflowID(
      context.Background(),
      getParam,
    )

    if err != pgx.ErrNoRows && err != nil {
      return nil, err
    }
  }

  // We didn't retrieve a valid notification group
  if g.ID == "" {
    if !param.WorkflowID.Valid {
      return nil, ErrNotificationGroupCreateUpdate
    }

    createParam := sqlc.CreateNotificationsGroupParams{
      UserID: param.UserID,
      WorkflowID: param.WorkflowID.String,
      EntityID: param.EntityID,
      EntityType: param.EntityType,
      TotalMainSteps: param.TotalMainSteps.NullInt64,
      TotalSubSteps: param.TotalSubSteps.NullInt64,
    }

    g, err = s.queries.CreateNotificationsGroup(
      context.Background(),
      createParam,
    )

    if err != nil {
      return nil, err
    }
  }

  return &g, nil
}

func (s *notificationStore) UpdateNotificationRedisID(notificationID string, redisID string) error {
  updateParam :=  sqlc.UpdateRedisIDByNotificationIDParams{
    RedisID: null.NewString(redisID, true).NullString,
    ID:      notificationID,
  }
  err := s.queries.UpdateRedisIDByNotificationID(
    context.Background(),
    updateParam,
  )

  return err
}

func (s *notificationStore) SaveNotification(p *notification.AddNotificationParam) (*notification.Notification, error) {
  groupParam := ProductGroupParam{
    ID: p.GroupID,
    WorkflowID: p.WorkflowID,
    UserID: p.UserID,
    EntityID: p.EntityID,
    EntityType: p.EntityType,
    TotalMainSteps: p.TotalMainSteps,
    TotalSubSteps: p.TotalSubSteps,
  }

  group, err := s.GetOrCreateProductGroup(groupParam)
  if err != nil {
    return nil, err
  }

  createParam := sqlc.CreateNotificationParams {
    NotificationGroupID: group.ID,
    Read: p.ReadTime.NullTime,
    CurrentMainStep: p.CurrentMainStep.NullInt64,
    CurrentSubStep: p.CurrentSubStep.NullInt64,
    MainMessage: p.MainMessage.NullString,
    SubMessage: p.SubMessage.NullString,
  }

  pModel, err := s.queries.CreateNotification(
    context.Background(),
    createParam,
  )
  if err != nil {
    return nil, err
  }

  return &notification.Notification{
    Model: &pModel,
    GroupModel: group,
  }, nil
}

