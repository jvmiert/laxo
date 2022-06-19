import Link from "next/link";
import { Fragment } from "react";
import { FormattedMessage } from "react-intl";
import {
  PlusIcon,
  HomeIcon,
  CollectionIcon,
  CogIcon,
} from "@heroicons/react/solid";

import NavLogo from "@/components/NavLogo";
import { useDashboard } from "@/providers/DashboardProvider";
import DashboardNavItem from "@/components/DashboardNavItem";
import DashboardNavPlatformItem from "@/components/DashboardNavPlatformItem";

interface navData {
  name: JSX.Element;
  href: string;
  icon?: JSX.Element;
}

interface navObject extends navData {
  children?: Array<navData>;
}

const navigationData: Array<navObject> = [
  {
    name: (
      <FormattedMessage
        description="Dashboard navigation home button"
        defaultMessage="Home"
      />
    ),
    icon: <HomeIcon />,
    href: "/dashboard/home",
  },
  {
    name: (
      <FormattedMessage
        description="Dashboard navigation products button"
        defaultMessage="Products"
      />
    ),
    icon: <CollectionIcon />,
    href: "/dashboard/products",
  },
];

export default function DashboardNavigation() {
  const { activeShop } = useDashboard();

  return (
    <div className="fixed left-0 top-0 z-[4] flex h-full w-52 flex-col space-y-3 overflow-y-auto border-r border-gray-200 bg-white py-4">
      <div className="flex items-center px-4 pb-4">
        <Link href="/" passHref>
          <a>
            <NavLogo />
          </a>
        </Link>
      </div>
      <div className="px-4 pb-2">
        <h1 className="mb-0 text-lg font-bold">
          <FormattedMessage
            description="Dashboard navigation title"
            defaultMessage="Dashboard"
          />
        </h1>
        {activeShop && (
          <span className="text-xs font-semibold text-gray-500">
            {activeShop.name}
          </span>
        )}
      </div>
      <hr className="mx-4 border-gray-200" />
      <div className="flex w-full items-stretch">
        <div className="flex w-full flex-col space-y-1 pb-3">
          {navigationData.map((item) =>
            item.children ? (
              <Fragment key={item.href}>
                <DashboardNavItem href={item.href} name={item.name} />
                {item.children.map((subItem) => (
                  <DashboardNavItem
                    key={subItem.href}
                    href={subItem.href}
                    icon={subItem.icon}
                    name={subItem.name}
                    depth={1}
                  />
                ))}
              </Fragment>
            ) : (
              <DashboardNavItem
                key={item.href}
                href={item.href}
                icon={item.icon}
                name={item.name}
              />
            ),
          )}
        </div>
      </div>
      <hr className="mx-4 border-gray-200" />
      <div className="flex w-full items-stretch">
        <div className="flex w-full flex-col space-y-1 pb-3">
          <div className="ml-9 py-2 px-3">
            <FormattedMessage
              description="Dashboard navigation platforms header"
              defaultMessage="Platforms"
            />
          </div>
          {activeShop &&
            activeShop.platforms.map((platform) => (
              <DashboardNavPlatformItem
                key={platform.name}
                platform={platform.name}
              />
            ))}
          <DashboardNavItem
            icon={<PlusIcon />}
            href="/dashboard/platforms"
            name={
              <FormattedMessage
                description="Dashboard navigation add platforms button"
                defaultMessage="Add Platforms"
              />
            }
          />
        </div>
      </div>
      <hr className="mx-4 border-gray-200" />
      <div className="grow" />
      <div className="flex w-full items-stretch">
        <div className="flex w-full flex-col space-y-1">
          <DashboardNavItem
            icon={<CogIcon />}
            href="/dashboard/settings"
            name={
              <FormattedMessage
                description="Dashboard navigation settings button"
                defaultMessage="Settings"
              />
            }
          />
        </div>
      </div>
    </div>
  );
}
