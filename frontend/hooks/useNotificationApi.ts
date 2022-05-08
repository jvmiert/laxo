import { grpc } from "@improbable-eng/grpc-web";
import { UserService } from "@/proto/user_pb_service";
import { NotificationUpdateRequest } from "@/proto/user_pb";

const myTransport = grpc.CrossBrowserHttpTransport({ withCredentials: true });

export default function useNotificationApi(): {
  getNotificationUpdate: (
    notificationID: string,
    notificationGroupID: string,
  ) => void;
} {
  const getNotificationUpdate = async (
    notificationID: string,
    notificationGroupID: string,
  ) => {
    const notificationUpdateRequest = new NotificationUpdateRequest();

    notificationUpdateRequest.setNotificationid(notificationID);
    notificationUpdateRequest.setNotificationgroupid(notificationGroupID);

    grpc.invoke(UserService.GetNotificationUpdate, {
      request: notificationUpdateRequest,
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

  return { getNotificationUpdate };
}
