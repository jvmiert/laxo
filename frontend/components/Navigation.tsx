import { useRouter } from "next/router";
import Link from "next/link";
import { useAuth } from "@/providers/AuthProvider";
import LangMenu from "@/components/LangMenu";
import NavLogo from "@/components/NavLogo";

interface NavLinkProps {
  currentPath: string;
  href: string;
  navText: string;
}

function NavLink(props: NavLinkProps) {
  const { currentPath, href, navText } = props;

  return (
    <li>
      <Link href={href}>
        <a className={`${currentPath == href ? "underline" : ""}`}>{navText}</a>
      </Link>
    </li>
  );
}

export default function Navigation() {
  const { pathname } = useRouter();
  const { auth } = useAuth();

  return (
    <div className="flex h-16 w-full justify-center border-b border-gray-200">
      <header className="m-auto flex w-4/5 max-w-5xl flex-row items-center justify-between px-6">
        <div className="flex items-center justify-center">
          <Link href="/" passHref>
            <span className="cursor-pointer">
              <NavLogo />
            </span>
          </Link>
        </div>
        <div className="flex items-center justify-center">
          <ul className="item-center flex list-none space-x-6">
            {!auth && (
              <>
                <NavLink
                  currentPath={pathname}
                  href="/login"
                  navText="Log In"
                />
                <li>
                  <Link href="/register">
                    <a className="w-full rounded-md bg-indigo-500 py-2 px-4 font-bold text-white shadow-md shadow-indigo-500/50 hover:bg-indigo-700 focus:outline-none focus:ring focus:ring-indigo-200">
                      Sign Up
                    </a>
                  </Link>
                </li>
              </>
            )}
            {auth && (
              <>
                <NavLink
                  currentPath={pathname}
                  href="/dashboard/home"
                  navText="Dashboard"
                />
                <NavLink
                  currentPath={pathname}
                  href="/logout"
                  navText="Logout"
                />
              </>
            )}
            <li>
              <LangMenu />
            </li>
          </ul>
        </div>
      </header>
    </div>
  );
}
