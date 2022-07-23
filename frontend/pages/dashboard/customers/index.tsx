import type { ReactElement } from "react";
import { InferGetServerSidePropsType, GetServerSideProps } from "next";
import { useIntl, defineMessage } from "react-intl";

import { withRedirectUnauth, withAuthPage } from "@/lib/withAuth";
import DashboardLayout from "@/components/DashboardLayout";

export const getServerSideProps: GetServerSideProps = withRedirectUnauth();

type CustomersPageProps = InferGetServerSidePropsType<
  typeof getServerSideProps
>;

function CustomersPage(props: CustomersPageProps) {
  const t = useIntl();

  return (
    <div>
      <p>Customers</p>
    </div>
  );
}

CustomersPage.getLayout = function getLayout(page: ReactElement) {
  return (
    <DashboardLayout
      title={defineMessage({
        description: "Dashboard customers title",
        defaultMessage: "Customers",
      })}
    >
      {page}
    </DashboardLayout>
  );
};

export default withAuthPage(CustomersPage);
