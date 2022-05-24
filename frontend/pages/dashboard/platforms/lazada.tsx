import cc from "classcat";
import type { ReactElement } from "react";
import { useIntl, defineMessage } from "react-intl";
import DashboardLayout from "@/components/DashboardLayout";
import { withRedirectUnauth, withAuthPage } from "@/lib/withAuth";
import { InferGetServerSidePropsType, GetServerSideProps } from "next";
import {
  useGetLazadaPlatform,
  useGetShop,
  useGetShopPlatformsRedirect,
} from "@/hooks/swrHooks";
import LazadaIcon from "@/components/icons/LazadaIcon";
import { useRouter } from "next/router";
import { formatDistance } from "date-fns";
import { enUS, vi } from "date-fns/locale";

export const getServerSideProps: GetServerSideProps = withRedirectUnauth();

type PlatformLazadaSettingsProps = InferGetServerSidePropsType<
  typeof getServerSideProps
>;

function PlatformLazadaSettings(props: PlatformLazadaSettingsProps) {
  const { locale } = useRouter();
  const t = useIntl();

  const { shops } = useGetShop();

  const shopID = shops.shops.length > 0 ? shops.shops[0].id : "";

  const { platforms } = useGetShopPlatformsRedirect(shopID);

  const lazadaLink = platforms.platforms.find(
    (p) => p.platform == "lazada",
  )?.url;

  const { platform: platformData } = useGetLazadaPlatform();

  //@TODO: loading state would be good
  if (!platformData) return <></>;

  const now = new Date();
  const connected =
    now < platformData.refreshExpiresIn || now < platformData.refreshExpiresIn;

  return (
    <>
      <div className="flex content-between">
        <div className="flex items-center">
          <div className="mr-5">
            <LazadaIcon className="h-10 w-10" />
          </div>
          <div>
            <h3 className="mb-2 text-2xl font-semibold">Lazada</h3>
            <div>
              <span
                className={cc([
                  "rounded-lg",
                  { "bg-green-50": connected },
                  { "bg-red-50": !connected },
                  "py-1 px-2 text-sm",
                  { "text-green-600": connected },
                  { "text-red-600": !connected },
                ])}
              >
                {connected
                  ? t.formatMessage({
                      description: "Lazada platform connected",
                      defaultMessage: "Connected",
                    })
                  : t.formatMessage({
                      description: "Lazada platform disconnected",
                      defaultMessage: "Disconnected",
                    })}
              </span>
              {!connected && (
                <p className="pt-2 text-xs italic">
                  {t.formatMessage({
                    description: "Lazada platform disconnected helper message",
                    defaultMessage:
                      "Your sales channel is disconnected, please update your credentials",
                  })}
                </p>
              )}
            </div>
          </div>
        </div>
        {lazadaLink && (
          <div className="ml-8">
            <a
              className="w-full rounded-xl border-white py-2 px-4 text-sm shadow-md hover:bg-gray-50 focus:outline-none focus:ring focus:ring-indigo-200"
              href={lazadaLink}
            >
              {t.formatMessage({
                defaultMessage: "Update Credentials",
                description: "Lazada update credentials button",
              })}
            </a>
          </div>
        )}
      </div>
      <div>
        {formatDistance(platformData.created, new Date(), {
          addSuffix: true,
          locale: locale == "vi" ? vi : enUS,
        })}
      </div>
    </>
  );
}

PlatformLazadaSettings.getLayout = function getLayout(page: ReactElement) {
  return (
    <DashboardLayout
      title={defineMessage({
        description: "Dashboard specific platform lazada title",
        defaultMessage: "Lazada",
      })}
    >
      {page}
    </DashboardLayout>
  );
};

export default withAuthPage(PlatformLazadaSettings);
