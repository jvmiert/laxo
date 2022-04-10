import type { ReactElement } from "react";
import { useIntl } from "react-intl";
import { InferGetServerSidePropsType, GetServerSideProps } from "next";
import { withRedirectUnauth, withAuthPage } from "@/lib/withAuth";
import DefaultLayout from "@/components/DefaultLayout";
import PlatformConnect from "@/components/PlatformConnect";

export const getServerSideProps: GetServerSideProps = withRedirectUnauth();

type ConnectStoreProps = InferGetServerSidePropsType<typeof getServerSideProps>;

function ConnectStore(props: ConnectStoreProps) {
  const t = useIntl();

  return (
    <>
      <h1 className="mb-3 mt-5 text-2xl font-bold">
        {t.formatMessage({
          defaultMessage: "Connect an e-commerce platform",
          description: "Connect shop: Form header",
        })}
      </h1>
      <p>
        {t.formatMessage({
          defaultMessage:
            "Add one of the platforms so we can use it to retrieve your shop information such as products.",
          description: "Connect shop: Form description",
        })}
      </p>
      <PlatformConnect />
    </>
  );
}

ConnectStore.getLayout = function getLayout(page: ReactElement) {
  return <DefaultLayout>{page}</DefaultLayout>;
};

export default withAuthPage(ConnectStore);
