import type { ReactElement } from "react";
import { useRouter } from "next/router";

import DefaultLayout from "@/components/DefaultLayout";
import TermsEn from "@/components/terms/TermsEn";
import TermsVi from "@/components/terms/TermsVi";

export default function Terms() {
  const { locale } = useRouter();

  return (
    <div className="mx-auto max-w-prose space-y-6 py-6">
      {locale === "en" ? <TermsEn /> : <TermsVi />}
    </div>
  );
}

Terms.getLayout = function getLayout(page: ReactElement) {
  return <DefaultLayout>{page}</DefaultLayout>;
};
