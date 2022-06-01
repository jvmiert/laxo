import cc from "classcat";
import Link from "next/link";
import { useRouter } from "next/router";
import LazadaIcon from "@/components/icons/LazadaIcon";
import ShopeeIcon from "@/components/icons/ShopeeIcon";
import TikiIcon from "@/components/icons/TikiIcon";

type DashboardNavPlatformItemProps = {
  platform: string;
};

function getIcon(platform: string): JSX.Element {
  switch (platform.toLowerCase()) {
    case "lazada":
      return <LazadaIcon className="h-4 w-4" />;
    case "shopee":
      return <ShopeeIcon className="h-4 w-4 fill-[#ff5422]" />;
    case "tiki":
      return <TikiIcon className="h-4 w-4 fill-[#1083e8]" />;
    default:
      return <></>;
  }
}

export default function DashboardNavPlatformItem({
  platform,
}: DashboardNavPlatformItemProps) {
  const { pathname } = useRouter();

  const href = `/dashboard/platforms/${platform.toLowerCase()}`;

  const active = href === pathname;

  return (
    <Link href={href}>
      <a href={href}>
        <div
          className={cc([
            "cursor-pointer rounded-l border-l-2 py-2 px-3 ml-4",
            {
              "border-indigo-700 bg-indigo-50 text-indigo-700": active,
            },
            {
              "border-transparent hover:border-gray-50 hover:bg-gray-50":
                !active,
            },
          ])}
        >
          <div className="flex items-center">
            <div className="ml-4 pr-3">{getIcon(platform)}</div>
            <div className="pl-2 font-medium capitalize">{platform}</div>
          </div>
        </div>
      </a>
    </Link>
  );
}
