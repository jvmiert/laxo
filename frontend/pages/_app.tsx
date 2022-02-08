import "../styles/globals.css";
import { IntlProvider } from "react-intl";
import { useRouter } from "next/router";
import type { AppProps } from "next/app";
import { AxiosProvider } from "@/providers/AxiosProvider";

function MyApp({ Component, pageProps }: AppProps) {
  const { locale, defaultLocale } = useRouter();
  return (
    <IntlProvider
      locale={locale!}
      defaultLocale={defaultLocale}
      messages={pageProps.intlMessages}
    >
      <AxiosProvider>
        <Component {...pageProps} />
      </AxiosProvider>
    </IntlProvider>
  );
}

export default MyApp;
