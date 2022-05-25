import { AxiosError, AxiosResponse } from "axios";
import { useMemo } from "react";
import useSWR from "swr";
import useSWRImmutable from "swr/immutable";
import { useAxios } from "@/providers/AxiosProvider";
import type {
  GetShopResponse,
  GetPlatformsResponse,
  GetNotificationResponse,
  NotificationResponseObject,
  GetLazadaPlatformResponse,
  LaxoProductResponse,
  LaxoProduct,
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
    platforms: data ? data.data : <GetPlatformsResponse>{},
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
    ["/product", offset, limit],
    (url) =>
      axiosClient.get<LaxoProductResponse>(url, {
        params: { offset, limit },
        transformResponse: transformLaxoProducts,
      }),
    {
      shouldRetryOnError: true,
    },
  );

  return {
    products: data ? data.data : <LaxoProductResponse>{},
    error,
    loading: isValidating,
  };
}
