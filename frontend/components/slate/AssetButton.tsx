import { PhotographIcon } from "@heroicons/react/outline";

export default function AssetButton({ openFunc }: { openFunc: () => void }) {
  return (
    <button
      className="relative -ml-px inline-flex items-center rounded-tr-md border-t border-l border-r px-4 py-2 text-sm font-medium focus:z-10 focus:border-indigo-500 focus:outline-none focus:ring-1 focus:ring-indigo-500"
      type="button"
      onMouseDown={(event) => {
        event.preventDefault();
        openFunc();
      }}
    >
      <PhotographIcon className="h-4 w-4 " />
    </button>
  );
}
