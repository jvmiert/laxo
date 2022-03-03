import { GetServerSideProps, NextPage } from "next";
import { AxiosServerClient } from "@/lib/axios";

// This wrapper function redirects users that are authenticated
export const withRedirectAuth = (
  url: string,
  getServerSidePropsInner: GetServerSideProps = async () => ({ props: {} }),
) => {
  const getServerSideProps: GetServerSideProps = async (ctx) => {
    const axios = AxiosServerClient(ctx);

    let redirect;
    try {
      await axios.get("/user");
      redirect = { redirect: { destination: url, permanent: false } };
    } catch (error) {
      redirect = {};
    }

    const result = await getServerSidePropsInner(ctx);

    let returnProps = {};
    if ("props" in result) {
      returnProps = {
        ...result.props,
      };
    }
    return {
      ...redirect,
      ...result,
      props: {
        ...returnProps,
      },
    };
  };
  return getServerSideProps;
};

// This wrapper function redirects users that are unauthenticated
export const withRedirectUnauth = (
  url: string,
  getServerSidePropsInner: GetServerSideProps = async () => ({ props: {} }),
) => {
  const getServerSideProps: GetServerSideProps = async (ctx) => {
    const axios = AxiosServerClient(ctx);

    let redirect = {};
    try {
      await axios.get("/user");
    } catch (error) {
      redirect = { redirect: { destination: url, permanent: false } };
    }

    const result = await getServerSidePropsInner(ctx);

    let returnProps = {};
    if ("props" in result) {
      returnProps = {
        ...result.props,
      };
    }
    return {
      ...redirect,
      ...result,
      props: {
        ...returnProps,
      },
    };
  };
  return getServerSideProps;
};
