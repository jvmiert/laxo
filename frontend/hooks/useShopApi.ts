import axios from "axios";
import { useSWRConfig } from "swr";
import { useAxios } from "@/providers/AxiosProvider";
import type { ResponseError } from "@/types/ApiResponse";

export default function useShopApi(): {
  doCreateShop: (shopName: string) => Promise<ResponseError>;
} {
  const { axiosClient } = useAxios();
  const { mutate } = useSWRConfig();

  const doCreateShop = async (shopName: string) => {
    try {
      await axiosClient.post("/shop", { shopName });
    } catch (error) {
      if (axios.isAxiosError(error)) {
        mutate("/shop");
        if (error.response?.data instanceof Object) {
          return {
            success: false,
            error: true,
            errorDetails: error.response.data.errorDetails,
          };
        }
        return { success: false, error: true, errorDetails: {} };
      }

      throw error;
    }

    mutate("/shop");
    return { success: true, error: false, errorDetails: {} };
  };

  return { doCreateShop };
}
