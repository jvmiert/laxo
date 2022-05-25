import cc from "classcat";
import { useState, Fragment } from "react";
import type { ReactElement } from "react";
import { useIntl, defineMessage } from "react-intl";
import { useRouter } from "next/router";
import DashboardLayout from "@/components/DashboardLayout";
import { withRedirectUnauth, withAuthPage } from "@/lib/withAuth";
import { InferGetServerSidePropsType, GetServerSideProps } from "next";
import { useGetLaxoProducts } from "@/hooks/swrHooks";
import { generatePaginateNumbers } from "@/lib/paginate";
import {
  ChevronRightIcon,
  ChevronLeftIcon,
  CheckIcon,
  SelectorIcon,
  SearchIcon,
  RefreshIcon,
} from "@heroicons/react/solid";
import { Listbox, Transition } from "@headlessui/react";

export const getServerSideProps: GetServerSideProps = withRedirectUnauth();

type DashboardProductsPageProps = InferGetServerSidePropsType<
  typeof getServerSideProps
>;

const limitCount = [
  { limit: 10 },
  { limit: 25 },
  { limit: 50 },
  { limit: 100 },
];

const generateURL = (l: number, p: number): string => {
  const limit = l != 10 ? `l=${l}` : null;
  const page = p > 1 ? `p=${p}` : null;

  let url = "/dashboard/products";

  if (limit || page) url += "?";
  if (page) url += page;
  if (limit && page) url += "&";
  if (limit) url += limit;

  return url;
};

