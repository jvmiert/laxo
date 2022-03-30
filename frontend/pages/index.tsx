import type { ReactElement } from "react";
import { useIntl } from "react-intl";
import Head from "next/head";
import DefaultLayout from "@/components/DefaultLayout";
import { useAuth } from "@/providers/AuthProvider";

export default function HomePage() {
  const t = useIntl();
  const { auth } = useAuth();
  return (
    <>
      <Head>
        <title>Laxo</title>
        <link rel="icon" href="/favicon.ico" />
      </Head>

      <p className="text-3xl font-bold underline">
        {t.formatMessage({
          defaultMessage: "Hello World",
          description: "Index Page: title",
        })}
      </p>
      <p>You are {!auth && "not"} authenticated!</p>
    </>
  );
}

HomePage.getLayout = function getLayout(page: ReactElement) {
  return <DefaultLayout>{page}</DefaultLayout>;
};
