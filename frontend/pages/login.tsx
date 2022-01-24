import { useEffect } from 'react';
import { useRouter } from 'next/router'
import { Form, Field  } from 'react-final-form'
import { ValidationErrors, SubmissionErrors } from "final-form";
import { z } from "zod";
import type { SafeParseReturnType } from "zod";
import type { NextPage } from 'next';
import Head from 'next/head';
import useLoginApi from "../hooks/loginUser";

const loginSchema = z.object({
  email: z.string().email().max(300),
  password: z.string().min(8).max(128),
});

type Values = z.infer<typeof loginSchema>;

const Login: NextPage = () => {
  const router = useRouter()

  const [isLoading, isError, isSuccess, doLogin] = useLoginApi();

  useEffect(() => {
   if(isSuccess) {
     router.push("/")
   }
  }, [isSuccess, router]);

  const submitForm = (values: Values ): SubmissionErrors => {
    //doLogin("test@example.com", "12345test");
    console.log(values)
    return {}
  }

  const validate = (values: Values ): ValidationErrors => {
    const errors: { [key: string | number]: string }= {};

    const validationResult  = loginSchema.safeParse(values);

    if (!validationResult.success) {
      validationResult.error.issues.forEach((validation) => {
        errors[validation.path[0]] = validation.message;
      });
    }
    return errors;
  }

  return (
    <div>
      <Head>
        <title>Laxo - Login</title>
        <link rel="icon" href="/favicon.ico" />
      </Head>

      <main>
        <p className="text-3xl font-bold underline">Login</p>
        <p>Loading: {isLoading.toString()}</p>
        <p>isError: {isError.toString()}</p>
        <p>isSuccess: {isSuccess.toString()}</p>
        <Form
          onSubmit={submitForm}
          validateOnBlur={false}
          validate={validate}
          initialValues={{ email: "", password: "" }}
          render={({ handleSubmit, hasValidationErrors }) => (
            <form
              onSubmit={handleSubmit}
              className="bg-white shadow-md rounded px-8 pt-6 pb-8 mb-4"
            >
              <div className="mb-4">
                <label className="block text-gray-700 text-sm font-bold mb-2" htmlFor="username">
                Email
                </label>
                <Field<string>
                  name="email"
                  render={({ input, meta }) => {
                    const showError = (meta.error || meta.submitError) && meta.touched;
                    return(
                      <>
                        <input className={`shadow appearance-none border ${showError ? "border-red-500" : null} rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline`} {...input} type="text" placeholder="Email"/>
                        {showError && (
                          <span className="text-red-500 text-xs italic">{meta.error || meta.submitError}</span>
                        )}
                        </>
                    )}}
                />
              </div>
              <div className="mb-6">
                <label className="block text-gray-700 text-sm font-bold mb-2" htmlFor="password">
                Password
                </label>
                <Field<string>
                  name="password"
                  render={({ input, meta }) => {
                    const showError = (meta.error || meta.submitError) && meta.touched;
                    return(
                      <>
                        <input className={`shadow appearance-none border ${showError ? "border-red-500" : null} rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline`} {...input} type="password" placeholder="******************"/>
                        {showError && (
                          <span className="text-red-500 text-xs italic">{meta.error || meta.submitError}</span>
                        )}
                        </>
                    )}}
                  />
              </div>
              <div className="flex items-center justify-between">
                <button disabled={hasValidationErrors} className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline" type="submit">
                Sign In
                </button>
              </div>
            </form>
          )}
          />
      </main>
    </div>
  )
}

export default Login
