import cc from "classcat";
import Link from "next/link";
import { useRouter } from "next/router";

import { getDashboardActiveItem } from "@/lib/getActiveNavigation";

type DashboardNavItemProps = {
  href: string;
  name: JSX.Element;
  depth?: number;
  icon?: JSX.Element;
};

export default function DashboardNavItem({
  href,
  name,
  depth = 0,
  icon,
}: DashboardNavItemProps) {
  const { pathname } = useRouter();

  const active = getDashboardActiveItem(pathname, href);

  return (
    <Link href={href}>
      <a href={href}>
        <div
          className={cc([
            "ml-4 cursor-pointer rounded-l border-l-2 py-2 px-3",
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
            <div className="ml-4 pr-3">
              <div className="h-4 w-4">{icon && icon}</div>
            </div>
            <div
              className={cc([
                "font-medium",
                { "pl-2": depth === 0 },
                { "pl-5": depth > 0 },
              ])}
            >
              {name}
            </div>
          </div>
        </div>
      </a>
    </Link>
  );
}
