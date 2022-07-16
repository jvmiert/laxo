import { AxiosError, AxiosResponse } from "axios";
import { useMemo, useCallback } from "react";
import useSWR, { KeyedMutator, Fetcher, Key } from "swr";
import useSWRImmutable from "swr/immutable";
import useSWRInfinite, { SWRInfiniteResponse } from "swr/infinite";
import { useAxios } from "@/providers/AxiosProvider";
import type {
  GetShopResponse,
  GetPlatformsResponse,
  GetNotificationResponse,
  NotificationResponseObject,
  GetLazadaPlatformResponse,
  LaxoProductResponse,
  LaxoProduct,
  LaxoProductDetails,
  LaxoProductDetailsResponse,
  LaxoAssetResponse,
  LaxoProductAsset,
} from "@/types/ApiResponse";

export function useGetAuth() {
  const { axiosClient } = useAxios();
  const { data, error } = useSWR("/user", (url) => axiosClient.get(url), {
    shouldRetryOnError: false,
  });

  // @TODO: the data is actually in data.data since we don't process the
  //        JSON anymore.

  const auth = useMemo(() => !!data && !error, [data, error]);

  return {
    auth,
  };
}

export function useGetShop(): {
  shops: GetShopResponse;
  error: AxiosError | undefined;
  loading: boolean;
} {
  const { axiosClient } = useAxios();
  const { data, error, isValidating } = useSWR<
    AxiosResponse<GetShopResponse>,
    AxiosError<unknown>
  >("/shop", (url) => axiosClient.get<GetShopResponse>(url), {
    shouldRetryOnError: true,
  });

  return {
    shops: data ? data.data : { shops: [], total: 0 },
    error,
    loading: isValidating,
  };
}

function transformNotification(r: any): GetNotificationResponse {
  const resp = JSON.parse(r);
  resp.notifications.forEach((n: NotificationResponseObject) => {
    n.notification.created = new Date(n.notification.created);
  });
  return resp;
}

export function useGetNotifications(): {
  notifications: GetNotificationResponse;
  error: AxiosError | undefined;
  loading: boolean;
} {
  const emptyNotifications = useMemo(
    () => ({ notifications: [], total: 0 }),
    [],
  );

  const { axiosClient } = useAxios();
  const { data, error, isValidating } = useSWRImmutable<
    AxiosResponse<GetNotificationResponse>,
    AxiosError<unknown>
  >(
    "/notifications",
    (url) =>
      axiosClient.get<GetNotificationResponse>(url, {
        transformResponse: transformNotification,
      }),
    {
      shouldRetryOnError: false,
    },
  );

  return {
    notifications: data ? data.data : emptyNotifications,
    error,
    loading: isValidating,
  };
}

export function useGetShopPlatformsRedirect(shopID: string): {
  platforms: GetPlatformsResponse;
  error: AxiosError | undefined;
  loading: boolean;
} {
  const { axiosClient } = useAxios();
  const { data, error, isValidating } = useSWR<
    AxiosResponse<GetPlatformsResponse>,
    AxiosError<unknown>
  >(
    shopID !== "" ? ["/oauth/redirects", shopID] : null,
    (url) => axiosClient.get<GetPlatformsResponse>(url, { params: { shopID } }),
    {
      shouldRetryOnError: true,
    },
  );

  return {
    platforms: data
      ? data.data
      : { shopID: "", platforms: [], connectedPlatforms: [] },
    error,
    loading: isValidating,
  };
}

function transformLazadaPlatform(r: any): GetLazadaPlatformResponse {
  const resp = JSON.parse(r);
  resp.refreshExpiresIn = new Date(resp.refreshExpiresIn);
  resp.accessExpiresIn = new Date(resp.accessExpiresIn);
  resp.created = new Date(resp.created);
  return resp;
}

export function useGetLazadaPlatform(): {
  platform: GetLazadaPlatformResponse | undefined;
  error: AxiosError | undefined;
  loading: boolean;
} {
  const { axiosClient } = useAxios();
  const { data, error, isValidating } = useSWR<
    AxiosResponse<GetLazadaPlatformResponse>,
    AxiosError<unknown>
  >(
    "/platforms/lazada",
    (url) =>
      axiosClient.get<GetLazadaPlatformResponse>(url, {
        transformResponse: transformLazadaPlatform,
      }),
    {
      shouldRetryOnError: true,
    },
  );

  return {
    platform: data ? data.data : undefined,
    error,
    loading: isValidating,
  };
}

