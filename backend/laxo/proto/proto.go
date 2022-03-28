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
  laxo.Logger.Info("Received CreateFrame", req.GetImgID())

  if err := stream.Send(&gen.CreateFrameReply{ImgID: "test1"}); err != nil {
    return err
  }

  time.Sleep(2 * time.Second)

  if err := stream.Send(&gen.CreateFrameReply{ImgID: "test2"}); err != nil {
    return err
  }


  return nil
}

