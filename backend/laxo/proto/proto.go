package proto

import (
	"context"
	"time"

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

  for {
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
      TotalProducts: uint32(result.Total),
      CurrentProducts: uint32(result.Current),
    }); err != nil {
      return err
    }

    if result.State == "complete" {
      break;
    }

    time.Sleep(1 * time.Second)
  }

  return nil
}
