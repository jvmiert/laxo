import type { NextPage } from 'next';
import Head from 'next/head';
import useLoginApi from "../hooks/loginUser";

const Login: NextPage = () => {
  const [isLoading, isError, isSuccess, doLogin] = useLoginApi();

  const submitForm = (e: React.SyntheticEvent) => {
    e.preventDefault();
    doLogin("test@example.com", "12345test");
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
        <form
          onSubmit={submitForm}
          className="bg-white shadow-md rounded px-8 pt-6 pb-8 mb-4"
        >
          <div className="mb-4">
            <label className="block text-gray-700 text-sm font-bold mb-2" htmlFor="username">
            Email
            </label>
            <input className="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline" id="username" type="text" placeholder="Email"/>
          </div>
          <div className="mb-6">
            <label className="block text-gray-700 text-sm font-bold mb-2" htmlFor="password">
            Password
            </label>
            <input className="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 mb-3 leading-tight focus:outline-none focus:shadow-outline" id="password" type="password" placeholder="******************"/>
          </div>
          <div className="flex items-center justify-between">
            <button className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline" type="submit">
            Sign In
            </button>
          </div>
        </form>
      </main>
    </div>
  )
}

export default Login
