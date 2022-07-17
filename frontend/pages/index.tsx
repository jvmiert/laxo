import type { ReactElement } from "react";
import Link from "next/link";
import { useIntl } from "react-intl";
import Head from "next/head";
import { AnnotationIcon } from "@heroicons/react/outline";

import Everywhere from "@/components/landingpage/Everywhere";
import DefaultLayout from "@/components/DefaultLayout";
import { useAuth } from "@/providers/AuthProvider";

export default function HomePage() {
  const t = useIntl();
  const { auth } = useAuth();
  return (
    <>
      <Head>
        <title>Laxo</title>
        <link rel="icon" href="/favicon.ico" />
      </Head>

      <section>
        <div className="md:pb-30 lg:pt-30 pt-16 pb-20 md:pt-20 lg:pb-40">
          <div className="flex flex-wrap items-center justify-between">
            <div className="mt-12 mb-16 flex max-w-2xl flex-col items-start text-left lg:w-1/2 lg:flex-grow">
              <span className="mb-8 text-xs font-bold uppercase tracking-widest text-indigo-600">
                Sell better, sell more.
              </span>
              <h1 className="mb-3 text-4xl font-bold tracking-tighter text-zinc-900 md:text-7xl lg:text-6xl lg:leading-snug">
                Making Your Selling Experience Easier
              </h1>
              <div className="max-w-md">
                <p className="mb-10 text-left text-lg leading-relaxed text-gray-500">
                  Laxo allows you to automatically manage your products over all
                  your sales channels, keep track of your orders, and improves
                  your marketing activities.
                </p>
              </div>
              <h3 className="mb-2 text-base font-medium underline decoration-indigo-400 decoration-2 underline-offset-2">
                Sign up to get notified when it&apos;s ready.
              </h3>
              <div className="flex-col sm:flex">
                <form className="border2 transform rounded-xl bg-gray-50 p-2 transition duration-500 ease-in-out sm:flex sm:max-w-lg">
                  <div className="revue-form-group min-w-0 flex-1">
                    <label htmlFor="email" className="sr-only">
                      Email address
                    </label>
                    <input
                      id="email"
                      type="email"
                      className="block w-full transform rounded-md border border-transparent bg-transparent px-5 py-3 text-base text-neutral-600 placeholder-gray-400 transition duration-500 ease-in-out focus:border-transparent focus:outline-none focus:ring-2 focus:ring-white focus:ring-offset-2 focus:ring-offset-gray-300"
                      placeholder="Enter your email  "
                    />
                  </div>
                  <div className="mt-4 sm:mt-0 sm:ml-3">
                    <button
                      type="submit"
                      value="Subscribe"
                      className="w-full rounded-lg bg-indigo-500 py-3 px-5 font-bold text-white shadow-lg shadow-indigo-500/50 hover:bg-indigo-700 focus:outline-none focus:ring focus:ring-indigo-200 disabled:cursor-not-allowed disabled:bg-indigo-200 sm:px-10"
                    >
                      Notify me
                    </button>
                  </div>
                </form>
                <div className="sm:flex sm:max-w-lg">
                  <p className="mt-3 text-xs text-gray-500">
                    By subscribing, you agree with our{" "}
                    <Link href="/terms">
                      <a className="underline">Terms of Service</a>
                    </Link>{" "}
                    and{" "}
                    <Link href="/privacy">
                      <a className="underline">Privacy Policy</a>
                    </Link>
                    .
                  </p>
                </div>
              </div>
            </div>
            <div className="w-full rounded-xl lg:w-1/2 lg:max-w-lg">
              <div>
                <div className="relative w-full max-w-lg">
                  <div className="absolute top-0 -right-4 h-72 w-72 rounded-full bg-violet-300 opacity-70 mix-blend-multiply blur-xl filter"></div>

                  <div className="absolute left-20 h-72 w-72 rounded-full bg-fuchsia-300 opacity-70 mix-blend-multiply blur-xl filter lg:-bottom-24"></div>
                  <div className="relative"></div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </section>

      <section>
        <div className="bg-white py-10">
          <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
            <div className="lg:text-center">
              <p className="mt-2 text-3xl font-extrabold leading-8 tracking-tight text-gray-900 sm:text-4xl">
                Sell more by selling everywhere
              </p>
              <p className="mt-4 max-w-2xl text-xl text-gray-500 lg:mx-auto">
                Stop selecting your sales channels based on your available
                resources. Manage your products on Laxo and let us take care of
                the rest.
              </p>
            </div>

            <div className="mt-10">
              <dl className="space-y-10 md:grid md:grid-cols-3 md:gap-x-8 md:gap-y-10 md:space-y-0">
                <div className="relative">
                  <dt>
                    <div className="absolute flex h-12 w-12 items-center justify-center rounded-md bg-indigo-500 text-white">
                      <AnnotationIcon className="h-6 w-6" aria-hidden="true" />
                    </div>
                    <p className="ml-16 text-lg font-medium leading-6 text-gray-900">
                      Update your product photos
                    </p>
                  </dt>
                  <dd className="mt-2 ml-16 text-base text-gray-500">
                    dolor sit amet consect adipisicing elit.
                  </dd>
                </div>

                <div className="relative">
                  <dt>
                    <div className="absolute flex h-12 w-12 items-center justify-center rounded-md bg-indigo-500 text-white">
                      <AnnotationIcon className="h-6 w-6" aria-hidden="true" />
                    </div>
                    <p className="ml-16 text-lg font-medium leading-6 text-gray-900">
                      Change your product descriptions
                    </p>
                  </dt>
                  <dd className="mt-2 ml-16 text-base text-gray-500">
                    dolor sit amet consect adipisicing elit.
                  </dd>
                </div>

                <div className="relative">
                  <dt>
                    <div className="absolute flex h-12 w-12 items-center justify-center rounded-md bg-indigo-500 text-white">
                      <AnnotationIcon className="h-6 w-6" aria-hidden="true" />
                    </div>
                    <p className="ml-16 text-lg font-medium leading-6 text-gray-900">
                      Adjust your product prices
                    </p>
                  </dt>
                  <dd className="mt-2 ml-16 text-base text-gray-500">
                    dolor sit amet consect adipisicing elit.
                  </dd>
                </div>
              </dl>
            </div>
          </div>
        </div>
      </section>

      <section>
        <div className="md:py-30 lg:py-40">
          <Everywhere />
        </div>
      </section>

      <section>
        <div className="relative mt-12 sm:mt-16 lg:mt-24">
          <div className="lg:grid lg:grid-flow-row-dense lg:grid-cols-2 lg:items-center lg:gap-8">
            <div className="lg:col-start-2">
              <h3 className="text-2xl font-extrabold tracking-tight text-gray-900 sm:text-3xl">
                Conveniently manage your{" "}
                <span className="whitespace-nowrap">orders & inventory</span>
              </h3>
              <p className="mt-3 text-lg text-gray-500">
                Laxo keeps track of all your sales channels to allow you to
                track, review, and manage everything in one place.
              </p>

              <ul className="mt-10 space-y-10">
                <li>
                  <div className="flex items-center">
                    <div className="flex h-12 w-12 items-center justify-center rounded-md bg-indigo-500 text-white">
                      <AnnotationIcon className="h-6 w-6" aria-hidden="true" />
                    </div>
                    <p className="ml-8 text-lg font-medium text-gray-900">
                      See all your sales in one place
                    </p>
                  </div>
                </li>
                <li>
                  <div className="flex items-center">
                    <div className="flex h-12 w-12 items-center justify-center rounded-md bg-indigo-500 text-white">
                      <AnnotationIcon className="h-6 w-6" aria-hidden="true" />
                    </div>
                    <p className="ml-8 text-lg font-medium text-gray-900">
                      Manage your inventory in one place
                    </p>
                  </div>
                </li>
                <li>
                  <div className="flex items-center">
                    <div className="flex h-12 w-12 items-center justify-center rounded-md bg-indigo-500 text-white">
                      <AnnotationIcon className="h-6 w-6" aria-hidden="true" />
                    </div>
                    <p className="ml-8 text-lg font-medium text-gray-900">
                      Fullfill your orders and returns in once place
                    </p>
                  </div>
                </li>
              </ul>
            </div>

            <div className="relative -mx-4 mt-10 lg:col-start-1 lg:mt-0">
              <svg
                className="relative mx-auto"
                xmlns="http://www.w3.org/2000/svg"
                width="490"
                height="570"
                viewBox="0 0 490 570"
              >
                <rect fill="#ddd" width="490" height="570" />
                <text
                  fill="rgba(0,0,0,0.5)"
                  fontFamily="sans-serif"
                  fontSize="30"
                  dy="10.5"
                  fontWeight="bold"
                  x="50%"
                  y="50%"
                  textAnchor="middle"
                >
                  ...
                </text>
              </svg>
            </div>
          </div>
        </div>
      </section>

      <section>
        <div className="relative mt-12 lg:mt-24 lg:grid lg:grid-cols-2 lg:items-center lg:gap-8">
          <div className="relative">
            <h3 className="text-2xl font-extrabold tracking-tight text-gray-900 sm:text-3xl">
              Create coordinated marketing campaigns
            </h3>
            <p className="mt-3 text-lg text-gray-500">
              Design and create your marketing campaign simply on Laxo and apply
              to your sales channels all at once.
            </p>

            <ul className="mt-10 space-y-10">
              <li>
                <div className="flex items-center">
                  <div className="flex h-12 w-12 items-center justify-center rounded-md bg-indigo-500 text-white">
                    <AnnotationIcon className="h-6 w-6" aria-hidden="true" />
                  </div>
                  <p className="ml-8 text-lg font-medium text-gray-900">
                    Easily create special product image frames
                  </p>
                </div>
              </li>
              <li>
                <div className="flex items-center">
                  <div className="flex h-12 w-12 items-center justify-center rounded-md bg-indigo-500 text-white">
                    <AnnotationIcon className="h-6 w-6" aria-hidden="true" />
                  </div>
                  <p className="ml-8 text-lg font-medium text-gray-900">
                    Easily apply special promotion prices
                  </p>
                </div>
              </li>
              <li>
                <div className="flex items-center">
                  <div className="flex h-12 w-12 items-center justify-center rounded-md bg-indigo-500 text-white">
                    <AnnotationIcon className="h-6 w-6" aria-hidden="true" />
                  </div>
                  <p className="ml-8 text-lg font-medium text-gray-900">
                    Easily distribute discount coupons
                  </p>
                </div>
              </li>
            </ul>
          </div>

          <div className="relative -mx-4 mt-10 lg:mt-0" aria-hidden="true">
            <svg
              className="relative mx-auto"
              xmlns="http://www.w3.org/2000/svg"
              width="490"
              height="570"
              viewBox="0 0 490 570"
            >
              <rect fill="#ddd" width="490" height="570" />
              <text
                fill="rgba(0,0,0,0.5)"
                fontFamily="sans-serif"
                fontSize="30"
                dy="10.5"
                fontWeight="bold"
                x="50%"
                y="50%"
                textAnchor="middle"
              >
                ...
              </text>
            </svg>
          </div>
        </div>
      </section>
    </>
  );
}

HomePage.getLayout = function getLayout(page: ReactElement) {
  return <DefaultLayout>{page}</DefaultLayout>;
};
