import type { CSSProperties } from "react";
import Image from "next/image";
import Link from "next/link";

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
}: OverviewTableRowProps) {
  const shownURL = imgURL.length > 0 ? imgURL[0] : null;
  return (
    <tr style={style} className="border-none hover:bg-gray-50">
      <td className="max-w-xs whitespace-nowrap p-0 text-sm font-medium text-gray-900">
        <Link href={`/dashboard/products/${id}`}>
          <a>
            <div className="relative flex items-center px-6 py-4">
              <div className="w-8">
                <div className="relative h-8 w-8">
                  <Image
                    className="rounded"
                    alt={"Product preview"}
                    src={`/api/assets/${shopToken}/products/${shownURL}`}
                    layout="fill"
                    placeholder="blur"
                    blurDataURL={`data:image/svg+xml;base64,${shimmerBase64()}`}
                    objectFit="contain"
                    objectPosition="center"
                  />
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
            <div className="px-6 py-4" />
          </a>
        </Link>
      </td>
      <td className="whitespace-nowrap p-0 text-right text-sm font-medium">
        <Link href={`/dashboard/products/${id}`}>
          <a>
            <div className="px-6 py-4">
              <span className="text-indigo-600 hover:text-indigo-900">
                Edit
              </span>
            </div>
          </a>
        </Link>
      </td>
    </tr>
  );
}
