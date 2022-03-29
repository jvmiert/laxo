import { useIntl } from "react-intl";
import Head from "next/head";
import useCreateFrame from "@/hooks/useCreateFrame";
import Navigation from "@/components/Navigation";
import { withRedirectUnauth, withAuthPage } from "@/lib/withAuth";
import { InferGetServerSidePropsType, GetServerSideProps } from "next";

export const getServerSideProps: GetServerSideProps = withRedirectUnauth();

type DashboardPageProps = InferGetServerSidePropsType<
  typeof getServerSideProps
>;

export default withAuthPage(function DashboardPage(props: DashboardPageProps) {
  const t = useIntl();
  const { createFrame } = useCreateFrame();
  return (
    <div>
      <Head>
        <title>Laxo - Dashboard</title>
        <link rel="icon" href="/favicon.ico" />
      </Head>

      <Navigation />
      <main>
        <p className="text-3xl font-bold underline">Dashboard</p>
        <button onClick={createFrame}>Create a frame test</button>
      </main>
    </div>
  );
});