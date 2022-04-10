import type { ReactElement } from "react";
import { useIntl } from "react-intl";
import DashboardLayout from "@/components/DashboardLayout";
import { withRedirectUnauth, withAuthPage } from "@/lib/withAuth";
import { InferGetServerSidePropsType, GetServerSideProps } from "next";
import { useGetShop } from "@/hooks/swrHooks";
import ShopNotMadeNotification from "@/components/ShopNotMadeNotification";

export const getServerSideProps: GetServerSideProps = withRedirectUnauth();

type DashboardSettingsProps = InferGetServerSidePropsType<
  typeof getServerSideProps
>;

function DashboardSettings(props: DashboardSettingsProps) {
  const t = useIntl();
  const { shops } = useGetShop();

  if (shops.total < 1) return <ShopNotMadeNotification />;

  return (
    <div className="flex flex-1">
      <div className="max-w-xl grow space-y-4 rounded-md border border-gray-100 p-6">
        <h3 className="text-lg font-medium">
          {t.formatMessage({
            defaultMessage: "Your Shop Name",
            description: "Dashboard shop settings: Shop name header",
          })}
        </h3>
        <p>
          {t.formatMessage({
            defaultMessage:
              "This is what your shop is called within our platform.",
            description: "Dashboard shop settings: Shop name description",
          })}
        </p>
        <input
          readOnly
          className="appearance-none rounded border bg-gray-50 py-2 px-3 leading-tight text-gray-700 text-gray-500 shadow"
          value={shops.shops[0].name}
        />
      </div>
    </div>
  );
}

DashboardSettings.getLayout = function getLayout(page: ReactElement) {
  return <DashboardLayout title="Settings">{page}</DashboardLayout>;
};

export default withAuthPage(DashboardSettings);
