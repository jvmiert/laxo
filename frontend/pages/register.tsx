import Head from "next/head";
import { useIntl } from "react-intl";
import { Form, Field } from "react-final-form";
import createDecorator from "final-form-focus";
import Navigation from "@/components/Navigation";
import useRegisterFuncs, {
  RegisterSchemaValues,
} from "@/hooks/registerFormFuncs";
import { InferGetServerSidePropsType, GetServerSideProps } from "next";
import { withRedirectAuth, withUnauthPage } from "@/lib/withAuth";

export const getServerSideProps: GetServerSideProps = withRedirectAuth("/");

type RegisterPageProps = InferGetServerSidePropsType<typeof getServerSideProps>;

const focusOnError = createDecorator<RegisterSchemaValues>();

export default withUnauthPage(
  "/dashboard",
  function RegisterPage(props: RegisterPageProps) {
    const t = useIntl();
    const [validate, submitForm] = useRegisterFuncs();

    return (
      <div>
        <Head>
          <title>Laxo - Register</title>
          <link rel="icon" href="/favicon.ico" />
        </Head>

        <Navigation />
        <main>
          <p className="text-3xl font-bold underline">Register</p>
          Register
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
                    htmlFor="name"
                  >
                    {t.formatMessage({
                      defaultMessage: "Full Name",
                      description: "Register Page: Name Field",
                    })}
                  </label>
                  <Field<string>
                    name="fullname"
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
                              defaultMessage: "Elliot Alderson",
                              description: "Register Page: Name Field",
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
                <div className="mb-4">
                  <label
                    className="mb-2 block text-sm font-bold text-gray-700"
                    htmlFor="email"
                  >
                    {t.formatMessage({
                      defaultMessage: "Email",
                      description: "Register Page: Email Field",
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
                              description: "Register Page: Email Field",
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
                      description: "Register Page: Password Field",
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
                      defaultMessage: "Register",

                      description: "Register Page: Register Button",
                    })}
                  </button>
                </div>
              </form>
            )}
          />
        </main>
      </div>
    );
  },
);
