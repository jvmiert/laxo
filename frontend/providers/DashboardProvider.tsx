import {
  ReactNode,
  useState,
  useEffect,
  Dispatch,
  useRef,
  useCallback,
  useMemo,
  MutableRefObject,
} from "react";
import { nanoid } from "nanoid";
import { Draft } from "immer";
import { useImmerReducer } from "use-immer";
import { grpc } from "@improbable-eng/grpc-web";
import createSafeContext from "@/lib/useSafeContext";
import { useGetNotifications, useGetShop } from "@/hooks/swrHooks";
import useNotificationApi from "@/hooks/useNotificationApi";
import type {
  NotificationResponseObject,
  GetShopResponseShops,
} from "@/types/ApiResponse";
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
  activeShop: GetShopResponseShops | null;
  slateResetRef: MutableRefObject<() => void>;
  toggleSlateDirtyState: () => void;
  slateIsDirty: boolean;
  productDetailFormResetRef: MutableRefObject<() => void>;
  toggleProductDetailFormDirtyState: () => void;
  productDetailFormIsDirty: boolean;
  productDetailSubmitIsDisabled: boolean;
  toggleProductDetailSubmitIsDisabled: () => void;
}

export type Alert = {
  id: string;
  type: "success" | "warning" | "error";
  message: string;
};

type AlertAction = Omit<Alert, "id">;

interface DashboardState {
  notifications: Array<NotificationResponseObject>;
  alerts: Array<Alert>;
}

export const InitialDashboardState: DashboardState = {
  notifications: [],
  alerts: [],
};

export type DashboardAction =
  | { type: "reset"; state: DashboardState }
  | { type: "add"; notification: NotificationResponseObject }
  | { type: "remove_alert"; id: string }
  | { type: "alert"; alert: AlertAction };

function reducer(draft: Draft<DashboardState>, action: DashboardAction) {
  switch (action.type) {
    case "reset":
      return action.state;
    case "alert":
      draft.alerts.push({ ...action.alert, id: nanoid() });
      break;
    case "remove_alert":
      const alertIndex = draft.alerts.findIndex((a) => a.id === action.id);
      if (alertIndex !== -1) draft.alerts.splice(alertIndex, 1);
      break;
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

  const productDetailFormResetRef = useRef<() => void>(() => {});
  const [productDetailFormIsDirty, setProductDetailFormIsDirty] =
    useState(false);
  const toggleProductDetailFormDirtyState = useCallback(() => {
    setProductDetailFormIsDirty((prevState) => !prevState);
  }, []);

  const [productDetailSubmitIsDisabled, setProductDetailSubmitIsDisabled] =
    useState(false);
  const toggleProductDetailSubmitIsDisabled = useCallback(() => {
    setProductDetailSubmitIsDisabled((prevState) => !prevState);
  }, []);

  const slateResetRef = useRef<() => void>(() => {});
  const [slateIsDirty, setSlateIsDirty] = useState(false);
  const toggleSlateDirtyState = useCallback(() => {
    setSlateIsDirty((prevState) => !prevState);
  }, []);

  const [state, dispatch] = useImmerReducer(reducer, InitialDashboardState);

  const { getNotificationUpdate } = useNotificationApi();

  const { shops } = useGetShop();
  const { auth } = useAuth();
  const { route } = useRouter();

  const activeShop = shops.total > 0 ? shops.shops[0] : null;

  const {
    notifications,
    loading: notificationLoading,
    error: notificationError,
  } = useGetNotifications();

  useEffect(() => {
    if (notifications.notifications.length > 0) {
      dispatch({
        type: "reset",
        state: { notifications: notifications.notifications, alerts: [] },
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
      activeShop: activeShop,
      slateResetRef: slateResetRef,
      toggleSlateDirtyState: toggleSlateDirtyState,
      slateIsDirty: slateIsDirty,
      productDetailFormResetRef: productDetailFormResetRef,
      toggleProductDetailFormDirtyState: toggleProductDetailFormDirtyState,
      productDetailFormIsDirty: productDetailFormIsDirty,
      productDetailSubmitIsDisabled: productDetailSubmitIsDisabled,
      toggleProductDetailSubmitIsDisabled: toggleProductDetailSubmitIsDisabled,
    }),
    [
      notificationOpen,
      closeNotification,
      openNotification,
      toggleNotification,
      state,
      dispatch,
      notificationLoading,
      activeShop,
      toggleSlateDirtyState,
      slateIsDirty,
      productDetailFormIsDirty,
      toggleProductDetailFormDirtyState,
      productDetailSubmitIsDisabled,
      toggleProductDetailSubmitIsDisabled,
    ],
  );

  return <Provider value={providerValues}>{children}</Provider>;
};
