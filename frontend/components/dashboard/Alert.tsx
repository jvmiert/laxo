import cc from "classcat";
import { Fragment, useState, useEffect } from "react";
import { Transition } from "@headlessui/react";
import {
  CheckCircleIcon,
  ExclamationIcon,
  XCircleIcon,
  XIcon,
} from "@heroicons/react/solid";
import { useIntl } from "react-intl";

import { useDashboard } from "@/providers/DashboardProvider";
import { Alert as AlertProp } from "@/providers/DashboardProvider";

export default function Alert({ id, type, message }: AlertProp) {
  const t = useIntl();

  let [isShowing, setIsShowing] = useState(true);
  const { dashboardDispatch } = useDashboard();

  useEffect(() => {
    let timeout: ReturnType<typeof setTimeout>;
    if (isShowing) {
      timeout = setTimeout(() => {
        setIsShowing(false);
      }, 5000);
    }

    return () => {
      if (timeout) {
        clearTimeout(timeout);
      }
    };
  }, [isShowing, dashboardDispatch, id]);

  // We do this in order to give the Transition animation
  // time to play out.
  useEffect(() => {
    let timeout: ReturnType<typeof setTimeout>;
    if (!isShowing) {
      timeout = setTimeout(() => {
        dashboardDispatch({
          type: "remove_alert",
          id: id,
        });
      }, 1000);
    }

    return () => {
      if (timeout) {
        clearTimeout(timeout);
      }
    };
  }, [id, isShowing, dashboardDispatch]);

  return (
    <Transition
      appear={true}
      as={Fragment}
      show={isShowing}
      enter="transition ease-out duration-200"
      enterFrom="opacity-0 translate-y-full"
      enterTo="opacity-100 translate-y-0"
      leave="transition ease-in duration-150"
      leaveFrom="opacity-100 translate-y-0"
      leaveTo="opacity-0 translate-y-full"
    >
      <div
        className={cc([
          "rounded-md p-4 shadow-md",
          { "bg-green-50": type === "success" },
          { "bg-yellow-50": type === "warning" },
          { "bg-red-50": type === "error" },
        ])}
      >
        <div className="flex">
          <div className="flex-shrink-0">
            {type === "success" && (
              <CheckCircleIcon
                className="h-5 w-5 text-green-400"
                aria-hidden="true"
              />
            )}
            {type === "warning" && (
              <ExclamationIcon
                className="h-5 w-5 text-yellow-400"
                aria-hidden="true"
              />
            )}
            {type === "error" && (
              <XCircleIcon
                className="h-5 w-5 text-red-400"
                aria-hidden="true"
              />
            )}
          </div>
          <div className="ml-3">
            <p
              className={cc([
                "text-sm font-medium",
                { "text-green-800": type === "success" },
                { "text-yellow-800": type === "warning" },
                { "text-red-800": type === "error" },
              ])}
            >
              {message}
            </p>
          </div>
          <div className="ml-auto pl-3">
            <div className="-mx-1.5 -my-1.5">
              <button
                type="button"
                onClick={() => setIsShowing(false)}
                className={cc([
                  "inline-flex rounded-md p-1.5 focus:outline-none focus:ring-2 focus:ring-offset-2 ",
                  {
                    "bg-green-50 text-green-500 hover:bg-green-100 focus:ring-green-600 focus:ring-offset-green-50":
                      type === "success",
                  },
                  {
                    "bg-yellow-50 text-yellow-500 hover:bg-yellow-100 focus:ring-yellow-600 focus:ring-offset-yellow-50":
                      type === "warning",
                  },
                  {
                    "bg-red-50 text-red-500 hover:bg-red-100 focus:ring-red-600 focus:ring-offset-red-50":
                      type === "error",
                  },
                ])}
              >
                <span className="sr-only">
                  {t.formatMessage({
                    defaultMessage: "Dismiss",
                    description: "Alert: dismiss button",
                  })}
                </span>
                <XIcon className="h-5 w-5" aria-hidden="true" />
              </button>
            </div>
          </div>
        </div>
      </div>
    </Transition>
  );
}
