import { ReactChildren, ReactNode } from "react";
import { useRouter } from "next/router";
import { AxiosInstance } from "axios";
import AxiosClient from "@/lib/axios";
import createSafeContext from "@/lib/useSafeContext";

export interface AxiosConsumerProps {
  axiosClient: AxiosInstance;
}

export const [useAxios, Provider] = createSafeContext<AxiosConsumerProps>();

export const AxiosProvider = ({
  children,
}: {
  children: ReactChildren | ReactNode;
}) => {
  const { locale } = useRouter();

  const axiosClient = AxiosClient(locale);

  const providerValues: AxiosConsumerProps = {
    axiosClient,
  };

  return <Provider value={providerValues}>{children}</Provider>;
};
