import { grpc } from "@improbable-eng/grpc-web";
import { ProductService } from "@/proto/product_pb_service";
import { CreateFrameRequest } from "@/proto/product_pb";

//@TODO: export this from another file and make it default transport
//       as per: https://github.com/improbable-eng/grpc-web/blob/master/client/grpc-web/docs/transport.md#specifying-the-default-transport
const myTransport = grpc.CrossBrowserHttpTransport({ withCredentials: true });

export default function useCreateFrame(): {
  createFrame: () => void;
} {
  const createFrame = async () => {
    const createFrameRequest = new CreateFrameRequest();
    createFrameRequest.setImgid("hello!");

    grpc.invoke(ProductService.CreateFrame, {
      request: createFrameRequest,
      host: "http://localhost:8081",
      transport: myTransport,
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
