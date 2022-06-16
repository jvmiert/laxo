import { Fragment } from "react";
import type { ReactElement } from "react";
import { defineMessage } from "react-intl";
import { withRedirectUnauth, withAuthPage } from "@/lib/withAuth";
import { InferGetServerSidePropsType, GetServerSideProps } from "next";
import DashboardLayout from "@/components/DashboardLayout";
import AssetManagement from "@/components/dashboard/product/AssetManagement";
import { useRouter } from "next/router";
import { useGetLaxoProductDetails } from "@/hooks/swrHooks";
import ErrorPage from "next/error";
import { ChevronUpIcon } from "@heroicons/react/solid";
import { Disclosure, Transition, Tab } from "@headlessui/react";
import cc from "classcat";
import { CubeIcon } from "@heroicons/react/solid";
import LazadaIcon from "@/components/icons/LazadaIcon";
import ShopeeIcon from "@/components/icons/ShopeeIcon";
import Editor from "@/components/slate/Editor";

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

export const getServerSideProps: GetServerSideProps = withRedirectUnauth();

type DashboardProductDetailsProps = InferGetServerSidePropsType<
  typeof getServerSideProps
>;

function DashboardProductDetails(props: DashboardProductDetailsProps) {
  const {
    query: { productID },
  } = useRouter();

  const { product, error } = useGetLaxoProductDetails(productID);

  if (!product?.product) return <></>;

  const notFound = error?.response?.status === 404;

  if (notFound) return <ErrorPage statusCode={404} />;
  if (!productID) return <></>;

  const costPrice =
    parseFloat(
      `${product.product.costPrice.Int}e${product.product.costPrice.Exp}`,
    ) || 0;

  const costString = costPrice.toLocaleString("vi-VN");

  return (
    <div className="mx-auto max-w-5xl">
      <div className="space-y-3">
        <div className="rounded-md bg-white py-4 px-3 shadow-sm">
          <Disclosure defaultOpen>
            {({ open }) => (
              <>
                <Disclosure.Button className="flex w-full justify-between rounded-xl bg-gray-50 px-4 py-3">
                  <h3 className="text-lg font-medium leading-6 text-gray-900">
                    General
                  </h3>
                  <ChevronUpIcon
                    className={cc([
                      "h-4 h-4",
                      { "rotate-180 transform": open },
                    ])}
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
                >
                  <Disclosure.Panel static className="p-4">
                    <div className="grid grid-cols-8 gap-4">
                      <div className="col-span-5">
                        <label
                          className="mb-1 block text-sm text-gray-700"
                          htmlFor="name"
                        >
                          Name
                        </label>
                        <input
                          className="focus:shadow-outline w-full appearance-none rounded border py-2 px-3 leading-tight text-gray-700 shadow focus:outline-none focus:ring focus:ring-indigo-200"
                          type="text"
                          defaultValue={product.product.name}
                        />
                      </div>
                      <div className="col-span-3">
                        <label
                          className="mb-1 block text-sm text-gray-700"
                          htmlFor="msku"
                        >
                          SKU
                        </label>
                        <input
                          className="focus:shadow-outline w-full appearance-none rounded border py-2 px-3 leading-tight text-gray-700 shadow focus:outline-none focus:ring focus:ring-indigo-200"
                          type="text"
                          defaultValue={product.product.msku}
                        />
                      </div>
                      <div className="col-span-8">
                        <label
                          className="mb-1 block text-sm text-gray-700"
                          htmlFor="description"
                        >
                          Description
                        </label>
                        <textarea
                          name="description"
                          rows={8}
                          defaultValue={product.product.description}
                          className="focus:shadow-outline block w-full appearance-none rounded border py-2 px-3 leading-tight text-gray-700 shadow focus:outline-none focus:ring focus:ring-indigo-200"
                        />
                      </div>
                      <div className="col-span-8">
                        <Editor />
                      </div>
                    </div>
                  </Disclosure.Panel>
                </Transition>
              </>
            )}
          </Disclosure>
        </div>

        <div className="rounded-md bg-white py-4 px-3 shadow-sm">
          <Disclosure defaultOpen>
            {({ open }) => (
              <>
                <Disclosure.Button className="flex w-full justify-between rounded-xl bg-gray-50 px-4 py-3">
                  <h3 className="text-lg font-medium leading-6 text-gray-900">
                    Price
                  </h3>
                  <ChevronUpIcon
                    className={cc([
                      "h-4 h-4",
                      { "rotate-180 transform": open },
                    ])}
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
                >
                  <Disclosure.Panel static className="p-4">
                    <div className="grid grid-cols-8 gap-4">
                      <div className="col-start-1 col-end-3">
                        <label
                          className="mb-1 block text-sm text-gray-700"
                          htmlFor="name"
                        >
                          Selling Price
                        </label>
                        <div className="flex rounded shadow">
                          <input
                            className="focus:shadow-outline z-10 block w-full w-full flex-1 appearance-none rounded-none rounded-l border py-2 px-3 leading-tight text-gray-700 focus:outline-none focus:ring focus:ring-indigo-200"
                            type="text"
                            defaultValue={parseFloat(
                              `${product.product.sellingPrice.Int}e${product.product.sellingPrice.Exp}`,
                            ).toLocaleString("vi-VN")}
                          />
                          <span className="inline-flex items-center rounded-r border border-l-0 bg-gray-50 py-2 px-3 text-gray-500">
                            ₫
                          </span>
                        </div>
                      </div>
                      <div className="col-start-6 col-end-8">
                        <label
                          className="mb-1 block text-sm text-gray-700"
                          htmlFor="name"
                        >
                          Cost Price
                        </label>
                        <div className="flex rounded shadow">
                          <input
                            className="focus:shadow-outline z-10 block w-full w-full flex-1 appearance-none rounded-none rounded-l border py-2 px-3 leading-tight text-gray-700 focus:outline-none focus:ring focus:ring-indigo-200"
                            type="text"
                            defaultValue={costString}
                          />
                          <span className="inline-flex items-center rounded-r border border-l-0 bg-gray-50 py-2 px-3 text-gray-500">
                            ₫
                          </span>
                        </div>
                      </div>
                    </div>
                  </Disclosure.Panel>
                </Transition>
              </>
            )}
          </Disclosure>
        </div>

        <div className="rounded-md bg-white py-4 px-3 shadow-sm">
          <Disclosure defaultOpen>
            {({ open }) => (
              <>
                <Disclosure.Button className="flex w-full justify-between rounded-xl bg-gray-50 px-4 py-3">
                  <h3 className="text-lg font-medium leading-6 text-gray-900">
                    Media
                  </h3>
                  <ChevronUpIcon
                    className={cc([
                      "h-4 h-4",
                      { "rotate-180 transform": open },
                    ])}
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
                >
                  <Disclosure.Panel static className="p-4">
                    <AssetManagement
                      productID={productID.toString()}
                      mediaList={product.mediaList}
                    />
                  </Disclosure.Panel>
                </Transition>
              </>
            )}
          </Disclosure>
        </div>

        <div className="rounded-md bg-white py-4 px-3 shadow-sm">
          <Disclosure defaultOpen>
            {({ open }) => (
              <>
                <Disclosure.Button className="flex w-full justify-between rounded-xl bg-gray-50 px-4 py-3">
                  <h3 className="text-lg font-medium leading-6 text-gray-900">
                    Platforms
                  </h3>
                  <ChevronUpIcon
                    className={cc([
                      "h-4 h-4",
                      { "rotate-180 transform": open },
                    ])}
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
                                  <span className="capitalize">
                                    {p.platformName}
                                  </span>
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
                                        "rounded-md py-1 px-2",
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
                                      {p.status}
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
        </div>
      </div>
    </div>
  );
}

DashboardProductDetails.getLayout = function getLayout(page: ReactElement) {
  return (
    <DashboardLayout
      title={defineMessage({
        description: "Dashboard home title",
        defaultMessage: "Home",
      })}
    >
      {page}
    </DashboardLayout>
  );
};

export default withAuthPage(DashboardProductDetails);
