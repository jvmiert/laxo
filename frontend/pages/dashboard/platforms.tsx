import type { ReactElement } from "react";
import { useIntl } from "react-intl";
import DashboardLayout from "@/components/DashboardLayout";
import { withRedirectUnauth, withAuthPage } from "@/lib/withAuth";
import { InferGetServerSidePropsType, GetServerSideProps } from "next";
import PlatformConnect from "@/components/PlatformConnect";
import ShopNotMadeNotification from "@/components/ShopNotMadeNotification";
import { useGetShop } from "@/hooks/swrHooks";

export const getServerSideProps: GetServerSideProps = withRedirectUnauth();

type PlatformSettingsProps = InferGetServerSidePropsType<
  typeof getServerSideProps
>;

function PlatformSettings(props: PlatformSettingsProps) {
  const t = useIntl();
  const { shops } = useGetShop();

  if (shops.total < 1) return <ShopNotMadeNotification />;

  return (
    <>
      <p>
        {t.formatMessage({
          defaultMessage:
            "Below you can see to which platforms your shop is connected and add new platforms.",
          description: "Dashboard platforms: Description",
        })}
      </p>
      <PlatformConnect />
    </>
  );
}

PlatformSettings.getLayout = function getLayout(page: ReactElement) {
  return <DashboardLayout title="Platforms">{page}</DashboardLayout>;
};

export default withAuthPage(PlatformSettings);
