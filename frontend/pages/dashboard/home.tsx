import { InferGetServerSidePropsType, GetServerSideProps } from "next";
import type { ReactElement } from "react";
import { useIntl, defineMessage } from "react-intl";
import Link from "next/link";
import ShopNotMadeNotification from "@/components/ShopNotMadeNotification";
import { ArrowNarrowRightIcon } from "@heroicons/react/outline";

import { withRedirectUnauth, withAuthPage } from "@/lib/withAuth";
import { useGetShop } from "@/hooks/swrHooks";
import DashboardLayout from "@/components/DashboardLayout";

export const getServerSideProps: GetServerSideProps = withRedirectUnauth();

type DashboardPageProps = InferGetServerSidePropsType<
  typeof getServerSideProps
>;

function DashboardPage(props: DashboardPageProps) {
  const t = useIntl();
  const { shops } = useGetShop();

  const connected = shops.total > 0 && shops.shops[0].platforms.length > 0;

  return (
    <>
      {shops.total < 1 && <ShopNotMadeNotification />}
      {!connected && shops.total > 0 && (
        <p>
          {t.formatMessage({
            defaultMessage:
              "Please connect at least one e-commerce platform to your shop.",
            description: "Dashboard: not connected",
          })}
          <br />
          <Link href="/setup-shop/connect" passHref>
            <a className="cursor-pointer font-semibold text-indigo-500">
              {t.formatMessage({
                defaultMessage: "Connect a platform",
                description: "Dashboard: connect store button",
              })}{" "}
              <ArrowNarrowRightIcon className="inline h-4 w-4" />
            </a>
          </Link>
        </p>
      )}
    </>
  );
}

DashboardPage.getLayout = function getLayout(page: ReactElement) {
  return (
    <DashboardLayout
      title={defineMessage({
        description: "Dashboard home title",
        defaultMessage: "Home",
      })}
    >
      {page}
    </DashboardLayout>
  );
};

export default withAuthPage(DashboardPage);
