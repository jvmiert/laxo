import type { ReactNode } from "react";
import Head from "next/head";
import Navigation from "@/components/Navigation";
import DashboardNavigation from "@/components/DashboardNavigation";
import DashboardLoadingEvent from "@/components/DashboardLoadingEvent";
import DashboardNotificationControl from "@/components/DashboardNotificationControl";
import DashboardNotificationArea from "@/components/DashboardNotificationArea";

type DefaultLayoutProps = {
  children: ReactNode;
  title: string;
};

export default function DashboardLayout({
  children,
  title,
}: DefaultLayoutProps) {
  return (
    <>
      <Head>
        <title>{`Laxo: ${title}`}</title>
        <link rel="icon" href="/favicon.ico" />
      </Head>
      <Navigation />
      <div className="container mx-auto px-4 pt-4">
        <div className="flex w-full flex-row flex-nowrap">
          <div className="flex shrink-0 grow-0 basis-auto">
            <DashboardNavigation />
          </div>
          <div className="ml-6 flex grow flex-col">
            <DashboardLoadingEvent />
            <div className="flex">
              <div className="flex-grow">
                <h1 className="mb-4 text-xl font-semibold">{title}</h1>
              </div>
              <div>
                <DashboardNotificationControl />
              </div>
            </div>
            <div className="flex">
              <main className="shrink">{children}</main>
              <div className="w-80 grow">
                <DashboardNotificationArea />
              </div>
            </div>
          </div>
        </div>
      </div>
    </>
  );
}
