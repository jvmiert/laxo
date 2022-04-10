import { useRouter } from "next/router";

export const redirectAllowList = new Set([
  "/dashboard/home",
  "/dashboard/settings",
  "/dashboard/platforms",
]);

export default function useRedirectSafely(): {
  redirectSafely: (url: string) => void;
} {
  const { push, locale = "en" } = useRouter();

  const redirectSafely = (url: string) => {
    if (redirectAllowList.has(url)) {
      push(url, url, { locale: locale });
    } else {
      push("/", "/", { locale: locale });
    }
  };
  return { redirectSafely };
}
