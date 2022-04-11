import { useEffect } from "react";
import type { ReactElement } from "react";
import { useIntl } from "react-intl";
import DefaultLayout from "@/components/DefaultLayout";
import LoadSpinner from "@/components/LoadSpinner";
import useOAuthApi from "@/hooks/useOAuthApi";
import { useRouter } from "next/router";

function PlatformCallbackPage() {
  const t = useIntl();
  const { query } = useRouter();

  const { doVerifyPlatform } = useOAuthApi();

  useEffect(() => {
    const { platform, code, state } = query;

    const stringState = state ? state.toString() : undefined;

    if (platform && code) {
      doVerifyPlatform(platform.toString(), code.toString(), stringState);
    }
  }, [doVerifyPlatform, query]);

  return (
    <div className="flex flex-col items-center space-y-2">
      <p>
        {t.formatMessage({
          defaultMessage:
            "Please wait a moment while we connect this platform to your store.",
          description: "Platform connect page: standby loading message",
        })}
      </p>
      <LoadSpinner />
    </div>
  );
}

PlatformCallbackPage.getLayout = function getLayout(page: ReactElement) {
  return <DefaultLayout>{page}</DefaultLayout>;
};

export default PlatformCallbackPage;
