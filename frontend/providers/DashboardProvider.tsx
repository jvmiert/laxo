import { ReactChildren, ReactNode, useState, useEffect, Dispatch } from "react";
import { Draft } from "immer";
import { useImmerReducer } from "use-immer";
import createSafeContext from "@/lib/useSafeContext";
import { useGetNotifications } from "@/hooks/swrHooks";
import type { NotificationResponseObject } from "@/types/ApiResponse";

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
      if (index !== -1) draft.notifications.splice(index, 1);
      draft.notifications.unshift(action.notification);
      break;
  }
}

export const [useDashboard, Provider] =
  createSafeContext<DashboardConsumerProps>();

export const DashboardProvider = ({
  children,
}: {
  children: ReactChildren | ReactNode;
}) => {
  const [notificationOpen, setNotificationOpen] = useState(false);

  const closeNotification = () => setNotificationOpen(false);
  const openNotification = () => setNotificationOpen(true);
  const toggleNotification = () => setNotificationOpen(!notificationOpen);

  const [state, dispatch] = useImmerReducer(reducer, initialState);

  const { notifications } = useGetNotifications();

  useEffect(() => {
    if (notifications.notifications.length > 0) {
      dispatch({
        type: "reset",
        state: { notifications: notifications.notifications },
      });
    }
  }, [notifications, dispatch]);

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
