import { useRouter } from "next/router";
import Link from "next/link";

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

  return (
    <div>
      <ul className="list-none">
        <NavLink currentPath={pathname} href="/" navText="Home" />
        <NavLink currentPath={pathname} href="/login" navText="Login" />
        <NavLink currentPath={pathname} href="/register" navText="Register" />
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
