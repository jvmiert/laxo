import { useRouter } from "next/router";
import Link from "next/link";
import { Transition, Popover } from "@headlessui/react";
import { ChevronDownIcon } from "@heroicons/react/solid";
import { GlobeAltIcon } from "@heroicons/react/outline";

export default function LangMenu() {
  const { pathname, locale } = useRouter();
  return (
    <Popover className="relative inline-block text-left">
      <div>
        <Popover.Button>
          <GlobeAltIcon className="inline h-5 w-5 text-gray-900" />
          <ChevronDownIcon className="inline h-5 w-5 text-gray-900" />
        </Popover.Button>
      </div>
      <Transition
        enter="transition ease-out duration-100"
        enterFrom="transform opacity-0 scale-95"
        enterTo="transform opacity-100 scale-100"
        leave="transition ease-in duration-75"
        leaveFrom="transform opacity-100 scale-100"
        leaveTo="transform opacity-0 scale-95"
      >
        <Popover.Panel className="absolute right-0 mt-2 w-56 origin-top-right rounded-md bg-white shadow-lg ring-1 ring-black ring-opacity-5 focus:outline-none">
          <div className="py-1">
            <Popover.Button as="div">
              <Link href={pathname} locale="en">
                <a
                  className={`block px-4 py-2 text-sm ${
                    locale == "en"
                      ? "bg-gray-100 text-gray-900"
                      : "text-gray-700"
                  }`}
                >
                  English
                </a>
              </Link>
            </Popover.Button>
            <Popover.Button as="div">
              <Link href={pathname} locale="vi">
                <a
                  className={`block px-4 py-2 text-sm ${
                    locale == "vi"
                      ? "bg-gray-100 text-gray-900"
                      : "text-gray-700"
                  }`}
                >
                  Tiếng Việt
                </a>
              </Link>
            </Popover.Button>
          </div>
        </Popover.Panel>
      </Transition>
    </Popover>
  );
}
