import { grpc } from "@improbable-eng/grpc-web";
import { ProductService } from "@/proto/product_pb_service";
import { CreateFrameRequest } from "@/proto/product_pb";

export default function useCreateFrame(): {
  createFrame: () => void;
} {
  const createFrame = async () => {
    const createFrameRequest = new CreateFrameRequest();
    createFrameRequest.setImgid("hello!");

    grpc.invoke(ProductService.CreateFrame, {
      request: createFrameRequest,
      host: "http://localhost:8081",
      onMessage: (res) => {
        console.log("onMessage", res);
      },
      onEnd: (res) => {
        console.log("onEnd", res);
      },
    });
  };

  return { createFrame };
}
