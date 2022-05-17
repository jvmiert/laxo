import { Fragment } from "React";
import { useGetShop } from "@/hooks/swrHooks";
import DashboardNavItem from "@/components/DashboardNavItem";

interface navData {
  name: string;
  href: string;
}

interface navObject extends navData {
  children?: Array<navData>;
}

const navigationData: Array<navObject> = [
  { name: "Home", href: "/dashboard/home" },
  {
    name: "Settings",
    href: "/dashboard/settings",
    children: [{ name: "Platforms", href: "/dashboard/platforms" }],
  },
];

export default function DashboardNavigation() {
  const { shops } = useGetShop();

  return (
    <div className="w-52">
      <div className="mb-3 border-b border-gray-200 px-4 pb-3">
        <h1 className="mb-0 text-lg font-bold">Dashboard</h1>
        {shops.total > 0 && (
          <span className="text-xs font-semibold text-gray-500">
            {shops.shops[0].name}
          </span>
        )}
      </div>
      <div className="flex w-full items-stretch">
        <div className="block w-full">
          {navigationData.map((item) =>
            item.children ? (
              <Fragment key={item.name}>
                <DashboardNavItem href={item.href} name={item.name} />
                {item.children.map((subItem) => (
                  <DashboardNavItem
                    key={subItem.name}
                    href={subItem.href}
                    name={subItem.name}
                    depth={1}
                  />
                ))}
              </Fragment>
            ) : (
              <DashboardNavItem
                key={item.name}
                href={item.href}
                name={item.name}
              />
            ),
          )}
        </div>
      </div>
    </div>
  );
}
