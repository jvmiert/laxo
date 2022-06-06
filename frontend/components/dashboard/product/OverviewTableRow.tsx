import Image from "next/image";

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
  name: string;
  msku: string;
  sellingPriceInt: number;
  sellingPriceExp: number;
  numberFormat: Intl.NumberFormat;
};

export default function OverviewTableRow({
  imgURL,
  shopToken,
  name,
  msku,
  sellingPriceInt,
  sellingPriceExp,
  numberFormat,
}: OverviewTableRowProps) {
  const shownURL = imgURL.length > 0 ? imgURL[0] : null;
  return (
    <tr>
      <td className="max-w-xs whitespace-nowrap px-6 py-4 text-sm font-medium text-gray-900">
        <div className="relative flex items-center">
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
          <span className="ml-4 truncate">{name}</span>
        </div>
      </td>
      <td className="whitespace-nowrap px-6 py-4 text-sm text-gray-500">
        {msku}
      </td>
      <td className="whitespace-nowrap px-6 py-4 text-sm text-gray-500">
        {numberFormat.format(
          parseFloat(`${sellingPriceInt}e${sellingPriceExp}`),
        )}
      </td>
      <td className="whitespace-nowrap px-6 py-4 text-sm text-gray-500"></td>
      <td className="whitespace-nowrap px-6 py-4 text-right text-sm font-medium">
        <a href="#" className="text-indigo-600 hover:text-indigo-900">
          Edit
        </a>
      </td>
    </tr>
  );
}
