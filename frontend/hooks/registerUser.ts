import axios from "axios";
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

  const doRegistration = async (
    email: string,
    password: string,
    fullname: string,
  ) => {
    try {
      await axiosClient.post("/user", { email, password, fullname });
    } catch (error) {
      if (axios.isAxiosError(error)) {
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

    return { success: true, error: false, errorDetails: {} };
  };

  return [doRegistration];
}
