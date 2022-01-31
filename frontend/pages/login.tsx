import { useRouter } from "next/router";
import { Form, Field } from "react-final-form";
import type { NextPage } from "next";
import Head from "next/head";
import useLoginFuncs from "../hooks/loginFormFuncs";

const Login: NextPage = () => {
  const router = useRouter();

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
          validateOnBlur={false}
          validate={validate}
          initialValues={{ email: "", password: "" }}
          render={({ handleSubmit, submitting, submitError }) => (
            <form
              onSubmit={handleSubmit}
              className="bg-white shadow-md rounded px-8 pt-6 pb-8 mb-4"
            >
              <div className="mb-4">
                <label
                  className="block text-gray-700 text-sm font-bold mb-2"
                  htmlFor="username"
                >
                  Email
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
                          className={`shadow appearance-none border ${
                            showError ? "border-red-500" : null
                          } rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline`}
                          {...input}
                          type="text"
                          placeholder="Email"
                        />
                        {showError && (
                          <span className="text-red-500 text-xs italic">
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
                  className="block text-gray-700 text-sm font-bold mb-2"
                  htmlFor="password"
                >
                  Password
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
                          className={`shadow appearance-none border ${
                            showError ? "border-red-500" : null
                          } rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline`}
                          {...input}
                          type="password"
                          placeholder="******************"
                        />
                        {showError && (
                          <span className="text-red-500 text-xs italic">
                            {meta.error || meta.submitError}
                          </span>
                        )}
                      </>
                    );
                  }}
                />
              </div>
              {submitError && (
                <p className="text-red-500 text-xs italic mb-2">
                  {submitError}
                </p>
              )}
              <div className="flex items-center justify-between">
                <button
                  disabled={submitting}
                  className="disabled:cursor-not-allowed disabled:bg-blue-200 bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline"
                  type="submit"
                >
                  Sign In
                </button>
              </div>
            </form>
          )}
        />
      </main>
    </div>
  );
};

export default Login;
