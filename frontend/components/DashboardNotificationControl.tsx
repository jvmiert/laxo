import { Fragment, useRef } from "react";
import { Popover, Transition } from "@headlessui/react";
import { BellIcon } from "@heroicons/react/solid";
import { useDashboard } from "@/providers/DashboardProvider";
import DashboardNotification from "@/components/DashboardNotification";
import NotificationEmpty from "@/components/NotificationEmpty";
import { XIcon } from "@heroicons/react/outline";

export default function DashboardNotificationControl() {
  const { dashboardState, notificationLoading } = useDashboard();

  const notiRef = useRef<HTMLButtonElement>(null);
  return (
    <Popover className="relative">
      {({ open, close }) => {
        return (
          <>
            <Popover.Button ref={notiRef}>
              <div className="relative mr-4 rounded-full bg-gray-100 p-1">
                <BellIcon className="h-5 w-5 text-gray-900" />
                {!notificationLoading &&
                  dashboardState.notifications.length > 0 && (
                    <span className="text-light absolute top-0 right-0 inline-flex translate-x-1/2 -translate-y-1/2 transform items-center justify-center rounded-full bg-indigo-600 px-1.5 py-1 text-xs leading-none text-indigo-100">
                      {dashboardState.notifications.length}
                    </span>
                  )}
              </div>
            </Popover.Button>
            <Transition
              appear={true}
              show={open}
              as={Fragment}
              enter="transition ease-out duration-200"
              enterFrom="opacity-0 translate-y-1"
              enterTo="opacity-100 translate-y-0"
              leave="transition ease-in duration-150"
              leaveFrom="opacity-100 translate-y-0"
              leaveTo="opacity-0 translate-y-1"
            >
              <Popover.Panel
                static
                className="absolute top-full -left-64 z-10 w-80 rounded bg-gray-50 pb-6"
              >
                <div className="absolute right-11 -z-10 h-4 w-4 origin-top-left -rotate-45 transform rounded-sm bg-gray-50" />
                <div className="absolute right-2 top-2">
                  <button className="p-2" onClick={() => close()}>
                    <XIcon className="h-4 w-4" />
                  </button>
                </div>
                <h2 className="pt-6 pb-4 text-center text-xl font-bold">
                  Notifications
                </h2>
                <div className="flex max-h-[45vh] flex-col items-center gap-y-4 self-center overflow-y-auto overscroll-y-contain">
                  {dashboardState.notifications.length == 0 && (
                    <NotificationEmpty />
                  )}
                  {dashboardState.notifications.map((n) => (
                    <DashboardNotification
                      key={n.notificationGroup.id}
                      notification={n.notification}
                      notificationGroup={n.notificationGroup}
                    />
                  ))}
                </div>
              </Popover.Panel>
            </Transition>
          </>
        );
      }}
    </Popover>
  );
}
