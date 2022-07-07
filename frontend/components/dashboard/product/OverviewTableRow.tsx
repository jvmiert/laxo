import type { CSSProperties } from "react";
import Image from "next/image";
import Link from "next/link";
import { useIntl } from "react-intl";

import LazadaIcon from "@/components/icons/LazadaIcon";
import ShopeeIcon from "@/components/icons/ShopeeIcon";
import { LaxoProductPlatforms } from "@/types/ApiResponse";

const shimmer = `
<svg width="48px" height="48px" version="1.1" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink">
  <defs>
    <linearGradient id="g">
      <stop stop-color="#E2E8F0" offset="20%" />
      <stop stop-color="#F1F5F9" offset="50%" />
      <stop stop-color="#E2E8F0" offset="70%" />
    </linearGradient>
  </defs>
  <rect width="48px" height="48px" fill="#E2E8F0" />
  <rect id="r" width="48px" height="48px" fill="url(#g)" />
  <animate xlink:href="#r" attributeName="x" from="-48px" to="48px" dur="1s" repeatCount="indefinite"  />
</svg>`;

const shimmerBase64 = () =>
  typeof window === "undefined"
    ? Buffer.from(shimmer).toString("base64")
    : window.btoa(shimmer);

const getPlatformIcon = (platform: string): JSX.Element => {
  switch (platform.toLowerCase()) {
    case "lazada":
      return <LazadaIcon key={platform} className="h-4 w-4" />;
    case "shopee":
      return <ShopeeIcon key={platform} className="h-4 w-4 fill-[#ff5422]" />;
    default:
      return <></>;
  }
};

type OverviewTableRowProps = {
  imgURL: Array<string>;
  shopToken: string;
  id: string;
  name: string;
  msku: string;
  sellingPriceInt: number;
  sellingPriceExp: number;
  numberFormat: Intl.NumberFormat;
  style: CSSProperties;
  platforms: Array<LaxoProductPlatforms>;
};

export default function OverviewTableRow({
  imgURL,
  shopToken,
  id,
  name,
  style,
  msku,
  sellingPriceInt,
  sellingPriceExp,
  numberFormat,
  platforms,
}: OverviewTableRowProps) {
  const t = useIntl();

  const shownURL = imgURL.length > 0 ? imgURL[0] : null;
  return (
    <tr style={style} className="border-none hover:bg-gray-50">
      <td className="max-w-xs whitespace-nowrap p-0 text-sm font-medium text-gray-900">
        <Link href={`/dashboard/products/${id}`}>
          <a>
            <div className="relative flex items-center px-6 py-4">
              <div className="w-8">
                <div className="relative h-8 w-8">
                  {shownURL && (
                    <Image
                      className="rounded"
                      alt=""
                      src={`/api/assets/${shopToken}/${shownURL}`}
                      layout="fill"
                      placeholder="blur"
                      blurDataURL={`data:image/svg+xml;base64,${shimmerBase64()}`}
                      objectFit="contain"
                      objectPosition="center"
                    />
                  )}
                </div>
              </div>
              <span className="ml-4 truncate">{name}</span>
            </div>
          </a>
        </Link>
      </td>
      <td className="whitespace-nowrap p-0 text-sm text-gray-500">
        <Link href={`/dashboard/products/${id}`}>
          <a>
            <div className="px-6 py-4">{msku}</div>
          </a>
        </Link>
      </td>
      <td className="whitespace-nowrap p-0 text-sm text-gray-500">
        <Link href={`/dashboard/products/${id}`}>
          <a>
            <div className="px-6 py-4">
              {numberFormat.format(
                parseFloat(`${sellingPriceInt}e${sellingPriceExp}`),
              )}
            </div>
          </a>
        </Link>
      </td>
      <td className="whitespace-nowrap p-0 text-sm text-gray-500">
        <Link href={`/dashboard/products/${id}`}>
          <a>
            <div className="px-6 py-4">
              {platforms.map((p) => getPlatformIcon(p.platformName))}
            </div>
          </a>
        </Link>
      </td>
      <td className="whitespace-nowrap p-0 text-right text-sm font-medium">
        <Link href={`/dashboard/products/${id}`}>
          <a>
            <div className="px-6 py-4">
              <span className="text-indigo-600 hover:text-indigo-900">
                {t.formatMessage({
                  defaultMessage: "Edit",
                  description: "product overview: row edit button label",
                })}
              </span>
            </div>
          </a>
        </Link>
      </td>
    </tr>
  );
}
