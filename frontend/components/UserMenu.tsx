import { useIntl } from "react-intl";
import { Fragment, forwardRef, ReactNode } from "react";
import Link, { LinkProps } from "next/link";
import { Menu, Transition } from "@headlessui/react";
import { ChevronDownIcon } from "@heroicons/react/solid";
import { UserCircleIcon } from "@heroicons/react/outline";

type EnhancedLinkProps = {
  children: ReactNode[] | ReactNode;
  className: string;
} & LinkProps;

const UserLink = forwardRef<HTMLAnchorElement, EnhancedLinkProps>(
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

UserLink.displayName = "UserLink";

export default function UserMenu() {
  const t = useIntl();

  return (
    <Menu as="div" className="relative z-10 inline-block text-left">
      <Menu.Button>
        <UserCircleIcon className="inline h-5 w-5 text-gray-900" />
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
              <UserLink
                href="/logout"
                locale="en"
                className="block px-4 py-2 text-sm text-gray-700"
              >
                {t.formatMessage({
                  defaultMessage: "Logout",
                  description: "User dropdown menu: logout button",
                })}
              </UserLink>
            </Menu.Item>
          </div>
        </Menu.Items>
      </Transition>
    </Menu>
  );
}
