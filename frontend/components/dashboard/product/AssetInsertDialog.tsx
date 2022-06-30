import { Fragment } from "react";
import { Dialog, Transition } from "@headlessui/react";
import { XIcon, PlusCircleIcon, SearchIcon } from "@heroicons/react/solid";

import { useDashboard } from "@/providers/DashboardProvider";
import { useGetShopAssets } from "@/hooks/swrHooks";

type ProductImageDetailsProps = {};

export default function AssetInsertDialog({}: ProductImageDetailsProps) {
  const {
    dashboardDispatch,
    dashboardState: { insertImageIsOpen },
  } = useDashboard();

  const { assets, loading } = useGetShopAssets(1, 10);

  const closeDialog = () => {
    dashboardDispatch({
      type: "close_image_insert",
    });
  };

  return (
    <Transition appear show={!!insertImageIsOpen} as={Fragment}>
      <Dialog as="div" className="relative z-10" onClose={closeDialog}>
        <Transition.Child
          as={Fragment}
          enter="ease-out duration-300"
          enterFrom="opacity-0"
          enterTo="opacity-100"
          leave="ease-in duration-200"
          leaveFrom="opacity-100"
          leaveTo="opacity-0"
        >
          <div className="fixed inset-0 h-screen w-screen bg-zinc-800 bg-opacity-80" />
        </Transition.Child>
        <div className="fixed inset-0 overflow-y-auto">
          <div className="flex h-full items-center justify-center p-14">
            <Transition.Child
              as={Fragment}
              enter="ease-out duration-300"
              enterFrom="opacity-0 scale-95"
              enterTo="opacity-100 scale-100"
              leave="ease-in duration-200"
              leaveFrom="opacity-100 scale-100"
              leaveTo="opacity-0 scale-95"
            >
              <Dialog.Panel className="h-full w-full transform overflow-hidden rounded-lg bg-white transition-all">
                <div className="px-4 py-5 sm:px-6">
                  <div className="flex items-center justify-between">
                    <div>
                      <Dialog.Title className="text-lg font-medium leading-6 text-gray-900">
                        Shop Assets
                      </Dialog.Title>
                      <p className="mt-1 text-sm text-gray-500">
                        Insert one of the images below into your description
                      </p>
                    </div>
                    <div className="flex items-center space-x-5">
                      <div className="relative rounded-md border shadow">
                        <div className="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-3">
                          <span className="text-gray-500">
                            <SearchIcon className="h-4 w-4" />
                          </span>
                        </div>
                        <input
                          onChange={() => {}}
                          type="text"
                          className="block w-full rounded-md py-2 pl-9 pr-9 focus:outline-none focus:ring focus:ring-indigo-200"
                          placeholder="Search for image name"
                          defaultValue={""}
                        />
                      </div>
                      <button
                        type="button"
                        className="inline-flex items-center rounded-md border border-indigo-500 bg-indigo-500 py-2 px-4 text-white shadow shadow-indigo-500/50 hover:bg-indigo-700 focus:outline-none focus:ring focus:ring-indigo-200 disabled:cursor-not-allowed disabled:bg-indigo-200"
                        onClick={() => {}}
                      >
                        <PlusCircleIcon
                          className="mr-2 -ml-1 h-4 w-4"
                          aria-hidden="true"
                        />
                        Add Image
                      </button>
                      <button
                        type="button"
                        className="rounded-md bg-white text-gray-400 hover:text-gray-500 focus:ring-2 focus:ring-indigo-500"
                        onClick={closeDialog}
                      >
                        <span className="sr-only">Close panel</span>
                        <XIcon className="h-6 w-6" aria-hidden="true" />
                      </button>
                    </div>
                  </div>
                  <div className="my-6 -ml-4 -mr-4 border-b border-gray-200 sm:-ml-6 sm:-mr-6" />
                  <div className="">
                    <ul
                      role="list"
                      className="grid grid-cols-4 gap-x-4 gap-y-8"
                    >
                      {assets.assets.map((a) => (
                        <li key={a.id} className="relative">
                          {a.originalFilename}
                        </li>
                      ))}
                    </ul>
                  </div>
                </div>
              </Dialog.Panel>
            </Transition.Child>
          </div>
        </div>
      </Dialog>
    </Transition>
  );
}
