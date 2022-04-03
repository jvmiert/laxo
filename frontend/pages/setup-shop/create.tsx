import { useRef, useEffect } from "react";
import type { ReactElement } from "react";
import { Form, Field } from "react-final-form";
import createDecorator from "final-form-focus";
import { useIntl } from "react-intl";
import { InferGetServerSidePropsType, GetServerSideProps } from "next";
import { withRedirectUnauth, withAuthPage } from "@/lib/withAuth";
import DefaultLayout from "@/components/DefaultLayout";
import useShopFuncs, { ShopSchemaValues } from "@/hooks/shopFormFuncs";

export const getServerSideProps: GetServerSideProps = withRedirectUnauth();

type CreateStoreProps = InferGetServerSidePropsType<typeof getServerSideProps>;

const focusOnError = createDecorator<ShopSchemaValues>();

function CreateStore(props: CreateStoreProps) {
  const t = useIntl();

  const [validate, submitForm] = useShopFuncs();

  const initialFocusRef = useRef<HTMLInputElement>(null);

  useEffect(() => {
    initialFocusRef.current && initialFocusRef.current.focus();
  }, []);

  return (
    <>
      <h1 className="mb-3 mt-5 text-2xl font-bold">
        {t.formatMessage({
          defaultMessage: "Enter your shop name",
          description: "Setup shop: Form header",
        })}
      </h1>
      <p>
        {t.formatMessage({
          defaultMessage:
            "Enter the name you want your shop to be known as within our platform.",
          description: "Setup shop: Form description",
        })}
      </p>
      <div className="mt-5 max-w-xl">
        <Form
          onSubmit={submitForm}
          decorators={[focusOnError]}
          validateOnBlur={false}
          validate={validate}
          initialValues={{ shopName: "" }}
          render={({ handleSubmit, submitting, submitError }) => (
            <form onSubmit={handleSubmit}>
              <div className="mb-4 min-h-[91px]">
                <label
                  className="mb-2 block text-sm font-bold text-gray-700"
                  htmlFor="shopName"
                >
                  {t.formatMessage({
                    defaultMessage: "Shop Name",
                    description: "Setup shop: Form header",
                  })}
                </label>
                <Field<string>
                  name="shopName"
                  autoFocus
                  render={({ input, meta }) => {
                    const showError =
                      (!meta.pristine || meta.submitFailed) &&
                      meta.touched &&
                      (meta.error ||
                        (meta.submitError && !meta.dirtySinceLastSubmit)) &&
                      !meta.submitting;

                    return (
                      <>
                        <input
                          autoFocus
                          ref={initialFocusRef}
                          className={`appearance-none border shadow ${
                            showError ? "border-red-500" : null
                          } focus:shadow-outline w-full rounded py-2 px-3 leading-tight text-gray-700 focus:outline-none focus:ring focus:ring-indigo-200`}
                          {...input}
                          type="text"
                          placeholder={t.formatMessage({
                            defaultMessage: "My Shop's Name",
                            description: "Shop name placeholder",
                          })}
                        />
                        {showError && (
                          <span className="text-xs italic text-red-500">
                            {meta.error || meta.submitError}
                          </span>
                        )}
                      </>
                    );
                  }}
                />
              </div>
              {submitError && (
                <p className="mb-2 text-xs italic text-red-500">
                  {submitError}
                </p>
              )}
              <div className="flex items-center justify-between">
                <button
                  disabled={submitting}
                  className="w-full rounded-md bg-indigo-500 py-2 px-4 font-bold text-white shadow-lg shadow-indigo-500/50 hover:bg-indigo-700 focus:outline-none focus:ring focus:ring-indigo-200 disabled:cursor-not-allowed disabled:bg-indigo-200"
                  type="submit"
                >
                  {t.formatMessage({
                    defaultMessage: "Continue",
                    description: "Setup shop: Continue Button",
                  })}
                </button>
              </div>
            </form>
          )}
        />
      </div>
    </>
  );
}

CreateStore.getLayout = function getLayout(page: ReactElement) {
  return <DefaultLayout>{page}</DefaultLayout>;
};

export default withAuthPage(CreateStore);
