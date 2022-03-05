import { useIntl } from "react-intl";
import Head from "next/head";
import Navigation from "@/components/Navigation";
import { useAuth } from "@/providers/AuthProvider";

export default function HomePage() {
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
