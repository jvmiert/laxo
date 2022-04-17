import type { ReactNode } from "react";
import Head from "next/head";
import Navigation from "@/components/Navigation";
import DashboardNavigation from "@/components/DashboardNavigation";
import DashboardLoadingEvent from "@/components/DashboardLoadingEvent";

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
        <div className="flex w-full flex-row flex-wrap">
          <div className="flex shrink-0 grow-0 basis-auto">
            <DashboardNavigation />
          </div>
          <div className="ml-6 flex grow flex-col">
            <DashboardLoadingEvent />
            <h1 className="mb-4 text-xl font-semibold">{title}</h1>
            <main className="">{children}</main>
          </div>
        </div>
      </div>
    </>
  );
}
