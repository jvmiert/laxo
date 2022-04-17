package proto

import (
	"context"
	"strconv"
	"time"

	"github.com/mediocregopher/radix/v4"
	"laxo.vn/laxo/laxo"
	gen "laxo.vn/laxo/laxo/proto/gen"
	"laxo.vn/laxo/processing"
)

type ProtoServer struct {
  gen.UnimplementedProductServiceServer
}

func (s *ProtoServer) CreateFrame(req *gen.CreateFrameRequest, stream gen.ProductService_CreateFrameServer) error {
  uID := stream.Context().Value(keyUID)
  laxo.Logger.Info("Received CreateFrame", "content", req.ImgID, "uID", uID)

  if err := stream.Send(&gen.CreateFrameReply{ImgID: "test1"}); err != nil {
    return err
  }

  time.Sleep(2 * time.Second)

  if err := stream.Send(&gen.CreateFrameReply{ImgID: "test2"}); err != nil {
    return err
  }

  return nil
}

func (s *ProtoServer) GetProductRetrieveUpdate(req *gen.ProductRetrieveUpdateRequest, stream gen.ProductService_GetProductRetrieveUpdateServer) error {
  uID := stream.Context().Value(keyUID)
  laxo.Logger.Info("Received GetProductRetrieveUpdate", "content", req.RetrieveID, "uID", uID)

  resp, err := laxo.TemporalClient.QueryWorkflow(context.Background(), req.RetrieveID, "", "current_state")
  if err != nil {
    laxo.Logger.Error("Unable to query workflow", "error", err)
  }

  var result processing.QueryStateResult

  if err := resp.Get(&result); err != nil {
    laxo.Logger.Error("Unable to decode query result", "error", err)
  }

  if err := stream.Send(&gen.ProductRetrieveUpdateReply{
    CurrentStatus: result.State,
    TotalProducts: int32(result.Total),
    CurrentProducts: int32(result.Current),
  }); err != nil {
    return err
  }

  if result.State == "complete" {
    return nil
  }

  var entries []radix.StreamEntry
  laxo.RedisClient.Do(context.Background(), radix.Cmd(&entries, "XRANGE", req.RetrieveID, "-", "+"))

  var latestID radix.StreamEntryID
  state := ""

  for _, e := range entries {
    latestID = e.ID

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
