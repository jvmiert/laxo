import cc from "classcat";
import prettyBytes from "pretty-bytes";
import Image from "next/image";
import { GripVertical } from "lucide-react";

import { LaxoProductAsset } from "@/types/ApiResponse";

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

type AssetManagementItemProps = {
  asset: LaxoProductAsset;
  assetsToken: string;
  setActiveAssetDetails: (arg: LaxoProductAsset) => void;
  setShowImageDetails: (arg: boolean) => void;
};

export default function AssetManagementItem({
  asset,
  assetsToken,
  setActiveAssetDetails,
  setShowImageDetails,
}: AssetManagementItemProps) {
  const openImageDetails = (asset: LaxoProductAsset) => {
    setActiveAssetDetails(asset);
    setShowImageDetails(true);
  };

  return (
    <li key={asset.id} className="relative">
      <div
        className={cc([
          "group aspect-w-10 aspect-h-7 relative block w-full overflow-hidden rounded-lg bg-gray-100",
          {
            "focus-within:ring-2 focus-within:ring-indigo-500 focus-within:ring-offset-2 focus-within:ring-offset-gray-100":
              true,
          }, //active
          { "ring-2 ring-indigo-500 ring-offset-2": false }, //active
        ])}
      >
        <Image
          className={cc([
            "pointer-events-none rounded",
            { "group-hover:opacity-75": false }, //active
          ])}
          alt=""
          src={`/api/assets/${assetsToken}/${asset.id}${asset.extension}`}
          layout="fill"
          placeholder="blur"
          blurDataURL={`data:image/svg+xml;base64,${shimmerBase64()}`}
          objectFit="cover"
          objectPosition="center"
        />
        <button
          type="button"
          className="absolute inset-0 focus:outline-none"
          onClick={() => openImageDetails(asset)}
        ></button>
      </div>
      <div className="flex w-full flex-nowrap items-center">
        <div className="min-w-0 shrink">
          <p className="pointer-events-none mt-2 block w-full truncate text-sm font-medium text-gray-900">
            {asset.originalFilename}
          </p>
          <p className="pointer-events-none block w-full text-sm font-medium text-gray-500">
            {prettyBytes(asset.fileSize, { locale: "vi" })}
          </p>
        </div>
        <div className="mx-1 grow">
          <button
            type="button"
            className="cursor-grab rounded p-2 hover:bg-gray-100 focus:border-indigo-500 focus:outline-none focus:ring-1 focus:ring-indigo-500"
          >
            <GripVertical className="h-4 w-4" />
          </button>
        </div>
      </div>
    </li>
  );
}
