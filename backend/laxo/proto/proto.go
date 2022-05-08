package proto

import (
	"context"
	"strconv"
	"time"

	"github.com/mediocregopher/radix/v4"
	"laxo.vn/laxo/laxo"
	gen "laxo.vn/laxo/laxo/proto/gen"
)

type ProtoServer struct {
  gen.UnimplementedUserServiceServer
}

//@TODO: adjust below function to read the current latest notification and start listening to Redis

func (s *ProtoServer) GetNotificationUpdate(req *gen.NotificationUpdateRequest, stream gen.UserService_GetNotificationUpdateServer) error {
  uID := stream.Context().Value(keyUID)
  laxo.Logger.Info("Received GetProductRetrieveUpdate", "NotificationID", req.NotificationID, "NotificationGroupID", req.NotificationGroupID, "uID", uID)

  var entries []radix.StreamEntry
  laxo.RedisClient.Do(context.Background(), radix.Cmd(&entries, "XRANGE", req.RetrieveID, "-", "+"))

  latestID := entries[len(entries)-1].ID
  state := ""

  for _, e := range entries {
    state = ""
    total, complete := -1, -1
    for _, f := range e.Fields {
      key := f[0]
      value := f[1]

      if key == "state" {
        state = value
      }

      if key == "complete" {
        i, errConvert := strconv.Atoi(value)
        if errConvert != nil {
          laxo.Logger.Error("couldn't convert complete string to int", "error", errConvert)
        }

        complete = i
      }

      if key == "total" {
        i, errConvert := strconv.Atoi(value)
        if errConvert != nil {
          laxo.Logger.Error("couldn't convert total string to int", "error", errConvert)
        }

        total = i
      }
    }

    if state != "" {
      if errStream := stream.Send(&gen.ProductRetrieveUpdateReply{
        CurrentStatus: state,
        TotalProducts: int32(total),
        CurrentProducts: int32(complete),
      }); errStream != nil {
        return errStream
      }
    }
  }

  if state == "complete" {
    return nil
  }

  streamReaderConfig := radix.StreamReaderConfig{
    Group:     "",
    Consumer:  "",
    Count:     -1,
    NoBlock: false,
  }
  r := streamReaderConfig.New(laxo.RedisClient, map[string]radix.StreamConfig{
    req.RetrieveID: {
      After: latestID,
      Latest: true,
    },
  })

  state = ""

  for state != "complete" {
    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)

    _, entry, err := r.Next(ctx)

    if err != nil {
      if err != radix.ErrNoStreamEntries {
        laxo.Logger.Error("Redis stream Next() returned error", "error", err)
      } else {
        laxo.Logger.Error("ErrNoStreamEntries")
      }
    }

    cancel()

    total, complete := -1, -1

    for _, f := range entry.Fields {
      key := f[0]
      value := f[1]

      if key == "state" {
        state = value
      }

      if key == "complete" {
        i, errConvert := strconv.Atoi(value)
        if errConvert != nil {
          laxo.Logger.Error("couldn't convert complete string to int", "error", errConvert)
        }

        complete = i
      }

      if key == "total" {
        i, errConvert := strconv.Atoi(value)
        if errConvert != nil {
          laxo.Logger.Error("couldn't convert total string to int", "error", errConvert)
        }

        total = i
      }

    }

    if state != "" {
      laxo.Logger.Debug("sending new state", "total", total, "complete", complete, "state", state)
      if errStream := stream.Send(&gen.ProductRetrieveUpdateReply{
        CurrentStatus: state,
        TotalProducts: int32(total),
        CurrentProducts: int32(complete),
      }); errStream != nil {
        return errStream
      }
    }
  }

  return nil
}
