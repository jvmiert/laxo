import LangMenu from "@/components/LangMenu";
import UserMenu from "@/components/UserMenu";
import DashboardNotificationControl from "@/components/DashboardNotificationControl";
import Breadcrumbs from "@/components/dashboard/Breadcrumbs";

export default function DashboardTopNavigation() {
  return (
    <div className="flex h-[72px] justify-center border-b border-gray-200 bg-white">
      <header className="m-auto flex w-full flex-row items-center justify-between px-6">
        <Breadcrumbs />
        <div className="flex items-center justify-center">
          <ul className="item-center flex list-none space-x-6">
            <li>
              <DashboardNotificationControl />
            </li>
            <li>
              <UserMenu />
            </li>
            <li>
              <LangMenu />
            </li>
          </ul>
        </div>
      </header>
    </div>
  );
}
