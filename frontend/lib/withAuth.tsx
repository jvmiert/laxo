import { useEffect } from "react";
import { GetServerSideProps, NextPage } from "next";
import { AxiosServerClient } from "@/lib/axios";
import { useAuth } from "@/providers/AuthProvider";
import { useRouter } from "next/router";

// The GetServerSideProps wrappers in this file will query the backend to see if the user
// is authenticated based on their cookie. Depending on which of the wrapper functions,
// they will redirect to either the login page or the url that is passed to the wrapper.
//
// The wrapper functions will only run during the initial load. Subsequent requests return
// an empty prop and perform no query to backend. We do this because authentication state
// is already known on the client and it would be wasteful to query the backend again.
//
// In order to properly redirect the client on subsequent page loads we have a HOC
// component that will redirect based on the SWR auth state. To make sure this initial
// state is correct, we also pass a fallback key to the prop that indicated the auth state.
// Not passing this fallback will cause the HOC component to redirect because SWR's default
// auth state is not authenticated.

// This wrapper function redirects users that are authenticated
export const withRedirectAuth = (
  url: string,
  getServerSidePropsInner: GetServerSideProps = async () => ({ props: {} }),
) => {
  const getServerSideProps: GetServerSideProps = async (ctx) => {
    const result = await getServerSidePropsInner(ctx);

    let returnProps = {};
    if ("props" in result) {
      returnProps = {
        ...result.props,
      };
    }

    // only run on the server
    if (ctx.req?.url?.indexOf("/_next/data/") === 0) {
      return {
        props: {
          ...returnProps,
        },
      };
    }
    const axios = AxiosServerClient(ctx);
    const { locale } = ctx;

    let redirect;
    try {
      await axios.get("/user");
      redirect = {
        redirect: { destination: `/${locale}/${url}`, permanent: false },
      };
    } catch (error) {
      redirect = {};
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
  getServerSidePropsInner: GetServerSideProps = async () => ({ props: {} }),
) => {
  const getServerSideProps: GetServerSideProps = async (ctx) => {
    const result = await getServerSidePropsInner(ctx);

    let returnProps = {};
    if ("props" in result) {
      returnProps = {
        ...result.props,
      };
    }

    // only run on the server
    if (ctx.req?.url?.indexOf("/_next/data/") === 0) {
      return {
        props: {
          ...returnProps,
        },
      };
    }

    const axios = AxiosServerClient(ctx);
    const { locale, resolvedUrl } = ctx;

    let redirect = {};
    try {
      await axios.get("/user");
    } catch (error) {
      redirect = {
        redirect: {
          destination: `/${locale}/login?next=${resolvedUrl}`,
          permanent: false,
        },
      };
    }

    return {
      ...redirect,
      ...result,
      props: {
        ...returnProps,
        fallback: {
          "/user": true,
        },
      },
    };
  };
  return getServerSideProps;
};

export const withAuthPage = (Page: NextPage) => {
  const WithAuthPage = (props: any) => {
    const { auth } = useAuth();
    const { push, locale, route } = useRouter();

    useEffect(() => {
      if (!auth) {
        push(`/${locale}/login?next=${route}`);
      }
    });
    return <Page {...props} />;
  };

  WithAuthPage.displayname = `WithAuthPage(${Page.displayName})`;
  return WithAuthPage;
};

export const withUnauthPage = (url: string, Page: NextPage) => {
  const WithAuthPage = (props: any) => {
    const { auth } = useAuth();
    const { push, locale } = useRouter();

    useEffect(() => {
      if (auth) {
        push(url, url, { locale: locale });
      }
    });
    return <Page {...props} />;
  };

  WithAuthPage.displayname = `WithUnauthPage(${Page.displayName})`;
  return WithAuthPage;
};
