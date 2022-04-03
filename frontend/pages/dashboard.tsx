import type { ReactElement } from "react";
import { useIntl } from "react-intl";
import Link from "next/link";
import Head from "next/head";
import DefaultLayout from "@/components/DefaultLayout";
import { withRedirectUnauth, withAuthPage } from "@/lib/withAuth";
import { InferGetServerSidePropsType, GetServerSideProps } from "next";
import { useGetShop } from "@/hooks/swrHooks";
import { ArrowNarrowRightIcon } from "@heroicons/react/outline";

export const getServerSideProps: GetServerSideProps = withRedirectUnauth();

type DashboardPageProps = InferGetServerSidePropsType<
  typeof getServerSideProps
>;

function DashboardPage(props: DashboardPageProps) {
  const t = useIntl();
  const { shops } = useGetShop();

  const connected =
    shops && shops?.total > 0 && shops.shops[0].platforms.length > 0;
  return (
    <>
      <Head>
        <title>Laxo: Dashboard</title>
        <link rel="icon" href="/favicon.ico" />
      </Head>

      <h1 className="mb-3 text-2xl font-bold">
        Dashboard{" "}
        {shops?.total > 0 && (
          <span className="text-sm font-semibold text-gray-500">
            {shops.shops[0].name}
          </span>
        )}
      </h1>
      {shops?.total < 1 && (
        <p>
          {t.formatMessage({
            defaultMessage:
              "You didn't create your shop yet, please make your shop.",
            description: "Dashboard: no shop notifcation",
          })}
          <br />
          <Link href="/setup-shop/create" passHref>
            <a className="cursor-pointer font-semibold text-indigo-500">
              {t.formatMessage({
                defaultMessage: "Create your store",
                description: "Dashboard: create store button",
              })}{" "}
              <ArrowNarrowRightIcon className="inline h-4 w-4" />
            </a>
          </Link>
        </p>
      )}
      {!connected && (
        <p>
          {t.formatMessage({
            defaultMessage:
              "Please connect at least one e-commerce platform to your shop.",
            description: "Dashboard: not connected",
          })}
          <br />
          <Link href="/setup-shop/connect" passHref>
            <a className="cursor-pointer font-semibold text-indigo-500">
              {t.formatMessage({
                defaultMessage: "Connect a platform",
                description: "Dashboard: connect store button",
              })}{" "}
              <ArrowNarrowRightIcon className="inline h-4 w-4" />
            </a>
          </Link>
        </p>
      )}
    </>
  );
}

DashboardPage.getLayout = function getLayout(page: ReactElement) {
  return <DefaultLayout>{page}</DefaultLayout>;
};

export default withAuthPage(DashboardPage);
