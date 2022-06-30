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
import DetailsGeneralEdit from "@/components/dashboard/product/DetailsGeneralEdit";
import DetailsChangedNotification from "@/components/dashboard/product/DetailsChangedNotification";
import DetailsPlatformsEdit from "@/components/dashboard/product/DetailsPlatformsEdit";
import AssetInsertDialog from "@/components/dashboard/product/AssetInsertDialog";

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

  return (
    <div className="mx-auto max-w-5xl">
      <DetailsChangedNotification />
      <AssetInsertDialog />
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
                  unmount={false}
                >
                  <Disclosure.Panel static className="p-4">
                    <DetailsGeneralEdit product={product.product} />
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
                  unmount={false}
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
          <DetailsPlatformsEdit product={product} />
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
