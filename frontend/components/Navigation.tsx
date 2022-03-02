import { useRouter } from "next/router";
import Link from "next/link";
import { useAuth } from "@/providers/AuthProvider";

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
  const { pathname, locale } = useRouter();
  const { auth } = useAuth();

  return (
    <div>
      <ul className="list-none">
        <NavLink currentPath={pathname} href="/" navText="Home" />
        {!auth && (
          <>
            <NavLink currentPath={pathname} href="/login" navText="Login" />
            <NavLink
              currentPath={pathname}
              href="/register"
              navText="Register"
            />
          </>
        )}
        {auth && (
          <NavLink currentPath={pathname} href="/logout" navText="Logout" />
        )}
      </ul>
      <ul className="list-none">
        <li>
          <Link href={pathname} locale="en">
            <a className={`${locale == "en" ? "underline" : ""}`}>English</a>
          </Link>
        </li>
        <li>
          <Link href={pathname} locale="vi">
            <a className={`${locale == "vi" ? "underline" : ""}`}>Vietnamese</a>
          </Link>
        </li>
      </ul>
    </div>
  );
}
