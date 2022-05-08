import { useEffect } from "react";
import type { ReactElement } from "react";
import { useIntl } from "react-intl";
import axios from "axios";
import Head from "next/head";
import DefaultLayout from "@/components/DefaultLayout";

export default function TestPage() {
  const t = useIntl();

  useEffect(() => {
    const executeTask = async () => {
      const result = await axios("/api/test");

      console.log(result.data);
    };

    executeTask();
  }, []);

  return (
    <>
      <Head>
        <title>Laxo: Test</title>
        <link rel="icon" href="/favicon.ico" />
      </Head>

      <p className="text-lg">Testing 1, 2, 3...</p>
    </>
  );
}

TestPage.getLayout = function getLayout(page: ReactElement) {
  return <DefaultLayout>{page}</DefaultLayout>;
};
