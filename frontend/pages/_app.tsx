import "../styles/fonts.css";
import "../styles/globals.css";
import { IntlProvider } from "react-intl";
import { IntlErrorCode } from "@formatjs/intl";
import { SWRConfig } from "swr";
import { useRouter } from "next/router";
import { AxiosProvider } from "@/providers/AxiosProvider";
import { AuthProvider } from "@/providers/AuthProvider";
import { DashboardProvider } from "@/providers/DashboardProvider";
import { AppPropsWithLayout } from "@/types/pages";

import messages_en from "../compiled-lang/en.json";
import messages_vi from "../compiled-lang/vi.json";
// psuedo locale for development
import messages_sw from "../compiled-lang/sw.json";

type LocalesType = {
  [key: string]: any;
};

const languages: LocalesType = {
  en: { ...messages_en },
  vi: { ...messages_vi },
  sw: { ...messages_sw },
};

function MyApp({ Component, pageProps }: AppPropsWithLayout) {
  const { locale } = useRouter();
  const laxoLocale = locale ? (locale === "default" ? "vi" : locale) : "vi";

  const getLayout = Component.getLayout ?? ((page) => page);

  return (
    <SWRConfig
      value={pageProps?.fallback ? { fallback: pageProps.fallback } : {}}
    >
      <IntlProvider
        locale={laxoLocale}
        defaultLocale={"vi"}
        messages={languages[laxoLocale]}
        onError={(err) => {
          // Disabling missing translation warning for development
          if (
            err.code === IntlErrorCode.MISSING_TRANSLATION &&
            process.env.NODE_ENV === "development"
          ) {
            //console.warn("Missing translation", err.message);
            return;
          }
          throw err;
        }}
      >
        <AxiosProvider>
          <AuthProvider>
            <DashboardProvider>
              {getLayout(<Component {...pageProps} />)}
            </DashboardProvider>
          </AuthProvider>
        </AxiosProvider>
      </IntlProvider>
    </SWRConfig>
  );
}

export default MyApp;
