import axios from "axios";
import { useAxios } from "@/providers/AxiosProvider";

export interface LoginErrorDetails {
  [key: string]: string;
}

export default function useLoginApi(): [
  doLogin: (
    email: string,
    password: string,
  ) => Promise<{
    success: boolean;
    error: boolean;
    errorDetails: LoginErrorDetails;
  }>,
] {
  const { axiosClient } = useAxios();

  const doLogin = async (email: string, password: string) => {
    try {
      await axiosClient.post("/login", { email, password });
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

  return [doLogin];
}
