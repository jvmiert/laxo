import Link from "next/link";
import LangMenu from "@/components/LangMenu";
import UserMenu from "@/components/UserMenu";
import NavLogo from "@/components/NavLogo";
import DashboardNotificationControl from "@/components/DashboardNotificationControl";

export default function DashboardTopNavigation() {
  return (
    <div className="flex h-[72px] w-full justify-center border-b border-gray-200">
      <header className="m-auto flex w-full flex-row items-center justify-between px-6">
        <div className="flex items-center justify-center">
          <Link href="/" passHref>
            <span className="cursor-pointer">
              <NavLogo />
            </span>
          </Link>
        </div>
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
