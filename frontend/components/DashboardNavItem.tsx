import cc from "classcat";
import Link from "next/link";
import { useRouter } from "next/router";

type DashboardNavItemProps = {
  href: string;
  name: string;
  depth?: number;
};

export default function DashboardNavItem({
  href,
  name,
  depth = 0,
}: DashboardNavItemProps) {
  const { pathname } = useRouter();

  const active = pathname === href;

  return (
    <div
      className={cc([
        "mb-1",
        { "ml-0.5": !active },
        { "border-l-2 border-indigo-700": active },
        `pl-${2 + 4 * depth}`,
        { "hover:ml-0 hover:border-l-2 hover:border-gray-50": !active },
      ])}
    >
      <Link href={href}>
        <a
          className={cc([
            "block rounded",
            { "bg-indigo-50": active },
            "py-2 px-3 font-medium",
            { "text-indigo-700": active },
            { "hover:bg-gray-50": !active },
          ])}
        >
          {name}
        </a>
      </Link>
    </div>
  );
}
