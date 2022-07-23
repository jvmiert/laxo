import type { FC, ComponentProps } from "react";
import Link from "next/link";

type PrimaryButtonProps = {
  href: string;
  buttonText: string;
  Icon: FC<ComponentProps<"svg">>;
};

export default function PrimaryButton({
  buttonText,
  href,
  Icon,
}: PrimaryButtonProps) {
  return (
    <Link href={href}>
      <a className="inline-flex items-center rounded-md border border-indigo-500 bg-indigo-500 py-2 px-4 text-white shadow shadow-indigo-500/50 hover:bg-indigo-700 focus:outline-none focus:ring focus:ring-indigo-200 disabled:cursor-not-allowed disabled:bg-indigo-200">
        <Icon className="mr-2 -ml-1 h-4 w-4" />
        {buttonText}
      </a>
    </Link>
  );
}
