import { ReactNode, useState, useEffect, Dispatch } from "react";
import { Draft } from "immer";
import { useImmerReducer } from "use-immer";
import { grpc } from "@improbable-eng/grpc-web";
import createSafeContext from "@/lib/useSafeContext";
import { useGetNotifications } from "@/hooks/swrHooks";
import useNotificationApi from "@/hooks/useNotificationApi";
import type { NotificationResponseObject } from "@/types/ApiResponse";
import { NotificationUpdateReply } from "@/proto/user_pb";

export interface DashboardConsumerProps {
  notificationOpen: boolean;
  closeNotification: () => void;
  openNotification: () => void;
  toggleNotification: () => void;
  dashboardState: DashboardState;
  dashboardDispatch: Dispatch<DashboardAction>;
}

interface DashboardState {
  notifications: Array<NotificationResponseObject>;
}

const initialState: DashboardState = {
  notifications: [],
};

export type DashboardAction =
  | { type: "reset"; state: DashboardState }
  | { type: "add"; notification: NotificationResponseObject };

function reducer(draft: Draft<DashboardState>, action: DashboardAction) {
  switch (action.type) {
    case "reset":
      return action.state;
    case "add":
      const index = draft.notifications.findIndex(
        (n) =>
          n.notificationGroup.id === action.notification.notificationGroup.id,
      );
      if (index !== -1) {
        draft.notifications[index] = action.notification;
      } else {
        draft.notifications.unshift(action.notification);
      }
      break;
  }
}

const onEnd = (code: grpc.Code, message: string, trailers: grpc.Metadata) => {};

const onMessageFunc = (dispatch: Dispatch<DashboardAction>) => {
  return (res: NotificationUpdateReply) => {
    const keepAlive = res.getKeepalive();
    const hasNotification = res.hasNotification();
    const hasNotificationgroup = res.hasNotificationgroup();

    if (keepAlive || !hasNotification || !hasNotificationgroup) {
      return;
    }

    const notification = res.getNotification();
    const notificationGroup = res.getNotificationgroup();

    if (!notification || !notificationGroup) {
      return;
    }

    const notiObject = notification.toObject();
    const notiGroupObject = notificationGroup.toObject();

    const notiResponseObject: NotificationResponseObject = {
      notification: {
        id: notiObject.id,
        redisID: notiObject.redisid,
        notificationGroupID: notiObject.groupid,
        created: new Date(notiObject.created * 1000),
        read:
          notiObject.read > 0 ? new Date(notiObject.read * 1000) : undefined,
        currentMainStep: notiObject.currentmainstep,
        currentSubStep: notiObject.currentsubstep,
        mainMessage: notiObject.mainmessage,
        subMessage: notiObject.submessage,
      },
      notificationGroup: {
        id: notiGroupObject.id,
        userID: notiGroupObject.userid,
        workflowID: notiGroupObject.workflowid,
        entityID: notiGroupObject.entityid,
        entityType: notiGroupObject.entitytype,
        totalMainSteps: notiGroupObject.totalmainsteps,
        totalSubSteps: notiGroupObject.totalsubsteps,
      },
    };

    dispatch({
      type: "add",
      notification: notiResponseObject,
    });
  };
};

export const [useDashboard, Provider] =
  createSafeContext<DashboardConsumerProps>();

export const DashboardProvider = ({ children }: { children: ReactNode }) => {
  const [notificationOpen, setNotificationOpen] = useState(false);

  const closeNotification = () => setNotificationOpen(false);
  const openNotification = () => setNotificationOpen(true);
  const toggleNotification = () => setNotificationOpen(!notificationOpen);

  const [notificationStreamActive, setNotificationStreamActive] =
    useState(false);

  const [state, dispatch] = useImmerReducer(reducer, initialState);

  const { getNotificationUpdate } = useNotificationApi();

  const {
    notifications,
    loading: notificationLoading,
    error: notificationError,
  } = useGetNotifications();

  useEffect(() => {
    if (notifications.notifications.length > 0) {
      dispatch({
        type: "reset",
        state: { notifications: notifications.notifications },
      });
    }
  }, [notifications, dispatch]);

  useEffect(() => {
    //@TODO: implement cleanup - close the notification request
    if (
      !notificationLoading &&
      !notificationError &&
      !notificationStreamActive
    ) {
      const onMessage = onMessageFunc(dispatch);
      if (notifications.notifications.length > 0) {
        const latestNotification = notifications.notifications[0];
        console.log(
          "listen for notifications",
          latestNotification.notification.redisID,
        );
        getNotificationUpdate(
          latestNotification.notification.redisID,
          onMessage,
          onEnd,
        );
      } else {
        getNotificationUpdate("", onMessage, onEnd);
      }
      setNotificationStreamActive(true);
    }
  }, [
    notificationLoading,
    notificationError,
    notificationStreamActive,
    getNotificationUpdate,
    notifications,
    dispatch,
  ]);

  const providerValues: DashboardConsumerProps = {
    notificationOpen,
    closeNotification,
    openNotification,
    toggleNotification,
    dashboardState: state,
    dashboardDispatch: dispatch,
  };

  return <Provider value={providerValues}>{children}</Provider>;
};
