import type { ReactElement } from "react";
import { InferGetServerSidePropsType, GetServerSideProps } from "next";
import { useIntl, defineMessage } from "react-intl";

import { withRedirectUnauth, withAuthPage } from "@/lib/withAuth";
import DashboardLayout from "@/components/DashboardLayout";

export const getServerSideProps: GetServerSideProps = withRedirectUnauth();

type OrdersPageProps = InferGetServerSidePropsType<typeof getServerSideProps>;

function OrdersPage(props: OrdersPageProps) {
  const t = useIntl();

  return (
    <div>
      <p>Orders</p>
    </div>
  );
}

OrdersPage.getLayout = function getLayout(page: ReactElement) {
  return (
    <DashboardLayout
      title={defineMessage({
        description: "Dashboard orders title",
        defaultMessage: "Orders",
      })}
    >
      {page}
    </DashboardLayout>
  );
};

export default withAuthPage(OrdersPage);
