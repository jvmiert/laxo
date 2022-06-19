import type { ReactNode } from "react";
import Head from "next/head";
import DashboardTopNavigation from "@/components/DashboardTopNavigation";
import DashboardNavigation from "@/components/DashboardNavigation";
import AlertContainer from "@/components/dashboard/AlertContainer";
import { useIntl, MessageDescriptor } from "react-intl";

type DefaultLayoutProps = {
  children: ReactNode;
  title: MessageDescriptor;
};

export default function DashboardLayout({
  children,
  title,
}: DefaultLayoutProps) {
  const t = useIntl();

  return (
    <>
      <Head>
        <title>{`Laxo: ${t.formatMessage(title)}`}</title>
        <link rel="icon" href="/favicon.ico" />
      </Head>
      <div>
        <AlertContainer />
        <DashboardNavigation />
        <div className="flex h-screen flex-col overflow-y-hidden bg-gray-100">
          <div className="ml-52 bg-zinc-100">
            <DashboardTopNavigation />
          </div>
          <div className="ml-52 grow overflow-x-auto overflow-y-scroll">
            <main className="my-8 mx-6">{children}</main>
          </div>
        </div>
      </div>
    </>
  );
}
