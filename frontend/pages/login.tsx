import { useRouter } from "next/router";
import { Form, Field } from "react-final-form";
import createDecorator from "final-form-focus";
import { useIntl } from "react-intl";
import Head from "next/head";
import loadIntlMessages from "@/helpers/loadIntlMessages";
import type { LoadI18nMessagesProps } from "@/helpers/loadIntlMessages";
import { InferGetStaticPropsType } from "next";
import useLoginFuncs, { LoginSchemaValues } from "@/hooks/loginFormFuncs";

export async function getStaticProps(ctx: LoadI18nMessagesProps) {
  return {
    props: {
      intlMessages: await loadIntlMessages(ctx),
    },
  };
}

type LoginPageProps = InferGetStaticPropsType<typeof getStaticProps>;

const focusOnError = createDecorator<LoginSchemaValues>();

export default function LoginPage(props: LoginPageProps) {
  const router = useRouter();
  const t = useIntl();

  const [validate, submitForm] = useLoginFuncs(router);

  return (
    <div>
      <Head>
        <title>Laxo - Login</title>
        <link rel="icon" href="/favicon.ico" />
      </Head>

      <main>
        <p className="text-3xl font-bold underline">Login</p>
        <Form
          onSubmit={submitForm}
          decorators={[focusOnError]}
          validateOnBlur={false}
          validate={validate}
          initialValues={{ email: "", password: "" }}
          render={({ handleSubmit, submitting, submitError }) => (
            <form
              onSubmit={handleSubmit}
              className="mb-4 rounded bg-white px-8 pt-6 pb-8 shadow-md"
            >
              <div className="mb-4">
                <label
                  className="mb-2 block text-sm font-bold text-gray-700"
                  htmlFor="username"
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
                            defaultMessage: "Email",
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
                  className="focus:shadow-outline rounded bg-blue-500 py-2 px-4 font-bold text-white hover:bg-blue-700 focus:outline-none disabled:cursor-not-allowed disabled:bg-blue-200"
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
      </main>
    </div>
  );
}
