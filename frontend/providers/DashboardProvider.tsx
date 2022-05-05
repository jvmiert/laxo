import { ReactChildren, ReactNode, useState } from "react";
import createSafeContext from "@/lib/useSafeContext";

export interface DashboardConsumerProps {
  notificationOpen: boolean;
  closeNotification: () => void;
  openNotification: () => void;
  toggleNotification: () => void;
}

interface DashboardState {}

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

  const providerValues: DashboardConsumerProps = {
    notificationOpen,
    closeNotification,
    openNotification,
    toggleNotification,
  };

  return <Provider value={providerValues}>{children}</Provider>;
};