function DashboardProductsPage(props: DashboardProductsPageProps) {
  const t = useIntl();
  const {
    push,
    query: { p: queryPageNumber, l: queryLimitNumber },
  } = useRouter();
  const { products } = useGetLaxoProducts(0, 50);

  const maxPage = 6;

  const [page, setPage] = useState<number>(
    Number(queryPageNumber) ? Number(queryPageNumber) : 1,
  );

  const pNumbers = generatePaginateNumbers(page, maxPage);

  const startingLimit = limitCount.find(
    (l) => l.limit == Number(queryLimitNumber),
  );
  const [limit, setLimit] = useState<{ limit: number }>(
    startingLimit ? startingLimit : limitCount[0],
  );

  const increasePage = () => {
    if (page + 1 > maxPage) return;
    setPage((p) => {
      const newP = p + 1;
      push(
        {
          pathname: "/dashboard/products",
          query: { p: newP, limit: limit.limit },
        },
        generateURL(limit.limit, newP),
        { shallow: true },
      );
      return newP;
    });
  };

  const decreasePage = () => {
    if (page - 1 < 1) return;
    setPage((p) => {
      const newP = p - 1;
      push(
        {
          pathname: "/dashboard/products",
          query: { p: newP, limit: limit.limit },
        },
        generateURL(limit.limit, newP),
        { shallow: true },
      );
      return newP;
    });
  };

  const handleLimit = (l: { limit: number }) => {
    push(
      {
        pathname: "/dashboard/products",
        query: { p: page, l: l.limit },
      },
      generateURL(l.limit, page),
      { shallow: true },
    );
    setLimit(l);
  };

  return (
    <>
      <div className="mb-5 flex items-center">
        <div className="relative rounded-md border shadow">
          <div className="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-3">
            <span className="text-gray-500">
              <SearchIcon className="h-4 w-4" />
            </span>
          </div>
          <input
            type="text"
            className="block w-full rounded-md py-2 pl-9 focus:outline-none focus:ring focus:ring-indigo-200"
          />
        </div>
        <div>
          <button className="inline-flex items-center rounded-md border border-indigo-500 bg-indigo-500 py-2 px-4 text-white shadow shadow-indigo-500/50 hover:bg-indigo-700 focus:outline-none focus:ring focus:ring-indigo-200 disabled:cursor-not-allowed disabled:bg-indigo-200">
            <RefreshIcon className="mr-2 -ml-1 h-4 w-4" />
            {t.formatMessage({
              defaultMessage: "Sync Products",

              description: "Products Page: sync products button",
            })}
          </button>
        </div>
      </div>
      <div className="flex items-center">
        <nav
          className="relative z-0 inline-flex -space-x-px rounded-md shadow-md"
          aria-label="Pagination"
        >
          <a
            onClick={decreasePage}
            className="cursor-pointer items-center rounded-l-md border bg-white px-2 py-2 text-sm font-medium text-gray-500 hover:bg-gray-50"
          >
            <span className="sr-only">Previous</span>
            <ChevronLeftIcon className="h-5 w-5" />
          </a>
          {pNumbers.map((p, i) => {
            if (p === "...") {
              return (
                <span
                  key={`${p}-${i}`}
                  className="box-content w-[2ch] border bg-white px-4 py-2 text-center text-sm font-medium text-gray-700"
                >
                  ...
                </span>
              );
            }
            return (
              <a
                key={`${p}-{p != page}`}
                className={cc([
                  {
                    "box-content w-[2ch] cursor-pointer border bg-white px-4 py-2 text-center text-sm font-medium text-gray-500 hover:bg-gray-50":
                      p != page,
                  },
                  {
                    "z-10 box-content w-[2ch] cursor-pointer border border-indigo-500 bg-indigo-50 px-4 py-2 text-center text-sm font-medium text-indigo-600":
                      p == page,
                  },
                ])}
              >
                {p}
              </a>
            );
          })}

          <a
            onClick={increasePage}
            className="cursor-pointer items-center rounded-r-md border bg-white px-2 py-2 text-sm font-medium text-gray-500 hover:bg-gray-50"
          >
            <span className="sr-only">Next</span>
            <ChevronRightIcon className="h-5 w-5" />
          </a>
        </nav>
        <div className="flex items-center">
          <span className="pr-3">Results per page: </span>
          <Listbox value={limit} onChange={handleLimit}>
            <div className="relative">
              <Listbox.Button className="relative box-content w-[3ch] cursor-default rounded-md border bg-white py-2 pl-5 pr-10 text-left text-sm shadow-md focus:outline-none focus-visible:border-indigo-500 focus-visible:ring-2 focus-visible:ring-white focus-visible:ring-opacity-75 focus-visible:ring-offset-2 focus-visible:ring-offset-indigo-300">
                <span className="block truncate">{limit.limit}</span>
                <span className="pointer-events-none absolute inset-y-0 right-0 flex items-center pr-2">
                  <SelectorIcon
                    className="h-5 w-5 text-gray-400"
                    aria-hidden="true"
                  />
                </span>
              </Listbox.Button>
              <Transition
                as={Fragment}
                leave="transition ease-in duration-100"
                leaveFrom="opacity-100"
                leaveTo="opacity-0"
              >
                <Listbox.Options className="absolute mt-1 max-h-60 w-full overflow-auto rounded-md bg-white py-1 text-base shadow-lg ring-1 ring-black ring-opacity-5 focus:outline-none">
                  {limitCount.map((l, limitIdx) => (
                    <Listbox.Option
                      key={limitIdx}
                      className={({ active }) =>
                        `relative cursor-default select-none py-2 pl-10 pr-4 ${
                          active
                            ? "bg-indigo-100 text-indigo-900"
                            : "text-gray-900"
                        }`
                      }
                      value={l}
                    >
                      {({ selected }) => (
                        <>
                          <span
                            className={`block truncate text-sm ${
                              selected ? "font-medium" : "font-normal"
                            }`}
                          >
                            {l.limit}
                          </span>
                          {selected ? (
                            <span className="absolute inset-y-0 left-0 flex items-center pl-3 text-sm text-indigo-600">
                              <CheckIcon
                                className="h-5 w-5"
                                aria-hidden="true"
                              />
                            </span>
                          ) : null}
                        </>
                      )}
                    </Listbox.Option>
                  ))}
                </Listbox.Options>
              </Transition>
            </div>
          </Listbox>
        </div>
      </div>
    </>
  );
}

DashboardProductsPage.getLayout = function getLayout(page: ReactElement) {
  return (
    <DashboardLayout
      title={defineMessage({
        description: "Dashboard products title",
        defaultMessage: "Products",
      })}
    >
      {page}
    </DashboardLayout>
  );
};

export default withAuthPage(DashboardProductsPage);
