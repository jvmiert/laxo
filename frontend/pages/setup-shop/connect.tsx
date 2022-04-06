import type { ReactElement } from "react";
import { useIntl } from "react-intl";
import { InferGetServerSidePropsType, GetServerSideProps } from "next";
import { useRouter } from "next/router";
import { fromUnixTime, formatDistance } from "date-fns";
import { enUS, vi } from "date-fns/locale";
import { withRedirectUnauth, withAuthPage } from "@/lib/withAuth";
import DefaultLayout from "@/components/DefaultLayout";
import { useGetShop } from "@/hooks/swrHooks";
import { DotsVerticalIcon } from "@heroicons/react/outline";
import LazadaIcon from "@/components/icons/LazadaIcon";
import ShopeeIcon from "@/components/icons/ShopeeIcon";
import TikiIcon from "@/components/icons/TikiIcon";

export const getServerSideProps: GetServerSideProps = withRedirectUnauth();

type ConnectStoreProps = InferGetServerSidePropsType<typeof getServerSideProps>;

function ConnectStore(props: ConnectStoreProps) {
  const t = useIntl();
  const { locale } = useRouter();
  const { shops, loading } = useGetShop();

  const getButtonDisabled = (platform: string): boolean => {
    if (shops.shops.length == 0) return true;

    if (shops.shops[0].platforms.length == 0) return false;

    return (
      shops.shops[0].platforms.filter((e) => e.name.toLowerCase() == platform)
        .length > 0
    );
  };

  return (
    <>
      <h1 className="mb-3 mt-5 text-2xl font-bold">
        {t.formatMessage({
          defaultMessage: "Connect an e-commerce platform",
          description: "Connect shop: Form header",
        })}
      </h1>
      <p>
        {t.formatMessage({
          defaultMessage:
            "Add one of the platforms so we can use it to retrieve your shop information such as products.",
          description: "Connect shop: Form description",
        })}
      </p>
      {!loading && (
        <div className="mt-6 max-w-xl rounded-md border border-gray-100 bg-gray-50 p-6">
          <p className="mb-4 border-b border-gray-200 pb-4 font-bold">
            Add New
          </p>
          <div className="flex space-x-4">
            {!getButtonDisabled("shopee") && (
              <button
                className="flex w-40 justify-center rounded-md bg-gradient-to-r from-[#ff9c68] to-[#ff5422] py-2 px-4 font-bold text-white shadow-lg shadow-orange-600/50 hover:from-orange-700 hover:to-orange-700 focus:outline-none focus:ring focus:ring-orange-200"
                type="submit"
              >
                {<ShopeeIcon />}{" "}
                {t.formatMessage({
                  defaultMessage: "Shopee",
                  description: "Connect Page: Connect Shopee Button",
                })}
              </button>
            )}
            {!getButtonDisabled("lazada") && (
              <button
                className="flex w-40 justify-center rounded-md bg-gradient-to-r from-[#37D8FF] to-[#972BFF] py-2 px-4 font-bold text-white shadow-lg shadow-blue-600/50 hover:from-blue-700 hover:to-blue-700 focus:outline-none focus:ring focus:ring-blue-200"
                type="submit"
              >
                {<LazadaIcon />}{" "}
                {t.formatMessage({
                  defaultMessage: "Lazada",
                  description: "Connect Page: Connect Lazada Button",
                })}
              </button>
            )}
            {!getButtonDisabled("tiki") && (
              <button
                className="flex w-40 justify-center rounded-md bg-gradient-to-r from-[#1a94ff] to-[#1083e8] py-2 px-4 font-bold text-white shadow-lg shadow-blue-600/50 hover:from-blue-700 hover:to-blue-700 focus:outline-none focus:ring focus:ring-blue-200"
                type="submit"
              >
                {<TikiIcon />}{" "}
                {t.formatMessage({
                  defaultMessage: "Tiki",
                  description: "Connect Page: Connect Tiki Button",
                })}
              </button>
            )}
          </div>
        </div>
      )}
      {!loading &&
        shops.shops.length > 0 &&
        shops.shops[0].platforms.length > 0 && (
          <div className="mt-6 max-w-xl rounded-md border border-gray-100 p-6">
            {shops?.shops[0].platforms.map((platform) => (
              <div key={platform.id} className="flex flex-row justify-between">
                <p className="capitalize">{platform.name}</p>
                <div>
                  <p>
                    {formatDistance(
                      fromUnixTime(platform.created),
                      new Date(),
                      {
                        addSuffix: true,
                        locale: locale == "vi" ? vi : enUS,
                      },
                    )}
                    <DotsVerticalIcon className="ml-4 inline h-4 w-4" />
                  </p>
                </div>
              </div>
            ))}
          </div>
        )}
    </>
  );
}

ConnectStore.getLayout = function getLayout(page: ReactElement) {
  return <DefaultLayout>{page}</DefaultLayout>;
};

export default withAuthPage(ConnectStore);
