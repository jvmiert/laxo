import { useRef, useEffect } from "react";
import Head from "next/head";
import Link from "next/link";
import { useIntl } from "react-intl";
import { Form, Field } from "react-final-form";
import createDecorator from "final-form-focus";
import useRegisterFuncs, {
  RegisterSchemaValues,
} from "@/hooks/registerFormFuncs";
import { InferGetServerSidePropsType, GetServerSideProps } from "next";
import { withRedirectAuth, withUnauthPage } from "@/lib/withAuth";
import NavLogo from "@/components/NavLogo";

export const getServerSideProps: GetServerSideProps = withRedirectAuth("/");

type RegisterPageProps = InferGetServerSidePropsType<typeof getServerSideProps>;

const focusOnError = createDecorator<RegisterSchemaValues>();

function RegisterPage(props: RegisterPageProps) {
  const t = useIntl();
  const [validate, submitForm] = useRegisterFuncs();

  const initialFocusRef = useRef<HTMLInputElement>(null);

  useEffect(() => {
    initialFocusRef.current && initialFocusRef.current.focus();
  }, []);

  return (
    <>
      <Head>
        <title>Laxo: Register</title>
        <link rel="icon" href="/favicon.ico" />
      </Head>

      <div className="mx-auto flex h-screen max-w-xl flex-col justify-center">
        <div>
          <div className="mb-6">
            <Link href="/" passHref>
              <span className="cursor-pointer">
                <NavLogo />
              </span>
            </Link>
          </div>
          <Form
            onSubmit={submitForm}
            decorators={[focusOnError]}
            validateOnBlur={false}
            validate={validate}
            initialValues={{ email: "", password: "" }}
            render={({ handleSubmit, submitting, submitError }) => (
              <form
                onSubmit={handleSubmit}
                className="mb-4 rounded-md bg-white px-14 py-16 shadow-lg shadow-gray-300"
              >
                <p className="pt-3 pb-7 text-2xl font-bold text-gray-700">
                  {t.formatMessage({
                    defaultMessage: "Create your account",
                    description: "Register Page: Form header",
                  })}
                </p>
                <div className="mb-4 min-h-[91px]">
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
                    autoFocus
                    name="fullname"
                    render={({ input, meta }) => {
                      const attemped = !meta.pristine || meta.submitFailed;
                      const unchangedAfterSubmit =
                        meta.submitError && !meta.dirtySinceLastSubmit;
                      const showError =
                        attemped &&
                        meta.touched &&
                        (meta.error || unchangedAfterSubmit) &&
                        !meta.submitting;

                      return (
                        <>
                          <input
                            ref={initialFocusRef}
                            className={`appearance-none border shadow ${
                              showError ? "border-red-500" : null
                            } focus:shadow-outline w-full rounded py-2 px-3 leading-tight text-gray-700 focus:outline-none focus:ring focus:ring-indigo-200`}
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
                <div className="mb-4 min-h-[91px]">
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
                        (!meta.pristine || meta.submitFailed) &&
                        meta.touched &&
                        (meta.error ||
                          (meta.submitError && !meta.dirtySinceLastSubmit)) &&
                        !meta.submitting;
                      console.log(meta);
                      return (
                        <>
                          <input
                            className={`appearance-none border shadow ${
                              showError ? "border-red-500" : null
                            } focus:shadow-outline w-full rounded py-2 px-3 leading-tight text-gray-700 focus:outline-none focus:ring focus:ring-indigo-200`}
                            {...input}
                            type="text"
                            placeholder={t.formatMessage({
                              defaultMessage: "my@email.com",
                              description: "Email placeholder",
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
                <div className="mb-4 min-h-[91px]">
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
                        (!meta.pristine || meta.submitFailed) &&
                        meta.touched &&
                        (meta.error ||
                          (meta.submitError && !meta.dirtySinceLastSubmit)) &&
                        !meta.submitting;

                      return (
                        <>
                          <input
                            className={`appearance-none border shadow ${
                              showError ? "border-red-500" : null
                            } focus:shadow-outline w-full rounded py-2 px-3 leading-tight text-gray-700 focus:outline-none focus:ring focus:ring-indigo-200`}
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
                    className="w-full rounded-md bg-indigo-500 py-2 px-4 font-bold text-white shadow-lg shadow-indigo-500/50 hover:bg-indigo-700 focus:outline-none focus:ring focus:ring-indigo-200 disabled:cursor-not-allowed disabled:bg-indigo-200"
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

          <div className="ml-4 pt-6">
            <span>
              {t.formatMessage(
                {
                  defaultMessage: "Have an account? {signIn}",
                  description: "Sign up Page: Log in bottom text",
                },
                {
                  signIn: (
                    <Link href={"/login"} passHref>
                      <a className="cursor-pointer text-indigo-500">Sign in</a>
                    </Link>
                  ),
                },
              )}
            </span>
          </div>
        </div>
      </div>
    </>
  );
}

export default withUnauthPage("/dashboard/home", RegisterPage);
