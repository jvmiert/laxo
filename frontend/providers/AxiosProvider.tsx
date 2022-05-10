import { ReactNode, useMemo } from "react";
import { useRouter } from "next/router";
import { AxiosInstance } from "axios";
import { AxiosClient } from "@/lib/axios";
import createSafeContext from "@/lib/useSafeContext";

export interface AxiosConsumerProps {
  axiosClient: AxiosInstance;
}

export const [useAxios, Provider] = createSafeContext<AxiosConsumerProps>();

export const AxiosProvider = ({ children }: { children: ReactNode }) => {
  const { locale } = useRouter();

  const axiosClient = useMemo(() => AxiosClient(locale), [locale]);

  const providerValues: AxiosConsumerProps = {
    axiosClient,
  };

  return <Provider value={providerValues}>{children}</Provider>;
};
