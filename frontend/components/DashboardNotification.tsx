import { useCallback } from "react";
import cc from "classcat";
import { defineMessage, useIntl, MessageDescriptor } from "react-intl";
import {
  CollectionIcon,
  CheckIcon,
  ExclamationIcon,
} from "@heroicons/react/outline";
import type { Notification, NotificationGroup } from "@/types/ApiResponse";
import { formatDistance } from "date-fns";
import { useRouter } from "next/router";
import { enUS, vi } from "date-fns/locale";

const notificationSubTranslate = function (key: string): MessageDescriptor {
  switch (key) {
    case "save":
      return defineMessage({
        description: "notification message: saving",
        defaultMessage: "Saving your product",
      });

    case "fetch":
      return defineMessage({
        description: "notification message: fetching",
        defaultMessage: "Retrieving your products",
      });

    default:
      return defineMessage({
        id: "notification_sub_message_default",
        description: "notification sub message: default",
        defaultMessage: "...",
      });
  }
};

const notificationCompleteTranslate = function (
  key: string,
): MessageDescriptor {
  switch (key) {
    case "product_add":
      return defineMessage({
        description: "notification complete message: product sync",
        defaultMessage: "Synchronized your products",
      });

    default:
      return defineMessage({
        description: "notification complete message: default",
        defaultMessage: "Success",
      });
  }
};

const notificationErrorTranslate = function (key: string): MessageDescriptor {
  switch (key) {
    case "product_add":
      return defineMessage({
        description: "notification error message: product sync",
        defaultMessage: "Could not synchronize",
      });

    default:
      return defineMessage({
        description: "notification error message: default",
        defaultMessage: "Something went wrong...",
      });
  }
};

const notificationMainTranslate = function (key: string): MessageDescriptor {
  switch (key) {
    case "product_add":
      return defineMessage({
        description: "notification message: product sync",
        defaultMessage: "Synchronizing your products...",
      });

    default:
      return defineMessage({
        description: "notification main message: default",
        defaultMessage: "Notification",
      });
  }
};

export interface DashboardNotificationProps {
  notification: Notification;
  notificationGroup: NotificationGroup;
}

export default function DashboardNotification({
  notification,
  notificationGroup,
}: DashboardNotificationProps) {
  const t = useIntl();
  const { locale } = useRouter();

  const done = notification.mainMessage == "complete" || notification.error;

  const getTitleMessage = () => {
    if (notification.mainMessage == "complete") {
      return t.formatMessage(
        notificationCompleteTranslate(notificationGroup.entityType),
      );
    }

    if (notification.error) {
      return t.formatMessage(
        notificationErrorTranslate(notificationGroup.entityType),
      );
    }

    return t.formatMessage(
      notificationMainTranslate(notificationGroup.entityType),
    );
  };

  const getIcon = () => {
    if (notification.mainMessage == "complete") {
      return <CheckIcon className="h-4 w-4 text-indigo-500" />;
    }

    if (notification.error) {
      return <ExclamationIcon className="h-4 w-4 text-indigo-500" />;
    }

    return <CollectionIcon className="h-4 w-4 text-indigo-500" />;
  };

  const showSubCount =
    !!notification.currentSubStep && notification.currentSubStep;

  const getProgress = useCallback((): number => {
    if (notification.mainMessage == "complete") return 100;

    const currentMainStep = notification.currentMainStep;
    const currentSubStep = notification.currentSubStep;

    if (currentMainStep == 0 && currentSubStep == 0) return 10;

    const totalMainSteps = notificationGroup.totalMainSteps;
    const totalSubSteps = notificationGroup.totalSubSteps;

    const mainPercentagePerStep = 50 / totalMainSteps;
    const mainPercentage = currentMainStep * mainPercentagePerStep;

    if (totalSubSteps == undefined || currentSubStep == undefined)
      return mainPercentage;

    if (totalSubSteps == 0) return mainPercentage;

    const subPercentagePerStep = 50 / totalSubSteps;
    const subPercentage = currentSubStep * subPercentagePerStep;

    return mainPercentage + subPercentage;
  }, [notification, notificationGroup]);

  return (
    <div className="flex min-w-full flex-row rounded px-4">
      <div className="flex items-center pr-4">
        <div className="rounded-full bg-indigo-100 p-4">{getIcon()}</div>
      </div>
      <div className="max-w-[400px] flex-grow space-y-2">
        <div
          className="relative"
          style={{
            transform: `translate(0, ${done ? 20 : 0}px)`,
            transitionDelay: "0.40s",
            transition: "transform 0.25s",
            transitionTimingFunction: "ease-in",
            transformOrigin: "top",
          }}
        >
          <p className="text-sm font-medium">{getTitleMessage()}</p>
          <div
            className="absolute left-0 top-full"
            style={{
              transform: `scaleY(${done ? 1.0 : 0})`,
              transitionDelay: "0.75s",
              transition: "transform 0.25s",
              transitionTimingFunction: "ease-in",
              transformOrigin: "bottom",
            }}
          >
            <span className="text-xs font-light">
              {formatDistance(notification.created, new Date(), {
                addSuffix: true,
                locale: locale == "vi" ? vi : enUS,
              })}
            </span>
          </div>
        </div>
        <div className="relative py-1">
          <div
            className="flex h-2 overflow-hidden rounded bg-indigo-200 shadow-md shadow-indigo-500/20"
            style={{
              transform: `scaleY(${done ? 0 : 1.0})`,
              transition: "transform 0.15s",
              transitionTimingFunction: "ease-in",
              transformOrigin: "top",
            }}
          >
            <div
              style={{
                width: `${getProgress()}%`,
                transitionTimingFunction: "ease-out",
                transition: "width 2s",
              }}
              className={cc([
                "flex",
                { "animate-pulse": !done },
                { "bg-gradient-to-l from-indigo-400 to-indigo-600": !done },
                "bg-indigo-600",
              ])}
            ></div>
          </div>
          <div className="flex place-content-between">
            <p
              className="mt-2 text-xs"
              style={{
                transform: `scaleY(${done ? 0 : 1.0})`,
                transition: "transform 0.15s",
                transitionTimingFunction: "ease-in",
                transformOrigin: "top",
              }}
            >
              {`${t.formatMessage(
                notificationSubTranslate(notification.mainMessage),
              )}...`}
            </p>
            <p
              className="mt-2 text-right text-xs"
              style={{
                transform: `scaleY(${showSubCount ? 1.0 : 0})`,
                transition: "transform 0.25s",
                transitionTimingFunction: "ease-in",
                transformOrigin: "top",
              }}
            >
              {notification.currentSubStep} / {notificationGroup.totalSubSteps}
            </p>
          </div>
        </div>
      </div>
    </div>
  );
}
