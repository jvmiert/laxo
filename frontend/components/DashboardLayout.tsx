import type { ReactNode } from "react";
import Head from "next/head";
import Navigation from "@/components/Navigation";
import DashboardNavigation from "@/components/DashboardNavigation";
import DashboardNotificationControl from "@/components/DashboardNotificationControl";
import DashboardNotificationArea from "@/components/DashboardNotificationArea";
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
      <Navigation />
      <div className="container mx-auto px-4 pt-4">
        <div className="flex w-full flex-row flex-nowrap">
          <div className="flex min-h-[55vh] shrink-0 grow-0 basis-auto">
            <DashboardNavigation />
          </div>
          <div className="ml-6 flex grow flex-col">
            <div className="flex">
              <div className="flex-grow">
                <h1 className="mb-4 text-xl font-semibold">{t.formatMessage(title)}</h1>
              </div>
              <div>
                <DashboardNotificationControl />
              </div>
            </div>
            <div className="flex justify-between">
              <main className="shrink">{children}</main>
              <DashboardNotificationArea />
            </div>
          </div>
        </div>
      </div>
    </>
  );
}
