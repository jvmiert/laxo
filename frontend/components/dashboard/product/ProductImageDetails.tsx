import { Fragment, useRef } from "react";
import { Dialog, Transition } from "@headlessui/react";
import { XIcon, TrashIcon } from "@heroicons/react/solid";
import Image from "next/image";
import prettyBytes from "pretty-bytes";
import { useIntl } from "react-intl";

import { LaxoProductAsset } from "@/types/ApiResponse";
import { Asset } from "@/hooks/useProductApi";

type ProductImageDetailsProps = {
  show: boolean;
  close: () => void;
  removeAsset: () => void;
  asset: LaxoProductAsset | Asset | undefined;
  assetsToken: string;
};

export default function ProductImageDetails({
  show,
  close,
  asset,
  assetsToken,
  removeAsset,
}: ProductImageDetailsProps) {
  const t = useIntl();

  let downloadButtonRef = useRef<HTMLAnchorElement>(null);

  return (
    <Transition appear show={show} as={Fragment}>
      <Dialog
        initialFocus={downloadButtonRef}
        as="div"
        className="relative z-10"
        onClose={() => close()}
      >
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
              <Dialog.Panel className="h-full w-full transform overflow-hidden transition-all">
                <div className="flex h-full">
                  <div className="grow rounded-l-lg bg-black p-6">
                    <div className="relative h-full">
                      {asset && (
                        <Image
                          alt=""
                          src={`/api/assets/${assetsToken}/${asset.id}${asset.extension}`}
                          layout="fill"
                          objectFit="contain"
                          objectPosition="center"
                        />
                      )}
                    </div>
                  </div>
                  <div className="flex w-96 flex-col rounded-r-lg bg-white p-6">
                    <div className="flex items-start justify-between">
                      <Dialog.Title className="text-lg font-medium text-gray-900">
                        {t.formatMessage({
                          defaultMessage: "Image Details",
                          description: "Product image details: title",
                        })}
                      </Dialog.Title>
                      <div className="ml-3 flex h-7 items-center">
                        <button
                          type="button"
                          className="rounded-md bg-white text-gray-400 hover:text-gray-500 focus:ring-2 focus:ring-indigo-500"
                          onClick={() => close()}
                        >
                          <span className="sr-only">
                            {t.formatMessage({
                              defaultMessage: "Close panel",
                              description:
                                "Product image details: close button",
                            })}
                          </span>
                          <XIcon className="h-6 w-6" aria-hidden="true" />
                        </button>
                      </div>
                    </div>
                    {asset && (
                      <div className="pt-12">
                        <h3 className="font-medium text-gray-900">
                          {t.formatMessage({
                            defaultMessage: "Information",
                            description:
                              "Product image details: asset details title",
                          })}
                        </h3>
                        <dl className="mt-2 divide-y divide-gray-200 border-t border-b border-gray-200">
                          <div className="flex justify-between py-3 text-sm font-medium">
                            <dt className="text-gray-500">
                              {t.formatMessage({
                                defaultMessage: "Name",
                                description:
                                  "Product image details: asset details name label",
                              })}
                            </dt>
                            <dd className="text-gray-900">
                              {asset.originalFilename}
                            </dd>
                          </div>
                          <div className="flex justify-between py-3 text-sm font-medium">
                            <dt className="text-gray-500">
                              {t.formatMessage({
                                defaultMessage: "Dimensions",
                                description:
                                  "Product image details: asset details dimension label",
                              })}
                            </dt>
                            <dd className="text-gray-900">
                              {asset.width} x {asset.height}
                            </dd>
                          </div>
                          <div className="flex justify-between py-3 text-sm font-medium">
                            <dt className="text-gray-500">
                              {t.formatMessage({
                                defaultMessage: "File Size",
                                description:
                                  "Product image details: asset details file size label",
                              })}
                            </dt>
                            <dd className="text-gray-900">
                              {prettyBytes(asset.fileSize, { locale: "vi" })}
                            </dd>
                          </div>
                        </dl>
                      </div>
                    )}
                    <div className="grow" />
                    {asset && (
                      <div className="flex">
                        <a
                          download
                          ref={downloadButtonRef}
                          className="flex-1 rounded-md border border-transparent bg-indigo-500 py-2 px-4 text-center text-sm font-medium text-white shadow-sm hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2"
                          target="_blank"
                          rel="noreferrer"
                          href={`/api/assets/${assetsToken}/${asset.id}${asset.extension}`}
                        >
                          {t.formatMessage({
                            defaultMessage: "Download",
                            description:
                              "Product image details: download button",
                          })}
                        </a>
                        <button
                          type="button"
                          className="ml-3 inline-flex shrink grow basis-0 items-center justify-center rounded-md border border-gray-300 bg-white py-2 px-4 text-sm font-medium leading-4 text-gray-700 shadow-sm hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2"
                          onClick={removeAsset}
                        >
                          <TrashIcon
                            className="-ml-0.5 mr-2 h-4 w-4"
                            aria-hidden="true"
                          />
                          {t.formatMessage({
                            defaultMessage: "Delete",
                            description: "Product image details: delete button",
                          })}
                        </button>
                      </div>
                    )}
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
