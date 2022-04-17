import { useEffect, useState } from "react";
import type { ReactElement } from "react";
import { useIntl } from "react-intl";
import axios from "axios";
import Head from "next/head";
import DefaultLayout from "@/components/DefaultLayout";
import useShopApi from "@/hooks/useShopApi";

export default function TestPage() {
  const t = useIntl();
  const { getProductRetrieveStatusUpdate } = useShopApi();

  const [state, setState] = useState(null);
  const [called, setCalled] = useState(false);

  useEffect(() => {
    if (state && !called) {
      getProductRetrieveStatusUpdate(state);
      setCalled(true);
    }
  }, [state, getProductRetrieveStatusUpdate, called]);

  useEffect(() => {
    const executeTask = async () => {
      const result = await axios("/api/test");

      setState(result.data);
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
