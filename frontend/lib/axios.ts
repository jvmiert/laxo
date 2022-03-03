import axios, { AxiosInstance } from "axios";
import { GetServerSidePropsContext } from "next";

export const AxiosClient = (locale: string | undefined): AxiosInstance => {
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

export const AxiosServerClient = (
  ctx: GetServerSidePropsContext,
): AxiosInstance => {
  const axiosClient = axios.create({
    baseURL: "http://127.0.0.1:8080/api/",
    withCredentials: true,
    timeout: 1000,
    headers: {
      "Content-Type": "application/json",
      locale: ctx.locale ? ctx.locale : "en",
    },
  });

  if (ctx.req?.headers?.cookie) {
    axiosClient.defaults.headers.post.cookie = ctx.req.headers.cookie;
    axiosClient.defaults.headers.get.cookie = ctx.req.headers.cookie;
  }
  return axiosClient;
};
