import {
  Fragment,
  CSSProperties,
  useMemo,
  useCallback,
  memo,
  useRef,
  useEffect,
  ChangeEvent,
  useState,
} from "react";
import { Dialog, Transition } from "@headlessui/react";
import { XIcon, PlusCircleIcon } from "@heroicons/react/solid";
import { FixedSizeGrid as Grid, areEqual } from "react-window";
import AutoSizer from "react-virtualized-auto-sizer";
import InfiniteLoader from "react-window-infinite-loader";
import { AxiosResponse } from "axios";
import Image from "next/image";
import prettyBytes from "pretty-bytes";
import { Transforms } from "slate";
import MurmurHash3 from "murmurhash3js-revisited";
import { useIntl } from "react-intl";

import { useDashboard } from "@/providers/DashboardProvider";
import useProductApi from "@/hooks/useProductApi";
import { useGetShopAssets } from "@/hooks/swrHooks";
import type { LaxoAssetResponse } from "@/types/ApiResponse";
import { LaxoImageElement } from "@/lib/laxoSlate";
import LoadSpinner from "@/components/LoadSpinner";

const shimmer = `
<svg width="100%" height="100%" version="1.1" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink">
  <defs>
    <linearGradient id="g">
      <stop stop-color="#F3F4F6" offset="20%" />
      <stop stop-color="#F9FAFB" offset="50%" />
      <stop stop-color="#F3F4F6" offset="70%" />
    </linearGradient>
  </defs>
  <rect width="100%" height="100%" fill="#F3F4F6" />
  <rect id="r" width="100%" height="100%" fill="url(#g)" />
  <animate xlink:href="#r" attributeName="x" from="-100%" to="100%" dur="1s" repeatCount="indefinite"  />
</svg>`;

const shimmerBase64 = () =>
  typeof window === "undefined"
    ? Buffer.from(shimmer).toString("base64")
    : window.btoa(shimmer);

type ItemDataShare = {
  isLoaded: (index: number) => boolean;
  data: AxiosResponse<LaxoAssetResponse>[] | undefined;
  colCount: number;
  pageCount: number;
  totalCount: number;
  assetsToken: string;
  onClick: (pageIndex: number, itemIndex: number) => void;
};

const ImageItem = memo(
  ({
    columnIndex,
    rowIndex,
    style,
    data,
  }: {
    columnIndex: number;
    rowIndex: number;
    style: CSSProperties;
    data: ItemDataShare;
  }) => {
    const index = rowIndex * data.colCount + columnIndex;

    if (index > data.totalCount) return <></>;
    if (!data?.data) return <></>;

    const loaded = data.isLoaded(index);
    if (!loaded) {
      return (
        <div style={style} className="p-4">
          <div className="relative flex h-full w-full flex-col">
            <div className="grow">
              <div className="group relative block h-full w-full animate-pulse overflow-hidden rounded-lg bg-gray-100"></div>
            </div>
            <div className="mt-2 block h-4 animate-pulse rounded bg-gray-100 text-sm text-gray-900"></div>
            <div className="mt-2 block h-4 animate-pulse rounded bg-gray-100 text-sm text-gray-900"></div>
          </div>
        </div>
      );
    }

    const pageIndex = Math.floor(index / data.pageCount);
    const pageItemIndex =
      pageIndex == 0 ? index : index - pageIndex * data.pageCount;
    const asset = data.data[pageIndex].data.assets[pageItemIndex];
    return (
      <div style={style} className="p-4">
        <div
          onClick={() => data.onClick(pageIndex, pageItemIndex)}
          className="relative flex h-full w-full cursor-pointer flex-col"
        >
          <div className="grow">
            <div className="group relative block h-full w-full overflow-hidden rounded-lg bg-gray-100">
              <Image
                className="rounded"
                alt=""
                src={`/api/assets/${data.assetsToken}/${asset.id}${asset.extension}`}
                layout="fill"
                placeholder="blur"
                blurDataURL={`data:image/svg+xml;base64,${shimmerBase64()}`}
                objectFit="cover"
                objectPosition="center"
              />
            </div>
          </div>
          <div>
            <p className="pointer-events-none mt-2 block truncate text-sm font-medium text-gray-900">
              {asset.originalFilename}
            </p>
            <p className="pointer-events-none block text-sm font-medium text-gray-500">
              {prettyBytes(asset.fileSize, { locale: "vi" })}
            </p>
          </div>
        </div>
      </div>
    );
  },
  areEqual,
);

ImageItem.displayName = "ImageItem";

// Keeping track of what asset we have already requested from backend
let itemStatusMap = new Map<number, boolean>();

type ProductImageDetailsProps = {};

