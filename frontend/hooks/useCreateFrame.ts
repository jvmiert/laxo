import { grpc } from "@improbable-eng/grpc-web";
import { ProductService } from "@/proto/product_pb_service";
import { CreateFrameRequest } from "@/proto/product_pb";

// This will eventually be used to create a new frame. I want to use GRPC for this because
// the frame requests should update the user on the progress of creating frames for each product
// that the frame applies to. I envision that the user will have many products, so this might take
// a while. Using GRPC abstracts away the server send logic plus gives a nicely typed API interface.
//
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
