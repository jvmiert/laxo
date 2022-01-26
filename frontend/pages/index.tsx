import { useEffect } from "react";
import type { NextPage } from "next";
import Link from "next/link";
import Head from "next/head";

const Home: NextPage = () => {
  return (
    <div>
      <Head>
        <title>Laxo</title>
        <link rel="icon" href="/favicon.ico" />
      </Head>

      <main>
        <Link href="/login">
          <a>Login</a>
        </Link>
        <p className="text-3xl font-bold underline">Hello world</p>
      </main>
    </div>
  );
};

export default Home;
