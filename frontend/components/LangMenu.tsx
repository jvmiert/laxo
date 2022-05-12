import { Fragment, forwardRef, ReactNode } from "react";
import cc from "classcat";
import { useRouter } from "next/router";
import Link, { LinkProps } from "next/link";
import { Menu, Transition } from "@headlessui/react";
import { ChevronDownIcon } from "@heroicons/react/solid";
import { GlobeAltIcon } from "@heroicons/react/outline";

type EnhancedLinkProps = {
  children: ReactNode[] | ReactNode;
  className: string;
} & LinkProps;

const LangLink = forwardRef<HTMLAnchorElement, EnhancedLinkProps>(
  (props, ref) => {
    let { href, children, locale, ...rest } = props;
    return (
      <Link href={href} locale={locale}>
        <a ref={ref} {...rest}>
          {children}
        </a>
      </Link>
    );
  },
);

LangLink.displayName = "LangLink";

export default function LangMenu() {
  const { pathname, locale } = useRouter();
  return (
    <Menu as="div" className="relative z-10 inline-block text-left">
      <Menu.Button>
        <GlobeAltIcon className="inline h-5 w-5 text-gray-900" />
        <ChevronDownIcon className="inline h-5 w-5 text-gray-900" />
      </Menu.Button>
      <Transition
        as={Fragment}
        enter="transition ease-out duration-100"
        enterFrom="transform opacity-0 scale-95"
        enterTo="transform opacity-100 scale-100"
        leave="transition ease-in duration-75"
        leaveFrom="transform opacity-100 scale-100"
        leaveTo="transform opacity-0 scale-95"
      >
        <Menu.Items className="absolute right-0 mt-2 w-56 origin-top-right rounded-md bg-white shadow-lg ring-1 ring-black ring-opacity-5 focus:outline-none">
          <div className="py-1">
            <Menu.Item>
              <LangLink
                href={pathname}
                locale="en"
                className={cc([
                  "block px-4 py-2 text-sm",
                  { "bg-gray-100 text-gray-900": locale == "en" },
                  { "text-gray-700": locale != "en" },
                ])}
              >
                English
              </LangLink>
            </Menu.Item>
            <Menu.Item>
              <LangLink
                href={pathname}
                locale="vi"
                className={cc([
                  "block px-4 py-2 text-sm",
                  { "bg-gray-100 text-gray-900": locale == "vi" },
                  { "text-gray-700": locale != "vi" },
                ])}
              >
                Tiếng Việt
              </LangLink>
            </Menu.Item>
          </div>
        </Menu.Items>
      </Transition>
    </Menu>
  );
}
