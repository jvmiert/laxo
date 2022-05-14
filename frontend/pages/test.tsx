import type { ReactElement } from "react";
import axios from "axios";
import Head from "next/head";
import DefaultLayout from "@/components/DefaultLayout";

export default function TestPage() {
  const taskExecution = () => {
    const executeTask = async () => {
      await axios("/api/test/task");
    };

    executeTask();
  };

  const taskLazada = () => {
    const executeTask = async () => {
      await axios("/api/test/lazada");
    };

    executeTask();
  };
  return (
    <>
      <Head>
        <title>Laxo: Test</title>
        <link rel="icon" href="/favicon.ico" />
      </Head>

      <p className="pb-5 text-lg">Testing 1, 2, 3...</p>
      <button
        type="button"
        onClick={taskExecution}
        className="rounded-md border border-gray-300 bg-white py-2 px-3 text-sm font-medium leading-4 text-gray-700 shadow-sm hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2"
      >
        Execute notification
      </button>
      <button
        type="button"
        onClick={taskLazada}
        className="rounded-md border border-gray-300 bg-white py-2 px-3 text-sm font-medium leading-4 text-gray-700 shadow-sm hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2"
      >
        Execute Lazada
      </button>
    </>
  );
}

TestPage.getLayout = function getLayout(page: ReactElement) {
  return <DefaultLayout>{page}</DefaultLayout>;
};
