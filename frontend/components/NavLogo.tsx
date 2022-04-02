import { ShoppingCartIcon } from "@heroicons/react/solid";

export default function Navigation() {
  return (
    <>
      <ShoppingCartIcon className="inline h-5 w-5 text-indigo-500" />{" "}
      <span className="font-bold underline decoration-indigo-500 decoration-2 dark:text-slate-200">
        Laxo
      </span>
    </>
  );
}
