import "../styles/globals.css";
import { IntlProvider } from "react-intl";
import { useRouter } from "next/router";
import type { AppProps } from "next/app";
import { AxiosProvider } from "@/providers/AxiosProvider";
import { AuthProvider } from "@/providers/AuthProvider";

function MyApp({ Component, pageProps }: AppProps) {
  const { locale, defaultLocale } = useRouter();
  return (
    <IntlProvider
      locale={locale!}
      defaultLocale={defaultLocale}
      messages={pageProps.intlMessages}
    >
      <AxiosProvider>
        <AuthProvider>
          <Component {...pageProps} />
        </AuthProvider>
      </AxiosProvider>
    </IntlProvider>
  );
}

export default MyApp;
