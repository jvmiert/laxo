import cc from "classcat";
import { Disclosure, Transition, Tab, Switch } from "@headlessui/react";
import { ChevronUpIcon, CubeIcon } from "@heroicons/react/solid";
import { Fragment, useState } from "react";
import { useIntl } from "react-intl";

import LazadaIcon from "@/components/icons/LazadaIcon";
import ShopeeIcon from "@/components/icons/ShopeeIcon";
import { LaxoProductDetails } from "@/types/ApiResponse";
import { LaxoProductPlatforms } from "@/types/ApiResponse";
import { useDashboard } from "@/providers/DashboardProvider";
import useProductApi from "@/hooks/useProductApi";

function getPlatformIcon(platform: string): JSX.Element {
  switch (platform.toLowerCase()) {
    case "lazada":
      return <LazadaIcon className="-ml-0.5 mr-2 h-5 w-5" />;
    case "shopee":
      return <ShopeeIcon className="-ml-0.5 mr-2 h-5 w-5 fill-[#ff5422]" />;
    default:
      return <CubeIcon className="-ml-0.5 mr-2 h-5 w-5" />;
  }
}

type PlatformSyncState = {
  [key: string]: boolean;
};

function getInitialSyncStatus(
  platformList: Array<LaxoProductPlatforms>,
): PlatformSyncState {
  let platformSyncInitialState: PlatformSyncState = {};
  platformList.forEach((p) => {
    if (p.platformName == "lazada")
      platformSyncInitialState[p.platformName] = p.syncStatus;
  });

  return platformSyncInitialState;
}

export type PlatformsEditProps = {
  product: LaxoProductDetails;
};

