import { useRouter } from "next/router";
import Link from "next/link";
import { useAuth } from "@/providers/AuthProvider";
import { ShoppingCartIcon } from "@heroicons/react/solid";
import LangMenu from "@/components/LangMenu";

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
              <ShoppingCartIcon className="inline h-5 w-5 text-pink-500" />{" "}
              <span className="font-bold underline decoration-pink-500 decoration-2 dark:text-slate-200">
                Laxo
              </span>
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
                <NavLink
                  currentPath={pathname}
                  href="/register"
                  navText="Sign Up"
                />
              </>
            )}
            {auth && (
              <>
                <NavLink
                  currentPath={pathname}
                  href="/dashboard"
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