export default function AssetInsertDialog({}: ProductImageDetailsProps) {
  const {
    activeShop,
    slateEditorRef,
    dashboardDispatch,
    dashboardState: { insertImageIsOpen },
  } = useDashboard();

  const pageCount = useMemo(() => 20, []);

  const t = useIntl();

  const { assetsPages, size, setSize, mutate, loading } =
    useGetShopAssets(pageCount);

  const { doCreateAsset, doUploadAsset, doGetAssetRank } = useProductApi();

  const totalAssets = assetsPages ? assetsPages[0].data.paginate.total : 0;

  const assetCount = assetsPages ? assetsPages[0].data.paginate.total : 0;
  const colCount = 4;
  const rowCount = Math.ceil(assetCount / colCount);

  const inputRef = useRef<HTMLInputElement>(null);
  const gridRef = useRef<Grid | null>(null);

  const [scrollToID, setScrollToID] = useState<string | null>(null);
  const [buttonLoading, setButtonLoading] = useState<boolean>(false);

  useEffect(() => {
    const scrollToData = async () => {
      if (scrollToID && gridRef.current) {
        const index = await doGetAssetRank(scrollToID);
        gridRef.current.scrollToItem({
          align: "center",
          columnIndex: 1,
          rowIndex: Math.floor(index / colCount),
        });
        setScrollToID(null);
      }
    };
    if (gridRef.current && !loading && scrollToID) {
      scrollToData();
    }
  }, [loading, scrollToID, assetsPages, doGetAssetRank, colCount]);

  const loadMoreAssets = (
    startIndex: number,
    stopIndex: number,
  ): Promise<void> => {
    const pageNeeded = Math.ceil(stopIndex / pageCount);
    for (let index = startIndex; index <= stopIndex; index++) {
      itemStatusMap.set(index, true);
    }
    return new Promise(async (resolve) => {
      if (pageNeeded > size) {
        await setSize(pageNeeded);
      }
      resolve();
    });
  };

  const isLoaded = useCallback(
    (index: number): boolean => {
      if (!assetsPages) return false;
      const pageIndex = Math.floor(index / pageCount);
      const pageLength = assetsPages.length;
      const loaded = pageIndex + 1 <= pageLength;
      return loaded;
    },
    [assetsPages, pageCount],
  );

  const isItemRequest = (index: number) => !!itemStatusMap.has(index);

  const onClick = useCallback(
    (pageIndex: number, itemIndex: number) => {
      if (!slateEditorRef.current) return;
      if (!assetsPages) return;
      const asset = assetsPages[pageIndex].data.assets[itemIndex];

      const text = { text: "" };
      const image: LaxoImageElement = {
        type: "image",
        src: `${asset.id}${asset.extension}`,
        height: asset.height,
        width: asset.width,
        children: [text],
      };
      Transforms.insertNodes(slateEditorRef.current, image);

      dashboardDispatch({
        type: "close_image_insert",
      });
    },
    [assetsPages, slateEditorRef, dashboardDispatch],
  );

  const gridItemData: ItemDataShare = useMemo(
    () => ({
      isLoaded,
      colCount,
      data: assetsPages,
      pageCount,
      totalCount: totalAssets - 1,
      assetsToken: activeShop ? activeShop.assetsToken : "",
      onClick: onClick,
    }),
    [assetsPages, isLoaded, pageCount, totalAssets, activeShop, onClick],
  );

  const closeDialog = () => {
    dashboardDispatch({
      type: "close_image_insert",
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
              description: "Asset insert: error upload",
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
              description: "Asset insert: server error upload",
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
        mutate();
        dashboardDispatch({
          type: "alert",
          alert: {
            type: "success",
            message: t.formatMessage({
              description: "Asset insert: successful upload",
              defaultMessage: "Successfully added your new image",
            }),
          },
        });

        setScrollToID(result.asset.id);

        URL.revokeObjectURL(url);
      };
      img.src = url;
    },
    [dashboardDispatch, doCreateAsset, doUploadAsset, t, mutate],
  );

  const onFileChange = (e: ChangeEvent<HTMLInputElement>) => {
    setButtonLoading(true);
    e.preventDefault();
    e.persist();

    if (!e.target?.files) return;

    Array.from(e.target.files).forEach(async (f) => {
      await uploadFile(f);
    });
    setButtonLoading(false);
  };

  const openFileDialog = () => {
    if (inputRef.current) {
      inputRef.current.click();
    }
  };

  return (
    <Transition appear show={!!insertImageIsOpen} as={Fragment}>
      <Dialog as="div" className="relative z-40" onClose={closeDialog}>
        <Transition.Child
          as={Fragment}
          enter="ease-out duration-300"
          enterFrom="opacity-0"
          enterTo="opacity-100"
          leave="ease-in duration-200"
          leaveFrom="opacity-100"
          leaveTo="opacity-0"
        >
          <div className="fixed inset-0 h-screen w-screen bg-zinc-800 bg-opacity-80" />
        </Transition.Child>
        <div className="fixed inset-0 overflow-y-auto">
          <div className="flex h-full items-center justify-center p-14">
            <Transition.Child
              as={Fragment}
              enter="ease-out duration-300"
              enterFrom="opacity-0 scale-95"
              enterTo="opacity-100 scale-100"
              leave="ease-in duration-200"
              leaveFrom="opacity-100 scale-100"
              leaveTo="opacity-0 scale-95"
            >
              <Dialog.Panel className="h-full w-full transform overflow-hidden rounded-lg bg-white transition-all">
                <div className="flex h-full flex-col px-4 py-5 sm:px-6">
                  <div className="flex items-center justify-between">
                    <div>
                      <Dialog.Title className="text-lg font-medium leading-6 text-gray-900">
                        {t.formatMessage({
                          defaultMessage: "Shop Assets",
                          description: "Asset insert: title",
                        })}
                      </Dialog.Title>
                      <p className="mt-1 text-sm text-gray-500">
                        {t.formatMessage({
                          defaultMessage:
                            "Insert one of the images below into your description",
                          description: "Asset insert: description",
                        })}
                      </p>
                    </div>
                    <div className="flex items-center space-x-5">
                      {/*<div className="relative rounded-md border shadow">
                        <div className="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-3">
                          <span className="text-gray-500">
                            <SearchIcon className="h-4 w-4" />
                          </span>
                        </div>
                        <input
                          onChange={() => {}}
                          type="text"
                          className="block w-full rounded-md py-2 pl-9 pr-9 focus:outline-none focus:ring focus:ring-indigo-200"
                          placeholder="Search for image name"
                          defaultValue={""}
                        />
                      </div>*/}
                      <button
                        type="button"
                        disabled={buttonLoading}
                        className="inline-flex min-w-[140px] items-center justify-center rounded-md border border-indigo-500 bg-indigo-500 py-2 px-4 text-white shadow shadow-indigo-500/50 hover:bg-indigo-700 focus:outline-none focus:ring focus:ring-indigo-200 disabled:cursor-not-allowed disabled:hover:bg-indigo-500"
                        onClick={openFileDialog}
                      >
                        {buttonLoading ? (
                          <>
                            <LoadSpinner className="mr-2 h-4 w-4 animate-spin fill-indigo-600 text-gray-200" />
                            {t.formatMessage({
                              defaultMessage: "Loading",
                              description:
                                "Asset insert: add image loading button",
                            })}
                          </>
                        ) : (
                          <>
                            <PlusCircleIcon
                              className="mr-2 -ml-1 h-4 w-4"
                              aria-hidden="true"
                            />
                            {t.formatMessage({
                              defaultMessage: "Add Image",
                              description: "Asset insert: add image button",
                            })}
                          </>
                        )}
                      </button>
                      <button
                        type="button"
                        className="rounded-md bg-white text-gray-400 hover:text-gray-500 focus:ring-2 focus:ring-indigo-500"
                        onClick={closeDialog}
                      >
                        <span className="sr-only">
                          {t.formatMessage({
                            defaultMessage: "Close panel",
                            description: "Asset insert: close panel button",
                          })}
                        </span>
                        <XIcon className="h-6 w-6" aria-hidden="true" />
                      </button>
                    </div>
                  </div>
                  <input
                    multiple
                    accept=".png,.gif,.jpeg,.jpg"
                    className="hidden"
                    ref={inputRef}
                    type="file"
                    onChange={onFileChange}
                  />
                  <div className="my-6 -ml-4 -mr-4 border-b border-gray-200 sm:-ml-6 sm:-mr-6" />
                  <div className="-mr-3 grow sm:-mr-5">
                    <div className="block h-full w-full">
                      {assetsPages && (
                        <AutoSizer>
                          {({ height, width }) => (
                            <InfiniteLoader
                              isItemLoaded={isItemRequest}
                              itemCount={assetCount}
                              loadMoreItems={loadMoreAssets}
                              minimumBatchSize={pageCount}
                            >
                              {({ ref, onItemsRendered }) => {
                                return (
                                  <Grid
                                    ref={(e) => {
                                      ref(e);
                                      gridRef.current = e;
                                    }}
                                    itemData={gridItemData}
                                    columnCount={colCount}
                                    columnWidth={width / colCount - 8}
                                    height={height}
                                    rowCount={rowCount}
                                    rowHeight={420}
                                    width={width}
                                    overscanRowCount={4}
                                    onItemsRendered={({
                                      visibleRowStartIndex,
                                      visibleRowStopIndex,
                                      overscanRowStopIndex,
                                      overscanRowStartIndex,
                                    }) => {
                                      onItemsRendered({
                                        overscanStartIndex:
                                          overscanRowStartIndex * colCount,
                                        overscanStopIndex:
                                          overscanRowStopIndex * colCount,
                                        visibleStartIndex:
                                          visibleRowStartIndex * colCount,
                                        visibleStopIndex:
                                          visibleRowStopIndex * colCount,
                                      });
                                    }}
                                  >
                                    {ImageItem}
                                  </Grid>
                                );
                              }}
                            </InfiniteLoader>
                          )}
                        </AutoSizer>
                      )}
                    </div>
                  </div>
                </div>
              </Dialog.Panel>
            </Transition.Child>
          </div>
        </div>
      </Dialog>
    </Transition>
  );
}
