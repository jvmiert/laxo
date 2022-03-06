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
