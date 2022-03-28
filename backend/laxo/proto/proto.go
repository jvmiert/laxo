package proto

import (
	"time"

	"laxo.vn/laxo/laxo"
	gen "laxo.vn/laxo/laxo/proto/gen"
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

