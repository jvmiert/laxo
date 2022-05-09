import { grpc } from "@improbable-eng/grpc-web";
import { UserService } from "@/proto/user_pb_service";
import {
  NotificationUpdateRequest,
  NotificationUpdateReply,
} from "@/proto/user_pb";

const myTransport = grpc.CrossBrowserHttpTransport({ withCredentials: true });

export default function useNotificationApi(): {
  getNotificationUpdate: (
    notificationRedisID: string,
    onMessage: (res: NotificationUpdateReply) => void,
    onEnd: (code: grpc.Code, message: string, trailers: grpc.Metadata) => void,
  ) => grpc.Request;
} {
  const getNotificationUpdate = (
    notificationRedisID: string,
    onMessage: (res: NotificationUpdateReply) => void,
    onEnd: (code: grpc.Code, message: string, trailers: grpc.Metadata) => void,
  ): grpc.Request => {
    const notificationUpdateRequest = new NotificationUpdateRequest();

    notificationUpdateRequest.setNotificationredisid(notificationRedisID);

    const request = grpc.invoke(UserService.GetNotificationUpdate, {
      request: notificationUpdateRequest,
      host: "http://localhost:8081",
      transport: myTransport,
      onMessage: onMessage,
      onEnd: onEnd,
    });

    return request;
  };

  return { getNotificationUpdate };
}
