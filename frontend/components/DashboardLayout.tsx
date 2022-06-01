import type { ReactNode } from "react";
import Head from "next/head";
import DashboardTopNavigation from "@/components/DashboardTopNavigation";
import DashboardNavigation from "@/components/DashboardNavigation";
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
      <DashboardTopNavigation />
      <div className="min-h-screen bg-zinc-100">
        <div className="flex w-full flex-row flex-nowrap">
          <div className="flex min-h-[55vh] shrink-0 grow-0 basis-auto bg-white">
            <DashboardNavigation />
          </div>
          <div className="flex grow flex-col">
            <div className="flex justify-between bg-zinc-100 p-8">
              <main className="grow">{children}</main>
            </div>
          </div>
        </div>
      </div>
    </>
  );
}
