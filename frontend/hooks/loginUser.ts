import { useState } from "react";
import AxiosClient from "../lib/axios";

export default function useLoginApi(): [
  isLoading: boolean,
  isError: boolean,
  isSuccess: boolean,
  doLogin: (email: string, password: string) => void,
] {
  const [isLoading, setIsLoading] = useState(false);
  const [isError, setIsError] = useState(false);
  const [isSuccess, setIsSuccess] = useState(false);

  const doLogin = async (email: string, password: string) => {
    setIsError(false);
    setIsLoading(true);
    try {
      const result = await AxiosClient.post("/login", { email, password });
    } catch (error) {
      setIsError(true);
      console.log(error);
    }
    setIsSuccess(true);
    setIsLoading(false);
  };

  return [isLoading, isError, isSuccess, doLogin];
}
