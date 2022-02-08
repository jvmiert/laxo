import { InferGetStaticPropsType } from "next";
import Head from "next/head";
import { Form, Field } from "react-final-form";
import loadIntlMessages from "@/helpers/loadIntlMessages";
import type { LoadI18nMessagesProps } from "@/helpers/loadIntlMessages";

export async function getStaticProps(ctx: LoadI18nMessagesProps) {
  return {
    props: {
      intlMessages: await loadIntlMessages(ctx),
    },
  };
}

type RegisterPageProps = InferGetStaticPropsType<typeof getStaticProps>;

export default function LoginPage(props: RegisterPageProps) {
  return (
    <div>
      <Head>
        <title>Laxo - Register</title>
        <link rel="icon" href="/favicon.ico" />
      </Head>

      <main>
        <p className="text-3xl font-bold underline">Register</p>
        Register
      </main>
    </div>
  );
}
