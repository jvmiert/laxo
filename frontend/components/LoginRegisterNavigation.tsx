import Link from "next/link";

import NavLogo from "@/components/NavLogo";
import LangMenu from "@/components/LangMenu";

export default function LoginRegisterNavigation() {
  return (
    <div className="mb-6 flex items-center justify-between">
      <Link href="/" passHref>
        <span className="cursor-pointer">
          <NavLogo />
        </span>
      </Link>
      <LangMenu />
    </div>
  );
}
