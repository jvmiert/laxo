import cc from "classcat";
import type { ReactElement } from "react";
import { FormattedMessage, MessageDescriptor } from "react-intl";
import { useRouter } from "next/router";
import { fromUnixTime, formatDistance } from "date-fns";
import { enUS, vi } from "date-fns/locale";
import { useGetShop, useGetShopPlatforms } from "@/hooks/swrHooks";
import { DotsVerticalIcon } from "@heroicons/react/outline";
import LazadaIcon from "@/components/icons/LazadaIcon";
import ShopeeIcon from "@/components/icons/ShopeeIcon";
import TikiIcon from "@/components/icons/TikiIcon";
import ShopNotMadeNotification from "@/components/ShopNotMadeNotification";

type PlatformVisualDataType = {
  icon: ReactElement;
  text: ReactElement<MessageDescriptor>;
  cssGradient: string;
  cssShadow: string;
  cssHover: string;
  cssRing: string;
};

type PlatformVisualObjectType = {
  [index: string]: PlatformVisualDataType;
};
const PlatformVisuals: PlatformVisualObjectType = {
  shopee: {
    icon: <ShopeeIcon />,
    text: (
      <FormattedMessage
        defaultMessage="Shopee"
        description="Connect Page: Connect Shopee Button"
      />
    ),
    cssGradient: "from-[#ff9c68] to-[#ff5422]",
    cssShadow: "shadow-orange-500/50",
    cssHover: "hover:from-orange-700 hover:to-orange-700",
    cssRing: "focus:ring-orange-200",
  },
  lazada: {
    icon: <LazadaIcon />,
    text: (
      <FormattedMessage
        defaultMessage="Lazada"
        description="Connect Page: Connect Lazada Button"
      />
    ),
    cssGradient: "from-[#37D8FF] to-[#972BFF]",
    cssShadow: "shadow-blue-600/50",
    cssHover: "hover:from-blue-700 hover:to-blue-700",
    cssRing: "focus:ring-blue-200",
  },
  tiki: {
    icon: <TikiIcon />,
    text: (
      <FormattedMessage
        defaultMessage="Tiki"
        description="Connect Page: Connect Tiki Button"
      />
    ),
    cssGradient: "from-[#1a94ff] to-[#1083e8]",
    cssShadow: "shadow-blue-600/50",
    cssHover: "hover:from-blue-700 hover:to-blue-700",
    cssRing: "focus:ring-blue-200",
  },
};

export default function PlatformConnect() {
  const { locale } = useRouter();
  const { shops } = useGetShop();

  const shopID = shops.shops.length > 0 ? shops.shops[0].id : "";

  const { platforms } = useGetShopPlatforms(shopID);

  if (shops.total < 1) return <ShopNotMadeNotification />;

  return (
    <>
      <div className="mt-6 max-w-xl rounded-md border border-gray-100 bg-gray-50 p-6">
        <p className="mb-4 border-b border-gray-200 pb-4 font-bold">Add New</p>
        <div className="flex flex-wrap justify-center gap-4">
          {platforms.platforms.map((p) => {
            if (p.platform in PlatformVisuals)
              return (
                <a
                  key={p.platform}
                  href={p.url}
                  className={cc([
                    "flex w-40 justify-center rounded-md bg-gradient-to-r",
                    PlatformVisuals[p.platform].cssGradient,
                    "py-2 px-4 font-bold text-white shadow-lg",
                    PlatformVisuals[p.platform].cssShadow,
                    PlatformVisuals[p.platform].cssHover,
                    "focus:outline-none focus:ring",
                    PlatformVisuals[p.platform].cssRing,
                  ])}
                  type="submit"
                >
                  {PlatformVisuals[p.platform].icon}{" "}
                  {PlatformVisuals[p.platform].text}
                </a>
              );
          })}
        </div>
      </div>
      {shops.shops.length > 0 && shops.shops[0].platforms.length > 0 && (
        <div className="mt-6 max-w-xl rounded-md border border-gray-100 p-6">
          {shops?.shops[0].platforms.map((platform) => (
            <div key={platform.id} className="flex flex-row justify-between">
              <p className="capitalize">{platform.name}</p>
              <div>
                <p>
                  {formatDistance(fromUnixTime(platform.created), new Date(), {
                    addSuffix: true,
                    locale: locale == "vi" ? vi : enUS,
                  })}
                  <DotsVerticalIcon className="ml-4 inline h-4 w-4" />
                </p>
              </div>
            </div>
          ))}
        </div>
      )}
    </>
  );
}
