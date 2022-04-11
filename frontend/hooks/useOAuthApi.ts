import axios from "axios";
import { useSWRConfig } from "swr";
import { useAxios } from "@/providers/AxiosProvider";
import type { ResponseError } from "@/types/ApiResponse";

export default function useOAuthApi(): {
  doVerifyPlatform: (
    platform: string,
    code: string,
    state?: string,
  ) => Promise<ResponseError>;
} {
  const { axiosClient } = useAxios();
  const { mutate } = useSWRConfig();

  const doVerifyPlatform = async (
    platform: string,
    code: string,
    state?: string,
  ) => {
    try {
      await axiosClient.post("/oauth/verify", { platform, code, state });
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

  return { doVerifyPlatform };
}
