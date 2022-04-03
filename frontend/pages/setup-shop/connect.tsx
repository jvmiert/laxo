import type { ReactElement } from "react";
import { useIntl } from "react-intl";
import { InferGetServerSidePropsType, GetServerSideProps } from "next";
import { useRouter } from "next/router";
import { fromUnixTime, formatDistance } from "date-fns";
import { enUS, vi } from "date-fns/locale";
import { withRedirectUnauth, withAuthPage } from "@/lib/withAuth";
import DefaultLayout from "@/components/DefaultLayout";
import { useGetShop } from "@/hooks/swrHooks";

export const getServerSideProps: GetServerSideProps = withRedirectUnauth();

type ConnectStoreProps = InferGetServerSidePropsType<typeof getServerSideProps>;

function ConnectStore(props: ConnectStoreProps) {
  const t = useIntl();
  const { locale } = useRouter();
  const { shops, loading } = useGetShop();

  const getButtonDisabled = (platform: string): boolean => {
    if (!shops) return true;

    if (shops.shops.length == 0) return true;

    if (shops.shops[0].platforms.length == 0) return true;

    // @TODO: type this stuff :(
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
            <button
              disabled={getButtonDisabled("shopee")}
              className="w-full rounded-md bg-indigo-500 py-2 px-4 font-bold text-white shadow-lg shadow-indigo-500/50 hover:bg-indigo-700 focus:outline-none focus:ring focus:ring-indigo-200 disabled:cursor-not-allowed disabled:bg-indigo-200"
              type="submit"
            >
              {t.formatMessage({
                defaultMessage: "Connect Shopee",
                description: "Connect Page: Connect Shopee Button",
              })}
            </button>
            <button
              disabled={getButtonDisabled("lazada")}
              className="w-full rounded-md bg-indigo-500 py-2 px-4 font-bold text-white shadow-lg shadow-indigo-500/50 hover:bg-indigo-700 focus:outline-none focus:ring focus:ring-indigo-200 disabled:cursor-not-allowed disabled:bg-indigo-200"
              type="submit"
            >
              {t.formatMessage({
                defaultMessage: "Connect Lazada",
                description: "Connect Page: Connect Lazada Button",
              })}
            </button>
            <button
              disabled={getButtonDisabled("tiki")}
              className="w-full rounded-md bg-indigo-500 py-2 px-4 font-bold text-white shadow-lg shadow-indigo-500/50 hover:bg-indigo-700 focus:outline-none focus:ring focus:ring-indigo-200 disabled:cursor-not-allowed disabled:bg-indigo-200"
              type="submit"
            >
              {t.formatMessage({
                defaultMessage: "Connect Tiki",
                description: "Connect Page: Connect Tiki Button",
              })}
            </button>
          </div>
        </div>
      )}
      <div className="mt-6 max-w-xl rounded-md border border-gray-100 p-6">
        {shops?.shops[0].platforms.map((platform) => (
          <div key={platform.id} className="flex flex-row justify-between">
            <p>{platform.name}</p>
            <p>
              {formatDistance(fromUnixTime(platform.created), new Date(), {
                addSuffix: true,
                locale: locale == "vi" ? vi : enUS,
              })}
            </p>
          </div>
        ))}
      </div>
    </>
  );
}

ConnectStore.getLayout = function getLayout(page: ReactElement) {
  return <DefaultLayout>{page}</DefaultLayout>;
};

export default withAuthPage(ConnectStore);
