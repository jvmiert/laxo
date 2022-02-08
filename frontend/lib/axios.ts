import axios, { AxiosInstance } from "axios";

const AxiosClient = (locale: string | undefined): AxiosInstance => {
  // FIXME: make default language configurable?
  const headerLocale = locale ? locale : "en";

  return axios.create({
    baseURL: "/api/",
    withCredentials: true,
    timeout: 1000,
    headers: {
      "Content-Type": "application/json",
      locale: headerLocale,
    },
  });
};

export default AxiosClient;
