import Link from "next/link";
import { useRouter } from "next/router";
import type { ReactElement } from "react";
import { FormattedMessage, MessageDescriptor } from "react-intl";
import { ChevronRightIcon } from "@heroicons/react/solid";

import { useGetLaxoProductDetails } from "@/hooks/swrHooks";

function getTranslation(key: string): ReactElement<MessageDescriptor> | "" {
  switch (key.toLowerCase()) {
    case "platforms":
      return (
        <FormattedMessage
          defaultMessage="Platforms"
          description="Breadcrumb indicator: platforms"
        />
      );
    case "products":
      return (
        <FormattedMessage
          defaultMessage="Products"
          description="Breadcrumb indicator: products"
        />
      );
    case "home":
      return (
        <FormattedMessage
          defaultMessage="Home"
          description="Breadcrumb indicator: home"
        />
      );
    default:
      return "";
  }
}

function getSecondCrumbTranslation(
  key: string,
): ReactElement<MessageDescriptor> | "" {
  switch (key.toLowerCase()) {
    case "new":
      return (
        <FormattedMessage
          defaultMessage="Add New Product"
          description="Breadcrumb indicator: add product"
        />
      );
    default:
      return "";
  }
}

export default function Breadcrumbs() {
  const {
    pathname,
    query: { productID },
  } = useRouter();

  const splitPath = pathname.split("/");

  const depth = splitPath.length - 2;

  const { product, loading } = useGetLaxoProductDetails(productID);

  const getSecondCrumb = (): string | ReactElement<MessageDescriptor> => {
    if (loading) return "";

    if (product) {
      return product.product.name;
    }

    if (depth > 1) {
      return getSecondCrumbTranslation(splitPath[3]);
    }

    return "";
  };

  return (
    <nav className="flex" aria-label="Breadcrumb">
      <ol role="list" className="flex items-center space-x-4">
        <li>
          <div className="flex items-center">
            {depth > 1 ? (
              <Link href={`/dashboard/${splitPath[2]}`} passHref>
                <a className="mr-4 text-sm font-medium capitalize capitalize text-gray-500 hover:text-gray-700">
                  {getTranslation(splitPath[2])}
                </a>
              </Link>
            ) : (
              <p className="mr-4 text-sm font-medium capitalize capitalize text-gray-500">
                {getTranslation(splitPath[2])}
              </p>
            )}
            {depth > 1 && (
              <>
                <ChevronRightIcon
                  className="h-5 w-5 flex-shrink-0 text-gray-400"
                  aria-hidden="true"
                />

                <p className="ml-4 truncate text-sm font-medium capitalize text-gray-500">
                  {getSecondCrumb()}
                </p>
              </>
            )}
          </div>
        </li>
      </ol>
    </nav>
  );
}
