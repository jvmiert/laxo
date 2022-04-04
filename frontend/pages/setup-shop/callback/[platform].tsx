import type { ReactElement } from "react";
import DefaultLayout from "@/components/DefaultLayout";
import { useRouter } from "next/router";

function PlatformCallbackPage() {
  const { query } = useRouter();
  const { platform } = query;
  return <div>Callback: {platform}</div>;
}

PlatformCallbackPage.getLayout = function getLayout(page: ReactElement) {
  return <DefaultLayout>{page}</DefaultLayout>;
};

export default PlatformCallbackPage;
