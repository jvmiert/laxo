import type { ReactNode } from "react";

import Navigation from "@/components/Navigation";
import Footer from "@/components/Footer";

type DefaultLayoutProps = {
  children: ReactNode;
};

export default function DefaultLayout({ children }: DefaultLayoutProps) {
  return (
    <>
      <Navigation />
      <div className="container mx-auto px-4 pt-4">
        <main>{children}</main>
      </div>
      <Footer />
    </>
  );
}
