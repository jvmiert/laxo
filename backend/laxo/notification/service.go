package notification

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/hashicorp/go-hclog"
	"github.com/mediocregopher/radix/v4"
	"gopkg.in/guregu/null.v4"
)

var ErrNotificationNeedsWorkflowOrGroup = errors.New("notification needs either a workflowID or groupID")

const (
  NotificationPrefix = "notifications_"

  EntityTypeProductAdd = "product_add"
)

type Store interface {
  SaveNotification(*NotificationCreateParam) (*Notification, error)
  UpdateNotificationRedisID(string, string) error
  GetNotificationGroupIDByWorkflowID(string, string) (string, error)
  CreateNotificationGroup(*NotificationGroupCreateParam) (string, error)
  UpdateNotificationGroup(*NotificationGroupUpdateParam) error
  GetNotifications(string, int32, int32) ([]Notification, error)
}

type NotificationGroupCreateParam struct {
  UserID           string
  WorkflowID       null.String
  EntityID         string
  EntityType       string
  TotalMainSteps   null.Int
  TotalSubSteps    null.Int
}

type NotificationGroupUpdateParam struct {
  UserID           null.String
  WorkflowID       null.String
  EntityID         null.String
  EntityType       null.String
  TotalMainSteps   null.Int
  TotalSubSteps    null.Int
  ID               string
}

type NotificationCreateParam struct {
  GroupID          string
  RedisID          null.String
  CurrentMainStep  null.Int
  CurrentSubStep   null.Int
  MainMessage      null.String
  SubMessage       null.String
  ReadTime         null.Time
}

type Service struct {
  store       Store
  logger      hclog.Logger
  redisClient radix.Client
}

func NewService(store Store, logger hclog.Logger, redisClient radix.Client) Service {
  return Service {
    store: store,
    logger: logger,
    redisClient: redisClient,
  }
}

func (s *Service) PublishNotification(n *Notification) (string, error) {
  bytes, err := n.JSON()
  if err != nil {
    return "", err
  }

  var StreamID radix.StreamEntryID
  s.redisClient.Do(context.Background(), radix.Cmd(
    &StreamID,
    "XADD",
    NotificationPrefix + n.GroupModel.UserID,
    "*",
    "notification", string(bytes),
  ))

  // expire the key in 2 hours
  s.redisClient.Do(context.Background(), radix.Cmd(
    nil,
    "EXPIRE",
    NotificationPrefix + n.GroupModel.UserID,
    "7200",
  ))

  return StreamID.String(), nil
}

func (s *Service) UpdateRedisIDToStore(notificationID string, redisID string) error {
  err := s.store.UpdateNotificationRedisID(notificationID, redisID)

  return err
}

func (s *Service) SaveNotificationToStore(param NotificationCreateParam) (*Notification, error) {
  n, err := s.store.SaveNotification(&param)

  return n, err
}

func (s *Service) UpdateNotificationGroup(param NotificationGroupUpdateParam) error {
  return s.store.UpdateNotificationGroup(&param)
}

func (s *Service) CreateNotificationGroup(param NotificationGroupCreateParam) (string, error) {
  return s.store.CreateNotificationGroup(&param)
}

func (s *Service) GetNotificationGroupIDByWorkflowID(workflowID, userID string) (string, error) {
  return s.store.GetNotificationGroupIDByWorkflowID(workflowID, userID)
}

func (s *Service) CreateNotification(param NotificationCreateParam) error {
  n, err := s.SaveNotificationToStore(param)
  if err != nil {
    return err
  }

  redisID, err := s.PublishNotification(n)
  if err != nil {
    return err
  }

  return s.UpdateRedisIDToStore(n.Model.ID, redisID)
}

func (s *Service) GetNotifications(userID string, offset, limit int32) ([]Notification, error) {
  nn, err := s.store.GetNotifications(userID, offset, limit)
  if err != nil {
    return nil, err
  }

  return nn, nil
}


func (s *Service) GetNotificationsJSON(userID string, offset, limit int32) ([]byte, error) {
  nn, err := s.GetNotifications(userID, offset, limit)
  if err != nil {
    return nil, err
  }

  nList := []json.RawMessage{}
  for _, s := range nn {
    b, errJ := s.JSON()
    if errJ != nil {
      return nil, errJ
    }
    j := json.RawMessage(b)
    nList = append(nList, j)
  }

	notificationData := map[string]interface{}{
		"notifications": nList,
		"total": len(nList),
	}

  bytes, err := json.Marshal(notificationData)

  if err != nil {
    s.logger.Error("Notification list marshal error", "error", err)
    return bytes, err
  }

  return bytes, nil
}

