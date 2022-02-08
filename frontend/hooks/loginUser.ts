import axios from "axios";
import { useRouter } from "next/router";
import AxiosClient from "@/lib/axios";

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
  const { locale } = useRouter();
  const doLogin = async (email: string, password: string) => {
    try {
      await AxiosClient(locale).post("/login", { email, password });
    } catch (error) {
      if (axios.isAxiosError(error)) {
        if (error.response?.data instanceof Object) {
          return {
            success: false,
            error: true,
            errorDetails: error.response.data,
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
