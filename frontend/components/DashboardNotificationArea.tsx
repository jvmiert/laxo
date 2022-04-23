import { Transition } from "@headlessui/react";
import { useDashboard } from "@/providers/DashboardProvider";
import DashboardLoadingEvent from "@/components/DashboardLoadingEvent";

export default function DashboardNotificationArea() {
  const { notificationOpen } = useDashboard();

  return (
    <Transition
      appear={true}
      show={notificationOpen}
      enter="transition duration-75 ease-out"
      enterFrom="transform origin-top-right scale-0 -translate-y-2"
      enterTo="transform origin-top-right scale-100 translate-y-0"
      leave="transition duration-75 ease-out"
      leaveFrom="transform origin-top-right scale-100 translate-y-0"
      leaveTo="transform origin-top-right scale-0 -translate-y-2"
    >
      <div className="relative ml-6 w-80 shrink-0 rounded bg-gray-100 pb-6">
        <div className="absolute right-6 -z-10 h-4 w-4 origin-top-left -rotate-45 transform rounded-sm bg-gray-100" />
        <h2 className="py-6 text-center text-xl font-bold">Notifications</h2>
        <div className="flex flex-col items-center gap-y-4 self-center">
          <DashboardLoadingEvent />
          <DashboardLoadingEvent />
          <DashboardLoadingEvent />
        </div>
      </div>
    </Transition>
  );
}
