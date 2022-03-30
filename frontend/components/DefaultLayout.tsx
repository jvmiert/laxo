import type { ReactNode } from "react";
import Navigation from "@/components/Navigation";

type DefaultLayoutProps = {
  children: ReactNode;
};

export default function DefaultLayout({ children }: DefaultLayoutProps) {
  return (
    <>
      <Navigation />
      <main>{children}</main>
    </>
  );
}
