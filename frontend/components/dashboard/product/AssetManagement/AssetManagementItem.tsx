import { forwardRef, HTMLAttributes } from "react";
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

export enum Position {
  Before = -1,
  After = 1,
}

export type AssetManagementItemProps = Omit<
  HTMLAttributes<HTMLButtonElement>,
  "id"
> & {
  clone?: boolean;
  active?: boolean;
  asset: LaxoProductAsset;
  assetsToken: string;
  setActiveAssetDetails: (arg: LaxoProductAsset) => void;
  setShowImageDetails: (arg: boolean) => void;
  insertPosition?: Position;
};

export const AssetManagementItem = forwardRef<
  HTMLLIElement,
  AssetManagementItemProps
>(
  (
    {
      asset,
      assetsToken,
      setActiveAssetDetails,
      setShowImageDetails,
      insertPosition,
      active,
      clone,
      ...props
    }: AssetManagementItemProps,
    ref,
  ) => {
    const openImageDetails = (asset: LaxoProductAsset) => {
      setActiveAssetDetails(asset);
      setShowImageDetails(true);
    };

    return (
      <li
        key={asset.id}
        ref={ref}
        className={cc([
          "relative",
          { "cursor-grab shadow-md": clone },
          { "animate-pop": clone },
          {
            "after:absolute after:inset-y-0 after:-left-2 after:w-0.5 after:bg-indigo-600 after:content-['']":
              insertPosition === Position.Before,
          },
          {
            "after:absolute after:inset-y-0 after:-right-2 after:w-0.5 after:bg-indigo-600 after:content-['']":
              insertPosition === Position.After,
          },
        ])}
        style={
          clone ? { transform: "translate3d(10px, 10px, 0) scale(1.025)" } : {}
        }
      >
        <div className="relative overflow-hidden rounded-md bg-white p-2">
          <div
            className={cc([
              { "absolute inset-0 z-20 bg-slate-100": active },
              { hidden: !active },
            ])}
          />
          <div
            className={cc([
              "aspect-w-10 aspect-h-7 relative block w-full overflow-hidden rounded-lg bg-gray-100 focus-within:ring-2 focus-within:ring-indigo-500 focus-within:ring-offset-2 focus-within:ring-offset-gray-100",
            ])}
          >
            <Image
              className="pointer-events-none rounded"
              alt=""
              src={`/api/assets/${assetsToken}/${asset.id}${asset.extension}`}
              layout="fill"
              placeholder={clone ? "empty" : "blur"}
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
          <div className={cc(["flex w-full flex-nowrap items-center"])}>
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
                {...props}
              >
                <GripVertical className="h-4 w-4" />
              </button>
            </div>
          </div>
        </div>
      </li>
    );
  },
);

AssetManagementItem.displayName = "AssetManagementItem";
