import axios from "axios";
import { useSWRConfig } from "swr";
import { useAxios } from "@/providers/AxiosProvider";
import type { ResponseError } from "@/types/ApiResponse";
import {
  useDashboard,
  InitialDashboardState,
} from "@/providers/DashboardProvider";

export interface LoginErrorDetails {
  [key: string]: string;
}

export default function useLoginApi(): {
  doLogin: (email: string, password: string) => Promise<ResponseError>;
  doLogout: () => Promise<ResponseError>;
} {
  const { axiosClient } = useAxios();
  const { mutate } = useSWRConfig();
  const { dashboardDispatch } = useDashboard();

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

    mutate("/user");
    mutate("/notifications");
    return { success: true, error: false, errorDetails: {} };
  };

  const doLogout = async () => {
    try {
      await axiosClient.post("/logout");
    } catch (error) {
      if (axios.isAxiosError(error)) {
        if (error.response?.data instanceof Object) {
          return {
            success: false,
            error: true,
            errorDetails: {},
          };
        }
        return { success: false, error: true, errorDetails: {} };
      }

      throw error;
    }
    mutate("/user");
    mutate("/notifications");
    dashboardDispatch({
      type: "reset",
      state: InitialDashboardState,
    });
    return { success: true, error: false, errorDetails: {} };
  };

  return { doLogin, doLogout };
}
