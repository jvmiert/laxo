import type { ReactElement } from "react";
import { useIntl } from "react-intl";
import Head from "next/head";
import DefaultLayout from "@/components/DefaultLayout";
import { withRedirectUnauth, withAuthPage } from "@/lib/withAuth";
import { InferGetServerSidePropsType, GetServerSideProps } from "next";
import { useGetShop } from "@/hooks/swrHooks";

export const getServerSideProps: GetServerSideProps = withRedirectUnauth();

type DashboardPageProps = InferGetServerSidePropsType<
  typeof getServerSideProps
>;

function DashboardPage(props: DashboardPageProps) {
  const t = useIntl();
  const { shops } = useGetShop();
  return (
    <>
      <Head>
        <title>Laxo - Dashboard</title>
        <link rel="icon" href="/favicon.ico" />
      </Head>

      <p className="text-2xl font-bold">Dashboard</p>
      <p>{JSON.stringify(shops)}</p>
    </>
  );
}

DashboardPage.getLayout = function getLayout(page: ReactElement) {
  return <DefaultLayout>{page}</DefaultLayout>;
};

export default withAuthPage(DashboardPage);
