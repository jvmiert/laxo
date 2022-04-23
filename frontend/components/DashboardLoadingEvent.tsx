import cc from "classcat";
import { useEffect, useState } from "react";
import { CollectionIcon, CheckIcon } from "@heroicons/react/outline";

export default function DashboardLoadingEvent() {
  const [progress, setProgress] = useState(0);
  const [pProduct, setPProduct] = useState(0);
  const [pText, setPText] = useState("Fetching");
  const [showDone, setShowDone] = useState(false);

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
    const timeoutM = setTimeout(() => setLevel(60, "Saving products", 8), 4000);
    const timeoutE = setTimeout(
      () => setLevel(100, "Processing products", 10),
      6000,
    );

    return function cleanup() {
      clearTimeout(timeout);
      clearTimeout(timeoutM);
      clearTimeout(timeoutE);
    };
  }, []);

  useEffect(() => {
    let timeOut: ReturnType<typeof setTimeout>;

    if (progress === 100) {
      timeOut = setTimeout(() => setShowDone(true), 1800);
    }

    return function cleanup() {
      if (timeOut) {
        clearTimeout(timeOut);
      }
    };
  }, [progress]);
  return (
    <div className="flex flex-row rounded px-4">
      <div className="pr-4">
        <div className="rounded-full bg-indigo-100 p-4">
          {showDone ? (
            <CheckIcon className="h-4 w-4 text-indigo-500" />
          ) : (
            <CollectionIcon className="h-4 w-4 text-indigo-500" />
          )}
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
              className={cc([
                "flex",
                { "animate-pulse": !showDone },
                { "bg-gradient-to-l from-indigo-400 to-indigo-600": !showDone },
                "bg-indigo-600",
              ])}
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
