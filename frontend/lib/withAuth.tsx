import { useEffect } from "react";
import { GetServerSideProps, NextPage } from "next";
import { useRouter } from "next/router";

import { AxiosServerClient } from "@/lib/axios";
import { NextPageWithLayout } from "@/types/pages";
import { useAuth } from "@/providers/AuthProvider";
import useRedirectSafely from "@/hooks/redirectSafely";

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

export const withAuthPage = (Page: NextPageWithLayout) => {
  const WithAuthPage = (props: any) => {
    const { auth } = useAuth();
    const { push, locale, asPath } = useRouter();

    useEffect(() => {
      if (!auth) {
        push(`/${locale}/login?next=${asPath}`);
      }
    });
    return <Page {...props} />;
  };

  WithAuthPage.displayname = `WithAuthPage(${Page.displayName})`;
  WithAuthPage.getLayout = Page.getLayout;
  return WithAuthPage;
};

export const withUnauthPage = (url: string, Page: NextPageWithLayout) => {
  const WithAuthPage = (props: any) => {
    const { auth } = useAuth();
    const { push, locale, query } = useRouter();

    const { redirectSafely } = useRedirectSafely();

    useEffect(() => {
      if (auth) {
        // @HACK: Sometimes we want to redirect to the page that the user was
        // visiting previously. We should actually have some state that indicates
        // when we are redirecting to avoid hijacking the redirect with the push below.
        if (!query?.next) {
          push(url, url, { locale: locale });
        }

        if (query?.next) {
          redirectSafely(query.next as string);
        }
      }
    }, [auth, locale, push, query, redirectSafely]);
    return <Page {...props} />;
  };

  WithAuthPage.displayname = `WithUnauthPage(${Page.displayName})`;
  WithAuthPage.getLayout = Page.getLayout;
  return WithAuthPage;
};
