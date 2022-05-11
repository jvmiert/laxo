import { AxiosError, AxiosResponse } from "axios";
import { useMemo } from "react";
import useSWR from "swr";
import useSWRImmutable from "swr/immutable";
import { useAxios } from "@/providers/AxiosProvider";
import type {
  GetShopResponse,
  GetPlatformsResponse,
  GetNotificationResponse,
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
  >("/notifications", (url) => axiosClient.get<GetNotificationResponse>(url), {
    shouldRetryOnError: false,
  });

  return {
    notifications: data ? data.data : emptyNotifications,
    error,
    loading: isValidating,
  };
}

export function useGetShopPlatforms(shopID: string): {
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
    platforms: data ? data.data : { shopID: "", platforms: [] },
    error,
    loading: isValidating,
  };
}
