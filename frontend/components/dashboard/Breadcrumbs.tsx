import { useRouter } from "next/router";

import { ChevronRightIcon } from "@heroicons/react/solid";

export default function Breadcrumbs() {
  const { pathname } = useRouter();

  const splitPath = pathname.split("/");

  const depth = splitPath.length - 2;

  return (
    <nav className="flex" aria-label="Breadcrumb">
      <ol role="list" className="flex items-center space-x-4">
        <li>
          <div className="flex items-center">
            <a
              href={"/test"}
              className="mr-4 text-sm font-medium capitalize text-gray-500 hover:text-gray-700"
            >
              {splitPath[2]}
            </a>
            {depth > 1 && (
              <>
                <ChevronRightIcon
                  className="h-5 w-5 flex-shrink-0 text-gray-400"
                  aria-hidden="true"
                />

                <a
                  href={"/test"}
                  className="ml-4 text-sm font-medium text-gray-500 hover:text-gray-700"
                >
                  {splitPath[3]}
                </a>
              </>
            )}
          </div>
        </li>
      </ol>
    </nav>
  );
}
