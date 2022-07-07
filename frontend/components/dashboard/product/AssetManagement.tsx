import cc from "classcat";
import prettyBytes from "pretty-bytes";
import { useEffect, useRef, useState, ChangeEvent, useCallback } from "react";
import Image from "next/image";
import { CloudUploadIcon } from "@heroicons/react/outline";
import MurmurHash3 from "murmurhash3js-revisited";
import { useIntl } from "react-intl";

import { LaxoProductAsset } from "@/types/ApiResponse";
import ProductImageDetails from "@/components/dashboard/product/ProductImageDetails";
import { useGetLaxoProductDetails, useGetShopAssets } from "@/hooks/swrHooks";
import { useDashboard } from "@/providers/DashboardProvider";
import useProductApi from "@/hooks/useProductApi";

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

type AssetManagementProps = {
  productID: string;
  mediaList: LaxoProductAsset[];
};

export default function AssetManagement({
  productID,
  mediaList,
}: AssetManagementProps) {
  const t = useIntl();
  const [showImageDetails, setShowImageDetails] = useState(false);
  const [dragActive, setDragActive] = useState(false);
  const [activeAssetDetails, setActiveAssetDetails] =
    useState<LaxoProductAsset>();
  const { activeShop, dashboardDispatch } = useDashboard();

  const dropRef = useRef<HTMLDivElement>(null);
  const inputRef = useRef<HTMLInputElement>(null);

  const { mutate: detailMutate } = useGetLaxoProductDetails(productID);
  const { mutate: assetsMutate } = useGetShopAssets(20);

  const { doCreateAsset, doUploadAsset, doAssetRequest } = useProductApi();

  const openImageDetails = (asset: LaxoProductAsset) => {
    setActiveAssetDetails(asset);
    setShowImageDetails(true);
  };

  const closeImageDetails = () => {
    setShowImageDetails(false);
    setTimeout(() => setActiveAssetDetails(undefined), 200);
  };

  const preventDefaultFunc = (e: DragEvent) => {
    e.preventDefault();
  };

  const openFileDialog = () => {
    if (inputRef.current) {
      inputRef.current.click();
    }
  };

  const removeAsset = async () => {
    if (!activeAssetDetails) return;

    closeImageDetails();

    const assignResult = await doAssetRequest({
      action: "delete",
      productID: productID,
      assetID: activeAssetDetails.id,
      order: 0,
    });

    if (assignResult.error) {
      dashboardDispatch({
        type: "alert",
        alert: {
          type: "error",
          message: t.formatMessage({
            description: "Asset management remove image server error",
            defaultMessage:
              "Something went wrong while removing your image, try again later",
          }),
        },
      });
      return;
    }

    detailMutate();
    dashboardDispatch({
      type: "alert",
      alert: {
        type: "success",
        message: t.formatMessage({
          description: "Asset management successful removed image",
          defaultMessage: "Successfully removed your image",
        }),
      },
    });
  };

  const uploadFile = useCallback(
    async (f: File) => {
      const split = f.type.split("/");

      const error = split.length >= 2 ? split[0] != "image" : true;
      if (error) {
        dashboardDispatch({
          type: "alert",
          alert: {
            type: "error",
            message: t.formatMessage({
              description: "Asset management error upload",
              defaultMessage: "Please only add images",
            }),
          },
        });
      }

      const url = URL.createObjectURL(f);
      const img = document.createElement("img");

      const arrayBuffer = await f.arrayBuffer();
      const byteArray = new Uint8Array(arrayBuffer);

      const murmur = MurmurHash3.x64.hash128(byteArray);

      const serverError = () => {
        dashboardDispatch({
          type: "alert",
          alert: {
            type: "error",
            message: t.formatMessage({
              description: "Asset management server error upload",
              defaultMessage:
                "Something went wrong, please make sure you're upload a valid image and try again",
            }),
          },
        });
      };

      img.onload = async (ev: Event) => {
        const target = ev.target as HTMLImageElement;
        const result = await doCreateAsset({
          originalName: f.name,
          size: f.size,
          width: target.width,
          height: target.height,
          hash: murmur,
        });

        if (result.error || !result?.asset) {
          serverError();
          return;
        }

        if (result.upload && result?.asset?.id) {
          const uploadResult = await doUploadAsset(result.asset.id, f);
          if (uploadResult.error) {
            serverError();
            return;
          }
        }

        const assignResult = await doAssetRequest({
          action: "active",
          productID: productID,
          assetID: result.asset.id,
          order: 0,
        });

        if (assignResult.error) {
          serverError();
          return;
        }

        detailMutate();
        assetsMutate();
        dashboardDispatch({
          type: "alert",
          alert: {
            type: "success",
            message: t.formatMessage({
              description: "Asset management successful upload",
              defaultMessage: "Successfully added your new image",
            }),
          },
        });

        URL.revokeObjectURL(url);
      };
      img.src = url;
    },
    [
      dashboardDispatch,
      doCreateAsset,
      doUploadAsset,
      doAssetRequest,
      t,
      productID,
      detailMutate,
      assetsMutate,
    ],
  );

  const onFileChange = (e: ChangeEvent<HTMLInputElement>) => {
    e.preventDefault();
    e.persist();

    if (!e.target?.files) return;

    Array.from(e.target.files).forEach(async (f) => {
      uploadFile(f);
    });
  };

  const dragEnterFunc = (e: DragEvent) => {
    const target = e.target as HTMLDivElement;
    e.preventDefault();
    if (target.id == "dropDiv") {
      setDragActive(true);
    }
  };

  const dragLeaveFunc = (e: DragEvent) => {
    const target = e.target as HTMLDivElement;
    e.preventDefault();
    if (target.id == "dropDiv") {
      setDragActive(false);
    }
  };

  const dropFunc = useCallback(
    (e: DragEvent) => {
      e.preventDefault();
      if (e.dataTransfer?.items) {
        for (let i = 0; i < e.dataTransfer.items.length; i++) {
          if (e.dataTransfer.items[i].kind === "file") {
            const file = e.dataTransfer.items[i].getAsFile();
            if (file) uploadFile(file);
          }
        }
      } else if (e.dataTransfer?.files) {
        for (let i = 0; i < e.dataTransfer.files.length; i++) {
          const file = e.dataTransfer.files[i];
          uploadFile(file);
        }
      }
      setDragActive(false);
    },
    [uploadFile],
  );

  useEffect(() => {
    window.addEventListener("dragover", preventDefaultFunc, false);
    window.addEventListener("drop", preventDefaultFunc, false);

    return () => {
      window.removeEventListener("dragover", preventDefaultFunc);
      window.removeEventListener("drop", preventDefaultFunc);
    };
  }, []);

  useEffect(() => {
    let dropRefStored: HTMLDivElement;

    if (dropRef.current) {
      dropRefStored = dropRef.current;
      dropRefStored.addEventListener("dragenter", dragEnterFunc, false);
      dropRefStored.addEventListener("dragleave", dragLeaveFunc, false);
      dropRefStored.addEventListener("dragover", preventDefaultFunc, false);
      dropRefStored.addEventListener("drop", dropFunc, false);
    }

    return () => {
      if (dropRefStored) {
        dropRefStored.removeEventListener("dragenter", dragEnterFunc);
        dropRefStored.removeEventListener("dragleave", dragLeaveFunc);
        dropRefStored.removeEventListener("dragover", preventDefaultFunc);
        dropRefStored.removeEventListener("drop", dropFunc);
      }
    };
  }, [dropRef, dropFunc]);

  if (!activeShop) return <></>;

  return (
    <div>
      <ProductImageDetails
        show={showImageDetails}
        close={closeImageDetails}
        asset={activeAssetDetails}
        assetsToken={activeShop.assetsToken}
        removeAsset={removeAsset}
      />
      <div className="mx-auto max-w-sm">
        <div
          ref={dropRef}
          className={cc([
            "border-grey relative flex cursor-pointer flex-col items-center rounded border px-4 py-8 text-center hover:border-transparent hover:ring-2 hover:ring-indigo-200",
            { "border-dashed border-slate-400": !dragActive },
            { "border-transparent ring-2 ring-indigo-200": dragActive },
          ])}
        >
          <div
            onClick={openFileDialog}
            id="dropDiv"
            className="absolute inset-0 z-10"
          />
          <CloudUploadIcon className="mb-5 h-8 w-8 text-gray-700" />
          <span>
            {t.formatMessage({
              description: "Asset management: dropzone description",
              defaultMessage: "Drop your images or click here to add them",
            })}
          </span>
          <input
            multiple
            accept=".png,.gif,.jpeg,.jpg"
            className="hidden"
            ref={inputRef}
            type="file"
            onChange={onFileChange}
          />
        </div>
      </div>
      <div className="mt-8">
        <ul role="list" className="grid grid-cols-4 gap-x-4 gap-y-8">
          {mediaList.map((m) => (
            <li key={m.id} className="relative">
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
                  src={`/api/assets/${activeShop.assetsToken}/${m.id}${m.extension}`}
                  layout="fill"
                  placeholder="blur"
                  blurDataURL={`data:image/svg+xml;base64,${shimmerBase64()}`}
                  objectFit="cover"
                  objectPosition="center"
                />
                <button
                  type="button"
                  className="absolute inset-0 focus:outline-none"
                  onClick={() => openImageDetails(m)}
                ></button>
              </div>
              <p className="pointer-events-none mt-2 block truncate text-sm font-medium text-gray-900">
                {m.originalFilename}
              </p>
              <p className="pointer-events-none block text-sm font-medium text-gray-500">
                {prettyBytes(m.fileSize, { locale: "vi" })}
              </p>
            </li>
          ))}
        </ul>
      </div>
    </div>
  );
}
