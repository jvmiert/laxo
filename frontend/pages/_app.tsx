import "../styles/globals.css";
import { IntlProvider } from "react-intl";
import { useRouter } from "next/router";
import type { AppProps } from "next/app";

function MyApp({ Component, pageProps }: AppProps) {
  const { locale, defaultLocale } = useRouter();
  return (
    <IntlProvider
      locale={locale!}
      defaultLocale={defaultLocale}
      messages={pageProps.intlMessages}
    >
      <Component {...pageProps} />
    </IntlProvider>
  );
}

export default MyApp;
