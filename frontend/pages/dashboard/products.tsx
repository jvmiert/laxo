import type { ReactElement } from "react";
import { useIntl, defineMessage } from "react-intl";
import DashboardLayout from "@/components/DashboardLayout";
import { withRedirectUnauth, withAuthPage } from "@/lib/withAuth";
import { InferGetServerSidePropsType, GetServerSideProps } from "next";
import { useGetLaxoProducts } from "@/hooks/swrHooks";

export const getServerSideProps: GetServerSideProps = withRedirectUnauth();

type DashboardProductsPageProps = InferGetServerSidePropsType<
  typeof getServerSideProps
>;

function DashboardProductsPage(props: DashboardProductsPageProps) {
  const t = useIntl();
  const { products } = useGetLaxoProducts(0, 50);
  return <div className="break-all">{JSON.stringify(products)}</div>;
}

DashboardProductsPage.getLayout = function getLayout(page: ReactElement) {
  return (
    <DashboardLayout
      title={defineMessage({
        description: "Dashboard products title",
        defaultMessage: "Products",
      })}
    >
      {page}
    </DashboardLayout>
  );
};

export default withAuthPage(DashboardProductsPage);