export default function DetailsPlatformsEdit({ product }: PlatformsEditProps) {
  const { doChangePlatformSync } = useProductApi();
  const { dashboardDispatch } = useDashboard();
  const t = useIntl();

  const [syncState, setSyncState] = useState<PlatformSyncState>(
    getInitialSyncStatus(product.platforms),
  );

  const changedSync = (checked: boolean, platformName: string) => {
    setSyncState((prevState) => ({
      ...prevState,
      ...{ [platformName]: checked },
    }));
    const success = doChangePlatformSync({
      productID: product.product.id,
      platform: platformName,
      state: checked,
    });

    if (!!success) {
      dashboardDispatch({
        type: "alert",
        alert: {
          type: "success",
          message: t.formatMessage({
            description: "Platform management successful changed sync",
            defaultMessage:
              "Successfully changed your platform synchronization",
          }),
        },
      });
    } else {
      dashboardDispatch({
        type: "alert",
        alert: {
          type: "error",
          message: t.formatMessage({
            description: "Platform management changed sync error",
            defaultMessage:
              "Could not changed your synchronization status, please try again later",
          }),
        },
      });
    }
  };
  return (
    <Disclosure defaultOpen>
      {({ open }) => (
        <>
          <Disclosure.Button className="flex w-full justify-between rounded-xl bg-gray-50 px-4 py-3">
            <h3 className="text-lg font-medium leading-6 text-gray-900">
              Platforms
            </h3>
            <ChevronUpIcon
              className={cc(["h-4 h-4", { "rotate-180 transform": open }])}
            />
          </Disclosure.Button>
          <Transition
            show={open}
            enter="transition ease-out duration-200"
            enterFrom="opacity-0 translate-y-0"
            enterTo="opacity-100 translate-y-2"
            leave="transition ease-in duration-150"
            leaveFrom="opacity-100 translate-y-2"
            leaveTo="opacity-0 translate-y-0"
            unmount={false}
          >
            <Disclosure.Panel static className="p-4">
              <div>
                <Tab.Group>
                  <Tab.List className="-mb-px flex w-full space-x-8 border-b border-gray-200">
                    {product.platforms.map((p) => (
                      <Tab key={p.platformName} as={Fragment}>
                        {({ selected }) => (
                          <button
                            className={cc([
                              "group inline-flex items-center border-b-2 py-4 px-1 focus:outline-none",
                              {
                                "border-transparent hover:border-gray-300":
                                  !selected,
                              },
                              {
                                "border-indigo-500": selected,
                              },
                            ])}
                          >
                            {getPlatformIcon(p.platformName)}
                            <span className="capitalize">{p.platformName}</span>
                          </button>
                        )}
                      </Tab>
                    ))}
                  </Tab.List>
                  <Tab.Panels className="py-4">
                    {product.platforms.map((p) => (
                      <Tab.Panel key={p.platformName}>
                        <dl className="grid grid-cols-6 gap-x-4 gap-y-8">
                          <div className="col-span-3">
                            <dt className="text-sm font-medium text-gray-500">
                              Name
                            </dt>
                            <dd className="mt-1 text-sm text-gray-900">
                              {p.name}
                            </dd>
                          </div>
                          <div className="col-span-2">
                            <dt className="text-sm font-medium text-gray-500">
                              Platform SKU
                            </dt>
                            <dd className="mt-1 text-sm text-gray-900">
                              {p.platformSKU}
                            </dd>
                          </div>
                          <div className="col-span-1">
                            <dt className="text-sm font-medium text-gray-500">
                              Status
                            </dt>
                            <dd className="mt-1 text-sm text-gray-900">
                              <span
                                className={cc([
                                  "rounded-md py-1 px-2 capitalize",
                                  {
                                    "bg-green-100":
                                      p.status.toLowerCase() == "active",
                                  },
                                  {
                                    "bg-gray-100":
                                      p.status.toLowerCase() != "active",
                                  },
                                ])}
                              >
                                {p.status.toLowerCase()}
                              </span>
                            </dd>
                          </div>
                          <div className="col-span-6">
                            <dt className="text-sm font-medium text-gray-500">
                              Product Link
                            </dt>
                            <dd className="mt-1 text-sm text-gray-900">
                              <a
                                href={p.productURL}
                                target="_blank"
                                rel="noreferrer"
                              >
                                {p.productURL}
                              </a>
                            </dd>
                          </div>
                          {p.platformName in syncState && (
                            <div className="col-span-6">
                              <dt className="text-sm font-medium text-gray-500">
                                Keep Product Updated
                              </dt>
                              <dd className="mt-2 text-sm text-gray-900">
                                <Switch
                                  checked={syncState[p.platformName]}
                                  onChange={(c) =>
                                    changedSync(c, p.platformName)
                                  }
                                  className={cc([
                                    {
                                      "bg-gray-200": !syncState[p.platformName],
                                    },
                                    {
                                      "bg-indigo-400":
                                        syncState[p.platformName],
                                    },
                                    "relative inline-flex h-6 w-11 flex-shrink-0 cursor-pointer rounded-full border-2 border-transparent transition-colors duration-200 ease-in-out focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2",
                                  ])}
                                >
                                  <span
                                    className={cc([
                                      {
                                        "translate-x-5":
                                          syncState[p.platformName],
                                      },
                                      {
                                        "translate-x-0":
                                          !syncState[p.platformName],
                                      },
                                      "pointer-events-none relative inline-block h-5 w-5 transform rounded-full bg-white shadow ring-0 transition duration-200 ease-in-out",
                                    ])}
                                  >
                                    <span
                                      className={cc([
                                        {
                                          "opacity-0 duration-100 ease-out":
                                            syncState[p.platformName],
                                        },
                                        {
                                          "opacity-100 duration-200 ease-in":
                                            !syncState[p.platformName],
                                        },
                                        "absolute inset-0 flex h-full w-full items-center justify-center transition-opacity",
                                      ])}
                                      aria-hidden="true"
                                    >
                                      <svg
                                        className="h-3 w-3 text-gray-400"
                                        fill="none"
                                        viewBox="0 0 12 12"
                                      >
                                        <path
                                          d="M4 8l2-2m0 0l2-2M6 6L4 4m2 2l2 2"
                                          stroke="currentColor"
                                          strokeWidth={2}
                                          strokeLinecap="round"
                                          strokeLinejoin="round"
                                        />
                                      </svg>
                                    </span>
                                    <span
                                      className={cc([
                                        {
                                          "opacity-0 duration-100 ease-out":
                                            !syncState[p.platformName],
                                        },
                                        {
                                          "opacity-100 duration-200 ease-in":
                                            syncState[p.platformName],
                                        },
                                        "absolute inset-0 flex h-full w-full items-center justify-center transition-opacity",
                                      ])}
                                      aria-hidden="true"
                                    >
                                      <svg
                                        className="h-3 w-3 text-indigo-600"
                                        fill="currentColor"
                                        viewBox="0 0 12 12"
                                      >
                                        <path d="M3.707 5.293a1 1 0 00-1.414 1.414l1.414-1.414zM5 8l-.707.707a1 1 0 001.414 0L5 8zm4.707-3.293a1 1 0 00-1.414-1.414l1.414 1.414zm-7.414 2l2 2 1.414-1.414-2-2-1.414 1.414zm3.414 2l4-4-1.414-1.414-4 4 1.414 1.414z" />
                                      </svg>
                                    </span>
                                  </span>
                                </Switch>
                              </dd>
                            </div>
                          )}
                        </dl>
                      </Tab.Panel>
                    ))}
                  </Tab.Panels>
                </Tab.Group>
              </div>
            </Disclosure.Panel>
          </Transition>
        </>
      )}
    </Disclosure>
  );
}
