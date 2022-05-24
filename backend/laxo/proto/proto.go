package proto

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/mediocregopher/radix/v4"
	"laxo.vn/laxo/laxo"
	"laxo.vn/laxo/laxo/notification"
	gen "laxo.vn/laxo/laxo/proto/gen"
)

type ProtoServer struct {
  gen.UnimplementedUserServiceServer
  service     *notification.Service
  logger      *laxo.Logger
  redisClient radix.Client
  ctx         context.Context
  server      *laxo.Server
}

func NewServer(service *notification.Service, logger *laxo.Logger, redisURI string, ctx context.Context, server *laxo.Server) (*ProtoServer, error) {
  client, err := (radix.PoolConfig{
    Size: 50,
  }).New(ctx, "tcp", redisURI)
  if err != nil {
    return nil, err
  }

  return &ProtoServer {
    service: service,
    logger: logger,
    redisClient: client,
    ctx: ctx,
    server: server,
  }, nil
}

func (s *ProtoServer) GetNotificationUpdate(req *gen.NotificationUpdateRequest, stream gen.UserService_GetNotificationUpdateServer) error {
  uID := stream.Context().Value(keyUID).(string)
  s.logger.Infow("Received GetNotificationUpdate",
    "NotificationRedisID", req.NotificationRedisID,
    "uID", uID,
  )

  channelID := notification.NotificationPrefix + uID

  var latestID radix.StreamEntryID
  streamConfig := make(map[string]radix.StreamConfig)
  sc := radix.StreamConfig{}

  if req.NotificationRedisID != "" {
    sRedisID := strings.Split(req.NotificationRedisID, "-")

    time, err := strconv.ParseUint(sRedisID[0], 10, 64)
    if err != nil {
      s.logger.Errorw("strconv error",
        "error", err,
      )
      return err
    }

    seq, err := strconv.ParseUint(sRedisID[1], 10, 64)
    if err != nil {
      s.logger.Errorw("strconv error",
        "error", err,
      )
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

  timeoutCtx, cancelTimeout := context.WithCancel(s.ctx)
  defer cancelTimeout()

  go func(ctx context.Context) {
    defer cancelTimeout()
    d := time.NewTicker(5 * time.Second)

    for {
      select {
        case <-d.C:
          err := stream.Send(&gen.NotificationUpdateReply{
            KeepAlive: true,
          })
          if err != nil {
            return
          }
        case <-ctx.Done():
          return
      }
    }
  }(timeoutCtx)

  for {
    ctx, cancel := context.WithCancel(timeoutCtx)
    defer cancel()

    select {
      case <-timeoutCtx.Done():
        return nil
      case <-s.ctx.Done():
        return nil
      default:
    }

    _, entry, err := r.Next(ctx)
    cancel()
    if err != nil {
      if err != radix.ErrNoStreamEntries {
        s.logger.Errorw("Redis stream Next() returned error",
          "error", err,
        )
        return err
      }
    }

    if err != radix.ErrNoStreamEntries {
      if len(entry.Fields) > 0 {
        if entry.Fields[0][0] == "notification" {
          var n notification.Notification

          if err = json.Unmarshal([]byte(entry.Fields[0][1]), &n); err != nil {
            s.logger.Errorw("notification json unmarshal error",
              "error", err,
            )
            return err
          }

          nn := &gen.Notification{
            ID: n.Model.ID,
            RedisID: n.Model.RedisID.ValueOrZero(),
            GroupID: n.Model.NotificationGroupID,
            Created: n.Model.Created.ValueOrZero().Unix(),
            CurrentMainStep: n.Model.CurrentMainStep.ValueOrZero(),
            MainMessage: n.Model.MainMessage.ValueOrZero(),
            SubMessage: n.Model.SubMessage.ValueOrZero(),
          }

          if n.Model.Read.Valid {
            read := n.Model.Read.ValueOrZero().Unix()
            nn.Read = &read
          }

          if n.Model.CurrentSubStep.Valid {
            nn.CurrentSubStep = &n.Model.CurrentSubStep.Int64
          }

          ng := &gen.NotificationGroup{
            ID: n.GroupModel.ID,
            UserID: n.GroupModel.UserID,
            WorkflowID: n.GroupModel.WorkflowID.ValueOrZero(),
            EntityID: n.GroupModel.EntityID,
            EntityType: n.GroupModel.EntityType,
            TotalMainSteps: n.GroupModel.TotalMainSteps.ValueOrZero(),
          }

          if n.GroupModel.TotalSubSteps.Valid {
            ng.TotalSubSteps = &n.GroupModel.TotalSubSteps.Int64
          }

          if err = stream.Send(&gen.NotificationUpdateReply{
            Notification: nn,
            NotificationGroup: ng,
          }); err != nil {
            return err
          }
        }
      }
    }
  }
}
