import {
  ReactNode,
  useState,
  useEffect,
  Dispatch,
  useRef,
  useCallback,
  useMemo,
} from "react";
import { Draft } from "immer";
import { useImmerReducer } from "use-immer";
import { grpc } from "@improbable-eng/grpc-web";
import createSafeContext from "@/lib/useSafeContext";
import { useGetNotifications } from "@/hooks/swrHooks";
import useNotificationApi from "@/hooks/useNotificationApi";
import type { NotificationResponseObject } from "@/types/ApiResponse";
import { NotificationUpdateReply } from "@/proto/user_pb";
import { useAuth } from "@/providers/AuthProvider";
import { useRouter } from "next/router";

export interface DashboardConsumerProps {
  notificationOpen: boolean;
  closeNotification: () => void;
  openNotification: () => void;
  toggleNotification: () => void;
  dashboardState: DashboardState;
  dashboardDispatch: Dispatch<DashboardAction>;
  notificationLoading: boolean;
}

interface DashboardState {
  notifications: Array<NotificationResponseObject>;
}

export const InitialDashboardState: DashboardState = {
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
        read: notification.hasRead()
          ? new Date(notiObject.read * 1000)
          : undefined,
        currentMainStep: notiObject.currentmainstep,
        currentSubStep: notification.hasCurrentsubstep()
          ? notiObject.currentsubstep
          : undefined,
        mainMessage: notiObject.mainmessage,
        subMessage: notiObject.submessage,
        error: notiObject.error,
      },
      notificationGroup: {
        id: notiGroupObject.id,
        userID: notiGroupObject.userid,
        workflowID: notiGroupObject.workflowid,
        entityID: notiGroupObject.entityid,
        entityType: notiGroupObject.entitytype,
        platformName: notiGroupObject.platformname,
        totalMainSteps: notiGroupObject.totalmainsteps,
        totalSubSteps: notificationGroup.hasTotalsubsteps()
          ? notiGroupObject.totalsubsteps
          : undefined,
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

  const closeNotification = useCallback(() => setNotificationOpen(false), []);
  const openNotification = useCallback(() => setNotificationOpen(true), []);
  const toggleNotification = useCallback(
    () => setNotificationOpen(!notificationOpen),
    [notificationOpen],
  );

  const notificationListenRef = useRef(false);
  const notificationCleanupRef = useRef<grpc.Request | undefined>(undefined);

  const [state, dispatch] = useImmerReducer(reducer, InitialDashboardState);

  const { getNotificationUpdate } = useNotificationApi();

  const { auth } = useAuth();
  const { route } = useRouter();

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
    const inDashboard = route.includes("dashboard");
    if (
      !notificationLoading &&
      !notificationError &&
      !notificationListenRef.current &&
      auth &&
      inDashboard
    ) {
      const onMessage = onMessageFunc(dispatch);
      if (notifications.notifications.length > 0) {
        const latestNotification = notifications.notifications[0];
        notificationCleanupRef.current = getNotificationUpdate(
          latestNotification.notification.redisID,
          onMessage,
          onEnd,
        );
      } else {
        notificationCleanupRef.current = getNotificationUpdate(
          "",
          onMessage,
          onEnd,
        );
      }
      notificationListenRef.current = true;
    }

    if (!inDashboard && notificationCleanupRef.current) {
      notificationCleanupRef.current.close();
      notificationListenRef.current = false;
      notificationCleanupRef.current = undefined;
    }
  }, [
    notificationLoading,
    notificationError,
    getNotificationUpdate,
    notifications,
    dispatch,
    auth,
    route,
  ]);

  const providerValues: DashboardConsumerProps = useMemo(
    () => ({
      notificationOpen,
      closeNotification,
      openNotification,
      toggleNotification,
      dashboardState: state,
      dashboardDispatch: dispatch,
      notificationLoading,
    }),
    [
      notificationOpen,
      closeNotification,
      openNotification,
      toggleNotification,
      state,
      dispatch,
      notificationLoading,
    ],
  );

  return <Provider value={providerValues}>{children}</Provider>;
};
