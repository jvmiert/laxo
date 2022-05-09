package proto

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/mediocregopher/radix/v4"
	"laxo.vn/laxo/laxo/notification"
	gen "laxo.vn/laxo/laxo/proto/gen"
)

type ProtoServer struct {
  gen.UnimplementedUserServiceServer
  service *notification.Service
  logger      hclog.Logger
  redisClient radix.Client
}

func NewServer(service *notification.Service, logger hclog.Logger, redisClient radix.Client) *ProtoServer {
  return &ProtoServer {
    service: service,
    logger: logger,
    redisClient: redisClient,
  }
}

func (s *ProtoServer) GetNotificationUpdate(req *gen.NotificationUpdateRequest, stream gen.UserService_GetNotificationUpdateServer) error {
  uID := stream.Context().Value(keyUID).(string)
  s.logger.Info("Received GetProductRetrieveUpdate", "NotificationRedisID", req.NotificationRedisID, "uID", uID)

  channelID := notification.NotificationPrefix + uID

  var latestID radix.StreamEntryID
  streamConfig := make(map[string]radix.StreamConfig)
  sc := radix.StreamConfig{}

  if req.NotificationRedisID != "" {
    sRedisID := strings.Split(req.NotificationRedisID, "-")

    time, err := strconv.ParseUint(sRedisID[0], 10, 64)
    if err != nil {
      s.logger.Error("strconv error", "error", err)
      return err
    }

    seq, err := strconv.ParseUint(sRedisID[1], 10, 64)
    if err != nil {
      s.logger.Error("strconv error", "error", err)
      return err
    }

    latestID.Time = time
    latestID.Seq = seq

    sc.After = latestID
  } else {
    sc.Latest = true
  }

  streamReaderConfig := radix.StreamReaderConfig{
    Group:     "",
    Consumer:  "",
    Count:     -1,
    NoBlock: false,
  }

  streamConfig[channelID] = sc

  r := streamReaderConfig.New(s.redisClient, streamConfig)

  keepAliveErrs := make(chan error, 1)
  go func() {
    for {
      time.Sleep(30 * time.Second)

      err := stream.Send(&gen.NotificationUpdateReply{
        KeepAlive: true,
      })
      if err != nil {
        keepAliveErrs <- err
        close(keepAliveErrs)
        return
      }
    }
  }()

  for {
    select {
    case msg := <-keepAliveErrs:
      s.logger.Error("Keepalive error", "error", msg)
      return msg
    default:
    }

    ctx, cancel := context.WithTimeout(context.Background(), 2 * time.Minute)
    _, entry, err := r.Next(ctx)

    if err != nil {
      if err != radix.ErrNoStreamEntries {
        s.logger.Error("Redis stream Next() returned error", "error", err)
        cancel()
        return err
      }
    }

    cancel()

    if err != radix.ErrNoStreamEntries {
      if len(entry.Fields) > 0 {
        if entry.Fields[0][0] == "notification" {
          var n notification.Notification

          if err = json.Unmarshal([]byte(entry.Fields[0][1]), &n); err != nil {
            s.logger.Error("notification json unmarshal error", "error", err)
            return err
          }

          //@TODO: The ValueOrZero approach is not correct for Created and Read

          if err = stream.Send(&gen.NotificationUpdateReply{
            Notification: &gen.Notification{
              ID: n.Model.ID,
              RedisID: n.Model.RedisID.ValueOrZero(),
              GroupID: n.Model.NotificationGroupID,
              Created: n.Model.Created.ValueOrZero().Unix(),
              Read: n.Model.Read.ValueOrZero().Unix(),
              CurrentMainStep: n.Model.CurrentMainStep.ValueOrZero(),
              CurrentSubStep: n.Model.CurrentSubStep.ValueOrZero(),
              MainMessage: n.Model.MainMessage.ValueOrZero(),
              SubMessage: n.Model.SubMessage.ValueOrZero(),
            },
            NotificationGroup: &gen.NotificationGroup{
              ID: n.GroupModel.ID,
              UserID: n.GroupModel.UserID,
              WorkflowID: n.GroupModel.WorkflowID.ValueOrZero(),
              EntityID: n.GroupModel.EntityID,
              EntityType: n.GroupModel.EntityType,
              TotalMainSteps: n.GroupModel.TotalMainSteps.ValueOrZero(),
              TotalSubSteps: n.GroupModel.TotalSubSteps.ValueOrZero(),
            },
          }); err != nil {
            return err
          }
        }
      }
    }
  }
}
