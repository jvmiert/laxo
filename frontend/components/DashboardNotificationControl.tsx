import { BellIcon } from "@heroicons/react/solid";

export default function DashboardNotificationControl() {
  return (
    <button type="button">
      <div className="relative rounded-full bg-gray-100 p-1 mr-4">
        <BellIcon className="h-5 w-5 text-gray-900" />
        <span className="text-light absolute top-0 right-0 inline-flex translate-x-1/2 -translate-y-1/2 transform items-center justify-center rounded-full bg-indigo-600 px-1.5 py-1 text-xs leading-none text-indigo-100">
          99
        </span>
      </div>
    </button>
  );
}
