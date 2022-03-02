import { useIntl } from "react-intl";
import Head from "next/head";
import Navigation from "@/components/Navigation";
import loadIntlMessages from "@/helpers/loadIntlMessages";
import type { LoadI18nMessagesProps } from "@/helpers/loadIntlMessages";
import { InferGetStaticPropsType } from "next";
import { useAuth } from "@/providers/AuthProvider";

export async function getStaticProps(ctx: LoadI18nMessagesProps) {
  return {
    props: {
      intlMessages: await loadIntlMessages(ctx),
    },
  };
}

type HomePageProps = InferGetStaticPropsType<typeof getStaticProps>;

export default function HomePage(props: HomePageProps) {
  const t = useIntl();
  const { auth } = useAuth();
  return (
    <div>
      <Head>
        <title>Laxo</title>
        <link rel="icon" href="/favicon.ico" />
      </Head>

      <Navigation />
      <main>
        <p className="text-3xl font-bold underline">
          {t.formatMessage({
            defaultMessage: "Hello World",
            description: "Index Page: title",
          })}
        </p>
        <p>You are {!auth && "not"} authenticated!</p>
      </main>
    </div>
  );
}
