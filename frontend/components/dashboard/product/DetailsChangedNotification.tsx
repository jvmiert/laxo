import { Fragment } from "react";
import { Transition } from "@headlessui/react";
import { SaveIcon, TrashIcon } from "@heroicons/react/outline";
import { useIntl } from "react-intl";

import { useDashboard } from "@/providers/DashboardProvider";
import LoadSpinner from "@/components/LoadSpinner";

export type DetailsChangedNotificationProps = {};

export default function DetailsChangedNotification({}: DetailsChangedNotificationProps) {
  const t = useIntl();
  const {
    slateResetRef,
    slateIsDirty,
    toggleSlateDirtyState,
    productDetailFormResetRef,
    productDetailFormIsDirty,
    toggleProductDetailFormDirtyState,
    productDetailSubmitIsDisabled,
    productDetailFormIsSubmitting,
  } = useDashboard();

  const resetFunc = () => {
    productDetailFormResetRef.current();
    if (productDetailFormIsDirty) {
      toggleProductDetailFormDirtyState();
    }

    slateResetRef.current();
    if (slateIsDirty) {
      toggleSlateDirtyState();
    }
  };

  const show = slateIsDirty || productDetailFormIsDirty;

  return (
    <Transition
      appear={true}
      as={Fragment}
      show={show}
      enter="transition ease-out duration-150"
      enterFrom="opacity-0 -translate-y-full"
      enterTo="opacity-100 translate-y-0"
      leave="transition ease-in duration-100"
      leaveFrom="opacity-100 translate-y-0"
      leaveTo="opacity-0 -translate-y-full"
    >
      <div className="fixed top-0 left-0 z-30 w-full bg-indigo-400 shadow">
        <div className="mx-auto flex h-[72px] max-w-lg items-center space-x-8">
          <div className="font-semibold text-white">
            {t.formatMessage({
              defaultMessage: "You have unsaved changes",
              description: "Save notification: unsaved changes",
            })}
          </div>
          <button
            form="generalEditForm"
            type="submit"
            disabled={productDetailSubmitIsDisabled}
            className="ml-3 inline-flex w-28 shrink items-center justify-center rounded-md border border-gray-300 bg-white py-2 px-4 text-sm font-medium leading-4 text-gray-700 shadow-sm hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-indigo-300 focus:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-75"
          >
            {!productDetailFormIsSubmitting ? (
              <>
                <SaveIcon className="-ml-0.5 mr-2 h-4 w-4" aria-hidden="true" />
                {t.formatMessage({
                  defaultMessage: "Save",
                  description: "Save notification: save button",
                })}
              </>
            ) : (
              <LoadSpinner className="h-4 w-4 animate-spin fill-indigo-600 text-gray-200" />
            )}
          </button>
          <button
            onClick={resetFunc}
            type="reset"
            className="ml-3 inline-flex shrink items-center justify-center rounded-md border border-gray-300 bg-white py-2 px-4 text-sm font-medium leading-4 text-gray-700 shadow-sm hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-indigo-300 focus:ring-offset-2"
          >
            <TrashIcon className="-ml-0.5 mr-2 h-4 w-4" aria-hidden="true" />
            {t.formatMessage({
              defaultMessage: "Reset",
              description: "Save notification: reset button",
            })}
          </button>
        </div>
      </div>
    </Transition>
  );
}
