import { Fragment, CSSProperties, useMemo, useCallback, memo } from "react";
import { Dialog, Transition } from "@headlessui/react";
import { XIcon, PlusCircleIcon, SearchIcon } from "@heroicons/react/solid";
import { FixedSizeGrid as Grid, areEqual } from "react-window";
import AutoSizer from "react-virtualized-auto-sizer";
import InfiniteLoader from "react-window-infinite-loader";
import { AxiosResponse } from "axios";
import Image from "next/image";
import prettyBytes from "pretty-bytes";
import { Transforms } from "slate";

import { useDashboard } from "@/providers/DashboardProvider";
import { useGetShopAssets } from "@/hooks/swrHooks";
import type { LaxoAssetResponse } from "@/types/ApiResponse";
import { LaxoImageElement } from "@/lib/laxoSlate";

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
                alt={"Product preview"}
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

  const { assetsPages, size, setSize } = useGetShopAssets(pageCount);

  const totalAssets = assetsPages ? assetsPages[0].data.paginate.total : 0;

  const assetCount = assetsPages ? assetsPages[0].data.paginate.total : 0;
  const colCount = 4;
  const rowCount = Math.ceil(assetCount / colCount);

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

  return (
    <Transition appear show={!!insertImageIsOpen} as={Fragment}>
      <Dialog as="div" className="relative z-10" onClose={closeDialog}>
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
                        Shop Assets
                      </Dialog.Title>
                      <p className="mt-1 text-sm text-gray-500">
                        Insert one of the images below into your description
                      </p>
                    </div>
                    <div className="flex items-center space-x-5">
                      <div className="relative rounded-md border shadow">
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
                      </div>
                      <button
                        type="button"
                        className="inline-flex items-center rounded-md border border-indigo-500 bg-indigo-500 py-2 px-4 text-white shadow shadow-indigo-500/50 hover:bg-indigo-700 focus:outline-none focus:ring focus:ring-indigo-200 disabled:cursor-not-allowed disabled:bg-indigo-200"
                        onClick={() => {}}
                      >
                        <PlusCircleIcon
                          className="mr-2 -ml-1 h-4 w-4"
                          aria-hidden="true"
                        />
                        Add Image
                      </button>
                      <button
                        type="button"
                        className="rounded-md bg-white text-gray-400 hover:text-gray-500 focus:ring-2 focus:ring-indigo-500"
                        onClick={closeDialog}
                      >
                        <span className="sr-only">Close panel</span>
                        <XIcon className="h-6 w-6" aria-hidden="true" />
                      </button>
                    </div>
                  </div>
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
                              {({ ref, onItemsRendered }) => (
                                <Grid
                                  ref={ref}
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
                              )}
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
