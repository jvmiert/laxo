import type { FC, ComponentProps } from "react";
import Link from "next/link";

type SecondaryButtonProps = {
  href: string;
  buttonText: string;
  Icon: FC<ComponentProps<"svg">>;
};

export default function SecondaryButton({
  buttonText,
  href,
  Icon,
}: SecondaryButtonProps) {
  return (
    <Link href={href}>
      <a className="inline-flex items-center rounded-md border border-slate-300 bg-white py-2 px-4 text-slate-700 shadow shadow-slate-300/50 hover:bg-slate-50 focus:outline-none focus:ring focus:ring-indigo-200 disabled:cursor-not-allowed disabled:bg-slate-200">
        <Icon className="mr-2 -ml-1 h-4 w-4" />
        {buttonText}
      </a>
    </Link>
  );
}
