import { useRouter } from "next/router";
import type { InputHTMLAttributes, ChangeEvent } from "react";
import { SearchIcon, XCircleIcon } from "@heroicons/react/solid";
import Link from "next/link";

type SearchInputProps = InputHTMLAttributes<HTMLInputElement> & {};

export default function SearchInput({ ...props }: SearchInputProps) {
  const {
    push,
    pathname,
    query: { p, l: queryLimitNumber, s: searchQuery },
  } = useRouter();

  const currentLimit = Number(queryLimitNumber) ? Number(queryLimitNumber) : 10;
  const currentSearchQuery = searchQuery ? searchQuery.toString() : "";

  const handleSearch = (e: ChangeEvent<HTMLInputElement>) => {
    push(
      {
        pathname: pathname,
        query: {
          ...(currentLimit > 10 && { l: queryLimitNumber }),
          ...(e.target.value != "" && { s: e.target.value }),
        },
      },
      undefined,
      { shallow: true, scroll: false },
    );
  };

  return (
    <div className="relative rounded-md border shadow">
      <div className="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-3">
        <span className="text-gray-500">
          <SearchIcon className="h-4 w-4" />
        </span>
      </div>
      <input
        {...props}
        value={currentSearchQuery}
        onChange={handleSearch}
        type="text"
        className="block w-full rounded-md py-2 pl-9 pr-11 focus:outline-none focus:ring focus:ring-indigo-200"
      />

      {currentSearchQuery != "" && (
        <div className="absolute inset-y-0 right-0 flex items-center pr-3">
          <Link
            shallow={true}
            scroll={false}
            href={{
              pathname: pathname,
              query: {
                ...(p && { p: p }),
                ...(queryLimitNumber && { l: queryLimitNumber }),
              },
            }}
          >
            <a className="z-10 rounded p-1 focus:outline-none focus:ring focus:ring-indigo-200">
              <XCircleIcon className="h-4 w-4" />
            </a>
          </Link>
        </div>
      )}
    </div>
  );
}