function transformLaxoProducts(r: any): LaxoProductResponse {
  const resp = JSON.parse(r);
  resp.products.forEach((p: LaxoProduct) => {
    p.product.created = new Date(p.product.created);
    p.product.updated = new Date(p.product.updated);
  });
  return resp;
}

export function useGetLaxoProducts(
  name: string,
  msku: string,
  offset: number,
  limit: number,
): {
  products: LaxoProductResponse;
  error: AxiosError | undefined;
  loading: boolean;
} {
  const { axiosClient } = useAxios();
  const { data, error, isValidating } = useSWR<
    AxiosResponse<LaxoProductResponse>,
    AxiosError<unknown>
  >(
    ["/product", name, msku, offset, limit],
    (url) =>
      axiosClient.get<LaxoProductResponse>(url, {
        params: { name, msku, offset, limit },
        transformResponse: transformLaxoProducts,
      }),
    {
      shouldRetryOnError: true,
    },
  );

  return {
    products: data
      ? data.data
      : { products: [], paginate: { total: 0, pages: 0, limit: 0, offset: 0 } },
    error,
    loading: isValidating,
  };
}

function transformLaxoProductDetails(
  r: any,
): LaxoProductDetailsResponse | undefined {
  let resp;

  try {
    resp = JSON.parse(r);
  } catch {
    return undefined;
  }

  if (resp?.product) {
    resp.product.created = new Date(resp.product.created);
    resp.product.updated = new Date(resp.product.updated);
  }

  return resp;
}

export function useGetLaxoProductDetails(
  productID: string | string[] | undefined,
): {
  product: LaxoProductDetails | undefined;
  error: AxiosError | undefined;
  loading: boolean;
  mutate: KeyedMutator<AxiosResponse<LaxoProductDetailsResponse>>;
} {
  const { axiosClient } = useAxios();
  const { data, error, isValidating, mutate } = useSWR<
    AxiosResponse<LaxoProductDetailsResponse>,
    AxiosError<unknown>
  >(
    productID ? `/product/${productID}` : null,
    (url) =>
      axiosClient.get<LaxoProductDetailsResponse>(url, {
        transformResponse: transformLaxoProductDetails,
      }),
    {
      shouldRetryOnError: true,
    },
  );

  return {
    product: data?.data,
    error,
    loading: isValidating,
    mutate: mutate,
  };
}

function transformLaxoAssets(r: any): LaxoAssetResponse {
  const resp = JSON.parse(r);

  resp.assets.forEach((a: LaxoProductAsset) => {
    a.created = new Date(a.created);
  });

  return resp;
}

export function useGetShopAssets(limit: number): {
  assetsPages: AxiosResponse<LaxoAssetResponse>[] | undefined;
  error: AxiosError | undefined;
  loading: boolean;
  size: number;
  mutate: KeyedMutator<AxiosResponse<LaxoAssetResponse>[]>;
  setSize: (
    size: number | ((_size: number) => number),
  ) => Promise<AxiosResponse<LaxoAssetResponse>[] | undefined>;
} {
  const { axiosClient } = useAxios();

  const getKey = useCallback(
    (
      pageIndex: number,
      previousPageData: AxiosResponse<LaxoAssetResponse>,
    ): Key => {
      if (previousPageData && previousPageData.data.paginate.pages < pageIndex)
        return null; // reached the end
      return `/asset/shop-assets?offset=${pageIndex * limit}&limit=${limit}`; // SWR key
    },
    [limit],
  );

  const fetcher = useCallback(
    (url: string) =>
      axiosClient.get(url, {
        transformResponse: transformLaxoAssets,
      }),
    [axiosClient],
  );

  const { data, error, isValidating, setSize, size, mutate } = useSWRInfinite<
    AxiosResponse<LaxoAssetResponse>,
    AxiosError<Error>
  >(getKey, fetcher, {
    shouldRetryOnError: true,
    revalidateAll: false,
  });

  return useMemo(
    () => ({
      assetsPages: data,
      error,
      loading: isValidating,
      size,
      setSize,
      mutate,
    }),
    [data, error, isValidating, setSize, size, mutate],
  );
}
