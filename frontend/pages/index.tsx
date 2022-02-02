import { useIntl } from "react-intl";
import Link from "next/link";
import Head from "next/head";
import loadIntlMessages from "@/helpers/loadIntlMessages";
import type { LoadI18nMessagesProps } from "@/helpers/loadIntlMessages";
import { InferGetStaticPropsType } from "next";

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
        <p className="text-3xl font-bold underline">
          {t.formatMessage({
            defaultMessage: "Hello World",
            description: "Index Page: title",
          })}
        </p>
      </main>
    </div>
  );
}
