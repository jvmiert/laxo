import { useEffect } from 'react';
import axios from 'axios';
import type { NextPage } from 'next'
import Head from 'next/head'

const Home: NextPage = () => {

  useEffect(() => {
    const fetchData = async () => {
        const result = await axios.post(
        "/api/login",
        {
          test: 'test',
        }
      )
      console.log(result);
    };
    fetchData();
  }, []);

  return (
    <div>
      <Head>
        <title>Laxo</title>
        <link rel="icon" href="/favicon.ico" />
      </Head>

      <main>
        <p className="text-3xl font-bold underline">Hello world</p>
      </main>
    </div>
  )
}

export default Home
