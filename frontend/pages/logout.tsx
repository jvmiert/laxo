import { useEffect } from "react";
import { useRouter } from "next/router";
import Head from "next/head";
import Navigation from "@/components/Navigation";
import useLoginApi from "@/hooks/useLoginApi";

export default function LogoutPage() {
  const { push } = useRouter();
  const { doLogout } = useLoginApi();

  useEffect(() => {
    const logout = async () => {
      // @TODO: handle error
      const { success } = await doLogout();

      if (success) {
        push("/");
      }
    };

    logout();
  }, [doLogout, push]);

  return (
    <div>
      <Head>
        <title>Laxo - Logout</title>
        <link rel="icon" href="/favicon.ico" />
      </Head>

      <Navigation />
      <main>
        <p>...</p>
      </main>
    </div>
  );
}
