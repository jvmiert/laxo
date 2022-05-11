import { useCallback } from "react";
import cc from "classcat";
import { defineMessage, useIntl, MessageDescriptor } from "react-intl";
import { CollectionIcon, CheckIcon } from "@heroicons/react/outline";
import type { Notification, NotificationGroup } from "@/types/ApiResponse";

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
        defaultMessage: "",
      });
  }
};

const notificationMainTranslate = function (key: string): MessageDescriptor {
  switch (key) {
    case "product_add":
      return defineMessage({
        description: "notification message: product sync",
        defaultMessage: "Synchronizing your products",
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

  const getProgress = useCallback((): number => {
    const currentMainStep = notification.currentMainStep;
    const currentSubStep = notification.currentSubStep;

    if (currentMainStep == 0 && currentSubStep == 0) return 10;

    const totalMainSteps = notificationGroup.totalMainSteps;
    const totalSubSteps = notificationGroup.totalSubSteps;

    const mainPercentagePerStep = 50 / totalMainSteps;
    const mainPercentage = currentMainStep * mainPercentagePerStep;

    if (totalSubSteps == undefined || currentSubStep == undefined) return mainPercentage;

    if (totalSubSteps == 0) return mainPercentage;

    const subPercentagePerStep = 50 / totalSubSteps;
    const subPercentage = currentSubStep * subPercentagePerStep;

    return mainPercentage + subPercentage;
  }, [notification, notificationGroup]);

  const done = getProgress() === 100;

  return (
    <div className="flex min-w-full flex-row rounded px-4">
      <div className="pr-4">
        <div className="rounded-full bg-indigo-100 p-4">
          {done ? (
            <CheckIcon className="h-4 w-4 text-indigo-500" />
          ) : (
            <CollectionIcon className="h-4 w-4 text-indigo-500" />
          )}
        </div>
      </div>
      <div className="max-w-[400px] flex-grow space-y-2">
        <p className="text-sm font-medium">
          {t.formatMessage(
            notificationMainTranslate(notificationGroup.entityType),
          )}
        </p>
        <div className="relative py-1">
          <div className="flex h-2 overflow-hidden rounded bg-indigo-200 shadow-md shadow-indigo-500/20">
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
            <p className="mt-2 text-xs">
              {t.formatMessage(
                notificationSubTranslate(notification.mainMessage),
              )}
              ...
            </p>
            {!!notification.currentSubStep && notification.currentSubStep > 0 && (
              <p className="mt-2 text-right text-xs">
                {notification.currentSubStep} /{" "}
                {notificationGroup.totalSubSteps}
              </p>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}
