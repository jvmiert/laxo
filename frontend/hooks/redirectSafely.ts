import { useRouter } from "next/router";

export const redirectAllowList = new Set([
  "/dashboard/home",
  "/dashboard/settings",
  "/dashboard/platforms",
  "/dashboard/products",
  "/dashboard/orders",
  "/dashboard/customers",
]);

export const redirectParamAllowList = new Set([]);

export default function useRedirectSafely(): {
  redirectSafely: (url: string) => void;
} {
  const { push, replace, locale = "en" } = useRouter();

  const redirectSafely = (url: string) => {
    const splitURL = url.split("/");
    let checkURL: string;

    if (splitURL.length > 3) {
      checkURL = splitURL.slice(0, 3).join("/");
    } else {
      checkURL = url;
    }

    if (redirectAllowList.has(checkURL)) {
      replace(url, url, { locale: locale });
    } else {
      push("/", "/", { locale: locale });
    }
  };
  return { redirectSafely };
}
