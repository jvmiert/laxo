import type { ReactElement } from "react";
import { Form, Field } from "react-final-form";
import createDecorator from "final-form-focus";
import { useIntl } from "react-intl";
import Head from "next/head";
import DefaultLayout from "@/components/DefaultLayout";
import { withRedirectAuth, withUnauthPage } from "@/lib/withAuth";
import { InferGetServerSidePropsType, GetServerSideProps } from "next";
import useLoginFuncs, { LoginSchemaValues } from "@/hooks/loginFormFuncs";

export const getServerSideProps: GetServerSideProps =
  withRedirectAuth("/dashboard");

type LoginPageProps = InferGetServerSidePropsType<typeof getServerSideProps>;

const focusOnError = createDecorator<LoginSchemaValues>();

function LoginPage(props: LoginPageProps) {
  const t = useIntl();

  const [validate, submitForm] = useLoginFuncs();

  return (
    <>
      <Head>
        <title>Laxo: Sign In</title>
        <link rel="icon" href="/favicon.ico" />
      </Head>

      <div className="mx-auto max-w-xl">
        <Form
          onSubmit={submitForm}
          decorators={[focusOnError]}
          validateOnBlur={false}
          validate={validate}
          initialValues={{ email: "", password: "" }}
          render={({ handleSubmit, submitting, submitError }) => (
            <form
              onSubmit={handleSubmit}
              className="mb-4 rounded bg-white px-14 py-16 shadow-lg shadow-gray-300"
            >
              <p className="pt-3 pb-7 text-2xl font-bold text-gray-700">
                Sign in to your account
              </p>
              <div className="mb-6">
                <label
                  className="mb-2 block text-sm font-bold text-gray-700"
                  htmlFor="email"
                >
                  {t.formatMessage({
                    defaultMessage: "Email",
                    description: "Login Page: Email Field",
                  })}
                </label>
                <Field<string>
                  name="email"
                  render={({ input, meta }) => {
                    const showError =
                      meta.touched &&
                      (meta.error ||
                        (meta.submitError && !meta.dirtySinceLastSubmit)) &&
                      !meta.submitting;

                    return (
                      <>
                        <input
                          className={`appearance-none border shadow ${
                            showError ? "border-red-500" : null
                          } focus:shadow-outline w-full rounded py-2 px-3 leading-tight text-gray-700 focus:outline-none`}
                          {...input}
                          type="text"
                          placeholder={t.formatMessage({
                            defaultMessage: "my@email.com",
                            description: "Login Page: Email Field",
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
              <div className="mb-6">
                <label
                  className="mb-2 block text-sm font-bold text-gray-700"
                  htmlFor="password"
                >
                  {t.formatMessage({
                    defaultMessage: "Password",
                    description: "Login Page: Password Field",
                  })}
                </label>
                <Field<string>
                  name="password"
                  render={({ input, meta }) => {
                    const showError =
                      meta.touched &&
                      (meta.error ||
                        (meta.submitError && !meta.dirtySinceLastSubmit)) &&
                      !meta.submitting;

                    return (
                      <>
                        <input
                          className={`appearance-none border shadow ${
                            showError ? "border-red-500" : null
                          } focus:shadow-outline w-full rounded py-2 px-3 leading-tight text-gray-700 focus:outline-none`}
                          {...input}
                          type="password"
                          placeholder="******************"
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
                  className="focus:shadow-outline w-full rounded bg-blue-500 py-2 px-4 font-bold text-white hover:bg-blue-700 focus:outline-none disabled:cursor-not-allowed disabled:bg-blue-200"
                  type="submit"
                >
                  {t.formatMessage({
                    defaultMessage: "Sign In",

                    description: "Login Page: Sign In Button",
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

LoginPage.getLayout = function getLayout(page: ReactElement) {
  return <DefaultLayout>{page}</DefaultLayout>;
};

export default withUnauthPage("/dashboard", LoginPage);
