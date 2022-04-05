import { AxiosError, AxiosResponse } from "axios";
import useSWR from "swr";
import { useAxios } from "@/providers/AxiosProvider";
import type { GetShopResponse } from "@/types/ApiResponse";

export function useGetAuth() {
  const { axiosClient } = useAxios();
  const { data, error } = useSWR("/user", (url) => axiosClient.get(url), {
    shouldRetryOnError: false,
  });

  // @TODO: the data is actually in data.data since we don't process the
  //        JSON anymore.

  return {
    auth: !!data && !error,
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
