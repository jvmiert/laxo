import type { ReactElement } from "react";
import { useIntl, FormattedMessage } from "react-intl";
import { useRouter } from 'next/router'
import DashboardLayout from "@/components/DashboardLayout";
import NestedDashboardPlatformTitle from "@/components/NestedDashboardPlatformTitle";
import { withRedirectUnauth, withAuthPage } from "@/lib/withAuth";
import { InferGetServerSidePropsType, GetServerSideProps } from "next";

export const getServerSideProps: GetServerSideProps = withRedirectUnauth();

type PlatformSpecificSettingsProps = InferGetServerSidePropsType<
  typeof getServerSideProps
>;

function PlatformSpecificSettings(props: PlatformSpecificSettingsProps) {
  const t = useIntl();

  const { query: { platform } } = useRouter();
  return (
    <>
      <p>Hello {platform}</p>
    </>
  );
}

PlatformSpecificSettings.getLayout = function getLayout(page: ReactElement) {
  return <DashboardLayout title={<NestedDashboardPlatformTitle />}>{page}</DashboardLayout>;
};

export default withAuthPage(PlatformSpecificSettings);
