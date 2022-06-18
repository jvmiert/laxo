import cc from "classcat";
import { Fragment, forwardRef, ReactNode } from "react";
import type { ReactElement, ChangeEvent } from "react";
import { useIntl, defineMessage } from "react-intl";
import { useRouter } from "next/router";
import DashboardLayout from "@/components/DashboardLayout";
import OverviewTableRow from "@/components/dashboard/product/OverviewTableRow";
import { withRedirectUnauth, withAuthPage } from "@/lib/withAuth";
import { InferGetServerSidePropsType, GetServerSideProps } from "next";
import Link, { LinkProps } from "next/link";
import { useGetLaxoProducts } from "@/hooks/swrHooks";
import { useDashboard } from "@/providers/DashboardProvider";
import useShopApi from "@/hooks/useShopApi";
import { generatePaginateNumbers } from "@/lib/paginate";
import {
  ChevronRightIcon,
  ChevronLeftIcon,
  CheckIcon,
  SelectorIcon,
  SearchIcon,
  RefreshIcon,
} from "@heroicons/react/solid";
import { Menu, Transition } from "@headlessui/react";

export const getServerSideProps: GetServerSideProps = withRedirectUnauth();

type DashboardProductsPageProps = InferGetServerSidePropsType<
  typeof getServerSideProps
>;

type EnhancedLinkProps = {
  children: ReactNode[] | ReactNode;
  className: string;
} & LinkProps;

const LimitLink = forwardRef<HTMLAnchorElement, EnhancedLinkProps>(
  (props, ref) => {
    let { href, children, locale, ...rest } = props;
    return (
      <Link href={href} locale={locale}>
        <a ref={ref} {...rest}>
          {children}
        </a>
      </Link>
    );
  },
);

LimitLink.displayName = "LimitLink";

const limitCount = [
  { limit: 10 },
  { limit: 25 },
  { limit: 50 },
  { limit: 100 },
];

