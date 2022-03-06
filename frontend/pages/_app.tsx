import "../styles/globals.css";
import { IntlProvider } from "react-intl";
import { SWRConfig } from "swr";
import { useRouter } from "next/router";
import type { AppProps } from "next/app";
import { AxiosProvider } from "@/providers/AxiosProvider";
import { AuthProvider } from "@/providers/AuthProvider";

import messages_en from "../compiled-lang/en.json";
import messages_vi from "../compiled-lang/vi.json";

type LocalesType = {
  [key: string]: any;
};

const languages: LocalesType = {
  en: { ...messages_en },
  vi: { ...messages_vi },
};

function MyApp({ Component, pageProps }: AppProps) {
  const { locale = "en", defaultLocale } = useRouter();

  return (
    <SWRConfig
      value={pageProps?.fallback ? { fallback: pageProps.fallback } : {}}
    >
      <IntlProvider
        locale={locale!}
        defaultLocale={defaultLocale}
        messages={languages[locale]}
      >
        <AxiosProvider>
          <AuthProvider>
            <Component {...pageProps} />
          </AuthProvider>
        </AxiosProvider>
      </IntlProvider>
    </SWRConfig>
  );
}

export default MyApp;
