import { AxiosInstance, AxiosPromise } from "axios";
import useSWR from "swr";
import { useAxios } from "@/providers/AxiosProvider";

interface FetcherReturnFunction {
  (url: string): AxiosPromise;
}

const axiosFetcher = (axios: AxiosInstance): FetcherReturnFunction => {
  return (url: string) => axios.get(url).then((res) => res.data);
};

export function useGetAuth() {
  const { axiosClient } = useAxios();
  const { data, error } = useSWR("/user", axiosFetcher(axiosClient), {
    shouldRetryOnError: false,
  });

  return {
    auth: !!data && !error,
  };
}

export function useGetShop(): {
  shops: any;
  error: any;
  loading: boolean;
} {
  const { axiosClient } = useAxios();
  const {
    data: shops,
    error,
    isValidating,
  } = useSWR("/shop", axiosFetcher(axiosClient), {
    shouldRetryOnError: true,
  });

  return {
    shops,
    error,
    loading: isValidating,
  };
}
