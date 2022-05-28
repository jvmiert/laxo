import { useEffect, useState } from "react";
import type { ReactElement } from "react";
import { useIntl } from "react-intl";
import DefaultLayout from "@/components/DefaultLayout";
import LoadSpinner from "@/components/LoadSpinner";
import useOAuthApi from "@/hooks/useOAuthApi";
import { useRouter } from "next/router";
import { XCircleIcon } from "@heroicons/react/solid";

function PlatformCallbackPage() {
  const [error, setError] = useState(false);

  const t = useIntl();
  const { query, push } = useRouter();

  const { doVerifyPlatform } = useOAuthApi();

  useEffect(() => {
    const { platform, code, state } = query;

    const stringState = state ? state.toString() : undefined;

    const verifyFunc = async () => {
      if (platform && code) {
        const result = await doVerifyPlatform(
          platform.toString(),
          code.toString(),
          stringState,
        );

        if (result.error) {
          setError(true);
        }

        if (result.success) {
          if (platform == "lazada") {
            push("/dashboard/platforms/lazada");
          }
        }
      }
    };

    verifyFunc();
  }, [doVerifyPlatform, query, push]);

  return (
    <div className="flex flex-col items-center space-y-2">
      {error ? (
        <div className="flex rounded-xl bg-red-100 p-6">
          <XCircleIcon className="mr-6 h-6 w-6 fill-red-600" />
          <div>
            <h3 className="align-top font-bold text-red-600">
              {t.formatMessage({
                defaultMessage: "Sorry something went wrong",
                description: "Platform connect page: error title",
              })}
            </h3>
            <p className="text-red-600">
              {t.formatMessage({
                defaultMessage:
                  "We are looking into this problem. Please try again later.",
                description: "Platform connect page: error message",
              })}
            </p>
          </div>
        </div>
      ) : (
        <>
          <p>
            {t.formatMessage({
              defaultMessage:
                "Please wait a moment while we connect this platform to your store.",
              description: "Platform connect page: standby loading message",
            })}
          </p>
          <LoadSpinner />
        </>
      )}
    </div>
  );
}

PlatformCallbackPage.getLayout = function getLayout(page: ReactElement) {
  return <DefaultLayout>{page}</DefaultLayout>;
};

export default PlatformCallbackPage;
