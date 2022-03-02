import axios from "axios";
import { useSWRConfig } from "swr";
import { useAxios } from "@/providers/AxiosProvider";
import type { ResponseError } from "@/types/ApiResponse";

export interface RegisterErrorDetails {
  [key: string]: string;
}

export default function useRegisterApi(): [
  doLogin: (
    email: string,
    password: string,
    fullname: string,
  ) => Promise<ResponseError>,
] {
  const { axiosClient } = useAxios();
  const { mutate } = useSWRConfig();

  const doRegistration = async (
    email: string,
    password: string,
    fullname: string,
  ) => {
    try {
      await axiosClient.post("/user", { email, password, fullname });
    } catch (error) {
      if (axios.isAxiosError(error)) {
        mutate("/user");
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

    mutate("/user");
    return { success: true, error: false, errorDetails: {} };
  };

  return [doRegistration];
}