function DashboardProductsPage(props: DashboardProductsPageProps) {
  const t = useIntl();
  const {
    push,
    query: { p: queryPageNumber, l: queryLimitNumber, s: searchQuery },
  } = useRouter();

  const currentPage = Number(queryPageNumber) ? Number(queryPageNumber) : 1;
  const currentLimit = Number(queryLimitNumber) ? Number(queryLimitNumber) : 10;

  const currentsearchQuery = searchQuery ? searchQuery.toString() : "";

  const offset = (currentPage - 1) * currentLimit;

  const { doPlatformSync } = useShopApi();
  const { activeShop } = useDashboard();

  const { products } = useGetLaxoProducts(
    currentsearchQuery,
    currentsearchQuery,
    offset,
    currentLimit,
  );

  if (!activeShop) return <></>;

  const { paginate } = products;
  const maxPage = paginate.pages;

  const handlePlatformSync = async () => {
    //@TODO: - Create user feedback for sync
    //       - Hide behind confirmation dialogue
    const result = await doPlatformSync();
  };

  const pNumbers = generatePaginateNumbers(currentPage, maxPage);

  const numberFormatter = new Intl.NumberFormat("vi-VI", {
    style: "currency",
    currency: "VND",
  });

  const handleSearch = (e: ChangeEvent<HTMLInputElement>) => {
    push(
      {
        pathname: "/dashboard/products",
        query: {
          ...(currentPage > 1 && { p: currentPage }),
          ...(currentLimit > 10 && { l: queryLimitNumber }),
          ...(e.target.value != "" && { s: e.target.value }),
        },
      },
      undefined,
      { shallow: true, scroll: false },
    );
  };

  const getDecreaseParams = () => {
    const searchParams = {
      ...(currentsearchQuery != "" && { s: currentsearchQuery }),
    };

    if (currentPage - 1 < 1) {
      if (currentLimit > 10) {
        return { l: currentLimit, ...searchParams };
      }
      return {};
    }

    if (currentLimit > 10) {
      return { p: currentPage - 1, l: currentLimit, ...searchParams };
    }

    return { p: currentPage - 1, ...searchParams };
  };

  const getIncreaseParams = () => {
    const searchParams = {
      ...(currentsearchQuery != "" && { s: currentsearchQuery }),
    };

    if (currentPage + 1 > maxPage) {
      if (currentLimit > 10) {
        return { p: currentPage, l: currentLimit, ...searchParams };
      }
      return { p: currentPage, ...searchParams };
    }

    if (currentLimit > 10) {
      return { p: currentPage + 1, l: currentLimit, ...searchParams };
    }

    return { p: currentPage + 1, ...searchParams };
  };

  const getLimitParams = (limit: number) => {
    const intendedLimit = limit * currentPage;
    const lastPage = Math.ceil(paginate.total / limit);
    const currentPageCeiling =
      intendedLimit > paginate.total ? lastPage : currentPage;
    return {
      ...(currentPageCeiling > 1 && { p: currentPageCeiling }),
      ...(limit > 10 && { l: limit }),
      ...(currentsearchQuery != "" && { s: currentsearchQuery }),
    };
  };

  return (
    <div className="pb-40">
      <div className="flex items-center justify-between">
        <div className="relative rounded-md border shadow">
          <div className="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-3">
            <span className="text-gray-500">
              <SearchIcon className="h-4 w-4" />
            </span>
          </div>
          <input
            onChange={handleSearch}
            type="text"
            className="block w-full rounded-md py-2 pl-9 pr-9 focus:outline-none focus:ring focus:ring-indigo-200"
            placeholder="Search for product name or SKU"
            defaultValue={currentsearchQuery}
          />
        </div>
        <div>
          <button
            onClick={handlePlatformSync}
            className="inline-flex items-center rounded-md border border-indigo-500 bg-indigo-500 py-2 px-4 text-white shadow shadow-indigo-500/50 hover:bg-indigo-700 focus:outline-none focus:ring focus:ring-indigo-200 disabled:cursor-not-allowed disabled:bg-indigo-200"
          >
            <RefreshIcon className="mr-2 -ml-1 h-4 w-4" />
            {t.formatMessage({
              defaultMessage: "Sync Products",

              description: "Products Page: sync products button",
            })}
          </button>
        </div>
      </div>
      <div className="my-6">
        <div className="flex flex-col">
          <div className="-my-2 overflow-x-auto sm:-mx-6 lg:-mx-8">
            <div className="inline-block min-w-full py-2 align-middle sm:px-6 lg:px-8">
              <div className="overflow-hidden border-b border-gray-200 shadow sm:rounded-lg">
                <table className="min-w-full divide-y divide-gray-200">
                  <thead className="bg-gray-50">
                    <tr>
                      <th
                        scope="col"
                        className="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500"
                      >
                        Name
                      </th>
                      <th
                        scope="col"
                        className="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500"
                      >
                        SKU
                      </th>
                      <th
                        scope="col"
                        className="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500"
                      >
                        Price
                      </th>
                      <th
                        scope="col"
                        className="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500"
                      >
                        Platforms
                      </th>
                      <th scope="col" className="relative px-6 py-3">
                        <span className="sr-only">Edit</span>
                      </th>
                    </tr>
                  </thead>
                  <tbody className="divide-y divide-gray-200 bg-white">
                    {products.products.map((p, i) => (
                      <OverviewTableRow
                        key={p.product.id}
                        imgURL={p.mediaList}
                        shopToken={activeShop.assetsToken}
                        id={p.product.id}
                        name={p.product.name}
                        msku={p.product.msku}
                        sellingPriceInt={p.product.sellingPrice.Int}
                        sellingPriceExp={p.product.sellingPrice.Exp}
                        platforms={p.platforms}
                        numberFormat={numberFormatter}
                        style={
                          i + 1 != products.products.length
                            ? {
                                boxShadow: "inset 0 -1px rgb(244, 244, 244)",
                              }
                            : {}
                        }
                      />
                    ))}
                  </tbody>
                </table>
              </div>
            </div>
          </div>
        </div>
      </div>
      <div className="flex items-center justify-between">
        <nav
          className="relative z-0 inline-flex -space-x-px rounded-md shadow-md"
          aria-label="Pagination"
        >
          <Link
            shallow={true}
            scroll={true}
            href={{
              pathname: "/dashboard/products",
              query: getDecreaseParams(),
            }}
          >
            <a className="cursor-pointer items-center rounded-l-md border bg-white px-2 py-2 text-sm font-medium text-gray-500 hover:bg-gray-50">
              <span className="sr-only">Previous</span>
              <ChevronLeftIcon className="h-5 w-5" />
            </a>
          </Link>
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
              <Link
                key={p}
                shallow={true}
                scroll={true}
                href={{
                  pathname: "/dashboard/products",
                  query: {
                    ...(p > 1 && { p: p }),
                    ...(currentLimit > 10 && { l: queryLimitNumber }),
                    ...(currentsearchQuery != "" && { s: currentsearchQuery }),
                  },
                }}
              >
                <a
                  className={cc([
                    {
                      "box-content w-[2ch] cursor-pointer border bg-white px-4 py-2 text-center text-sm font-medium text-gray-500 hover:bg-gray-50":
                        p != queryPageNumber,
                    },
                    {
                      "z-10 box-content w-[2ch] cursor-pointer border border-indigo-500 bg-indigo-50 px-4 py-2 text-center text-sm font-medium text-indigo-600":
                        p == queryPageNumber || (p == 1 && !queryPageNumber),
                    },
                  ])}
                >
                  {p}
                </a>
              </Link>
            );
          })}
          <Link
            shallow={true}
            scroll={true}
            href={{
              pathname: "/dashboard/products",
              query: getIncreaseParams(),
            }}
          >
            <a className="cursor-pointer items-center rounded-r-md border bg-white px-2 py-2 text-sm font-medium text-gray-500 hover:bg-gray-50">
              <span className="sr-only">Next</span>
              <ChevronRightIcon className="h-5 w-5" />
            </a>
          </Link>
        </nav>
        <div className="flex items-center">
          <span className="pr-3">Results per page: </span>
          <Menu>
            <div className="relative">
              <Menu.Button className="relative box-content w-[3ch] cursor-default rounded-md border bg-white py-2 pl-5 pr-10 text-left text-sm shadow-md focus:outline-none focus-visible:border-indigo-500 focus-visible:ring-2 focus-visible:ring-white focus-visible:ring-opacity-75 focus-visible:ring-offset-2 focus-visible:ring-offset-indigo-300">
                <span className="block truncate">{currentLimit}</span>
                <span className="pointer-events-none absolute inset-y-0 right-0 flex items-center pr-2">
                  <SelectorIcon
                    className="h-5 w-5 text-gray-400"
                    aria-hidden="true"
                  />
                </span>
              </Menu.Button>
              <Transition
                as={Fragment}
                leave="transition ease-in duration-100"
                leaveFrom="opacity-100"
                leaveTo="opacity-0"
              >
                <Menu.Items className="absolute mt-1 max-h-60 w-full overflow-auto rounded-md bg-white py-1 text-base shadow-lg ring-1 ring-black ring-opacity-5 focus:outline-none">
                  {limitCount.map((l, limitIdx) => (
                    <Menu.Item key={limitIdx}>
                      <LimitLink
                        href={{
                          pathname: "/dashboard/products",
                          query: getLimitParams(l.limit),
                        }}
                        className={cc([
                          "block select-none hover:bg-indigo-100",
                          {
                            "font-medium text-indigo-900":
                              l.limit == currentLimit,
                          },
                          {
                            "font-normal text-gray-900":
                              l.limit != currentLimit,
                          },
                        ])}
                      >
                        <div className="relative py-2 pl-10 pr-4">
                          {l.limit}
                          {l.limit == currentLimit ? (
                            <span className="absolute inset-y-0 left-0 flex items-center pl-3 text-sm text-indigo-600">
                              <CheckIcon
                                className="h-5 w-5"
                                aria-hidden="true"
                              />
                            </span>
                          ) : null}
                        </div>
                      </LimitLink>
                    </Menu.Item>
                  ))}
                </Menu.Items>
              </Transition>
            </div>
          </Menu>
        </div>
      </div>
    </div>
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
