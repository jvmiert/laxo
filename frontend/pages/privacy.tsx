import type { ReactElement } from "react";
import { useRouter } from "next/router";

import DefaultLayout from "@/components/DefaultLayout";
import PrivacyEn from "@/components/terms/PrivacyEn";
import PrivacyVi from "@/components/terms/PrivacyVi";

export default function Privacy() {
  const { locale } = useRouter();

  return (
    <div className="mx-auto max-w-prose space-y-6 py-6">
      {locale === "en" ? <PrivacyEn /> : <PrivacyVi />}
    </div>
  );
}

Privacy.getLayout = function getLayout(page: ReactElement) {
  return <DefaultLayout>{page}</DefaultLayout>;
};
