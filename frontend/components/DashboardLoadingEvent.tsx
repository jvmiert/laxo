import { useEffect, useState } from "react";
import { CollectionIcon } from "@heroicons/react/outline";

export default function DashboardLoadingEvent() {
  const [progress, setProgress] = useState(0);
  const [pProduct, setPProduct] = useState(0);
  const [pText, setPText] = useState("Fetching");

  const setLevel = (n: number, s: string, i: number) => {
    setProgress(n);
    setPText(s);
    setPProduct(i);
  };
  useEffect(() => {
    const timeout = setTimeout(
      () => setLevel(20, "Fetching data from Lazada", 4),
      2000,
    );
    const timeoutM = setTimeout(() => setLevel(60, "Saving products", 8), 8000);
    const timeoutE = setTimeout(
      () => setLevel(100, "Processing products", 10),
      14000,
    );

    return function cleanup() {
      clearTimeout(timeout);
      clearTimeout(timeoutM);
      clearTimeout(timeoutE);
    };
  }, []);
  return (
    <div className="mb-4 flex flex-row rounded bg-gray-50 p-4">
      <div className="pr-4">
        <div className="rounded-full bg-indigo-100 p-4">
          <CollectionIcon className="h-6 w-6 text-indigo-500" />
        </div>
      </div>
      <div className="max-w-[400px] flex-grow space-y-2">
        <p className="text-sm font-medium">
          Adding your products from Lazada...
        </p>
        <div className="relative py-1">
          <div className="flex h-2 overflow-hidden rounded bg-indigo-200 shadow-md shadow-indigo-500/20">
            <div
              style={{
                width: `${progress}%`,
                transitionTimingFunction: "ease-out",
                transition: "width 2s",
              }}
              className="flex animate-pulse bg-gradient-to-l from-indigo-400 to-indigo-600"
            ></div>
          </div>
          <div className="flex place-content-between">
            <p className="mt-2 text-xs">{pText}...</p>
            <p className="mt-2 text-right text-xs">{pProduct} / 10</p>
          </div>
        </div>
      </div>
    </div>
  );
}
