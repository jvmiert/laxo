package notification

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/mediocregopher/radix/v4"
	"gopkg.in/guregu/null.v4"
	"laxo.vn/laxo/laxo"
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
  PlatformName     string
}

type NotificationGroupUpdateParam struct {
  UserID           null.String
  WorkflowID       null.String
  EntityID         null.String
  EntityType       null.String
  TotalMainSteps   null.Int
  TotalSubSteps    null.Int
  ID               string
  PlatformName     null.String
}

type NotificationCreateParam struct {
  GroupID          string
  RedisID          null.String
  CurrentMainStep  null.Int
  CurrentSubStep   null.Int
  MainMessage      null.String
  SubMessage       null.String
  ReadTime         null.Time
  Error            null.Bool
}

type Service struct {
  store       Store
  logger      *laxo.Logger
  server      *laxo.Server
}

func NewService(store Store, logger *laxo.Logger, server *laxo.Server) Service {
  return Service {
    store: store,
    logger: logger,
    server: server,
  }
}

func (s *Service) PublishNotification(n *Notification) (string, error) {
  bytes, err := n.JSON()
  if err != nil {
    return "", err
  }

  var StreamID radix.StreamEntryID
  s.server.RedisClient.Do(context.Background(), radix.Cmd(
    &StreamID,
    "XADD",
    NotificationPrefix + n.GroupModel.UserID,
    "*",
    "notification", string(bytes),
  ))

  // expire the key in 2 hours
  s.server.RedisClient.Do(context.Background(), radix.Cmd(
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
    s.logger.Errorw("Notification list marshal error",
      "error", err,
    )
    return bytes, err
  }

  return bytes, nil
}

func (s *Service) ErrorNotification(groupID string) error {
  notifyParam := NotificationCreateParam{
    GroupID: groupID,
    CurrentMainStep: null.IntFrom(1),
    CurrentSubStep: null.NewInt(0, false),
    MainMessage: null.StringFrom(""),
    SubMessage: null.NewString("", false),
    ReadTime: null.NewTime(time.Time{}, false),
    Error: null.BoolFrom(true),
  }

  err := s.CreateNotification(notifyParam)
  if err != nil {
    return err
  }

  return nil
}
