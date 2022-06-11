import cc from "classcat";
import { useEffect, useRef, useState, ChangeEvent } from "react";
import { useDashboard } from "@/providers/DashboardProvider";
import Image from "next/image";
import { CloudUploadIcon } from "@heroicons/react/outline";
import { useAxios } from "@/providers/AxiosProvider";

const shimmer = `
<svg width="48px" height="48px" version="1.1" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink">
  <defs>
    <linearGradient id="g">
      <stop stop-color="#E2E8F0" offset="20%" />
      <stop stop-color="#F1F5F9" offset="50%" />
      <stop stop-color="#E2E8F0" offset="70%" />
    </linearGradient>
  </defs>
  <rect width="48px" height="48px" fill="#E2E8F0" />
  <rect id="r" width="48px" height="48px" fill="url(#g)" />
  <animate xlink:href="#r" attributeName="x" from="-48px" to="48px" dur="1s" repeatCount="indefinite"  />
</svg>`;

const shimmerBase64 = () =>
  typeof window === "undefined"
    ? Buffer.from(shimmer).toString("base64")
    : window.btoa(shimmer);

type AssetManagementProps = {
  mediaList: string[];
};

export default function AssetManagement({ mediaList }: AssetManagementProps) {
  const [dragActive, setDragActive] = useState(false);
  const { activeShop } = useDashboard();

  const dropRef = useRef<HTMLDivElement>(null);
  const inputRef = useRef<HTMLInputElement>(null);

  const { axiosClient } = useAxios();

  const preventDefaultFunc = (e: DragEvent) => {
    e.preventDefault();
  };

  const openFileDialog = () => {
    if (inputRef.current) {
      inputRef.current.click();
    }
  };

  const onFileChange = (e: ChangeEvent<HTMLInputElement>) => {
    e.preventDefault();
    e.persist();

    if (!e.target?.files) return;

    Array.from(e.target.files).forEach(async (f) => {
      console.log(f);
      axiosClient.post("/manage-assets", f, {
        headers: {
          "Content-Type": f.type,
        },
      });
    });
  };

  const dragEnterFunc = (e: DragEvent) => {
    const target = e.target as HTMLDivElement;
    e.preventDefault();
    if (target.id == "dropDiv") {
      setDragActive(true);
    }
  };

  const dragLeaveFunc = (e: DragEvent) => {
    const target = e.target as HTMLDivElement;
    e.preventDefault();
    if (target.id == "dropDiv") {
      setDragActive(false);
    }
  };

  const dropFunc = (e: DragEvent) => {
    e.preventDefault();
    if (e.dataTransfer?.items) {
      for (let i = 0; i < e.dataTransfer.items.length; i++) {
        if (e.dataTransfer.items[i].kind === "file") {
          const type = e.dataTransfer.items[i].type;
          const file = e.dataTransfer.items[i].getAsFile();
          console.log(file, type);
        }
      }
    } else if (e.dataTransfer?.files) {
      for (let i = 0; i < e.dataTransfer.files.length; i++) {
        const type = e.dataTransfer.types[i];
        const file = e.dataTransfer.files[i];
        console.log(file, type);
      }
    }
    setDragActive(false);
  };

  useEffect(() => {
    window.addEventListener("dragover", preventDefaultFunc, false);
    window.addEventListener("drop", preventDefaultFunc, false);

    return () => {
      window.removeEventListener("dragover", preventDefaultFunc);
      window.removeEventListener("drop", preventDefaultFunc);
    };
  }, []);

  useEffect(() => {
    let dropRefStored: HTMLDivElement;

    if (dropRef.current) {
      dropRefStored = dropRef.current;
      dropRefStored.addEventListener("dragenter", dragEnterFunc, false);
      dropRefStored.addEventListener("dragleave", dragLeaveFunc, false);
      dropRefStored.addEventListener("dragover", preventDefaultFunc, false);
      dropRefStored.addEventListener("drop", dropFunc, false);
    }

    return () => {
      if (dropRefStored) {
        dropRefStored.removeEventListener("dragenter", dragEnterFunc);
        dropRefStored.removeEventListener("dragleave", dragLeaveFunc);
        dropRefStored.removeEventListener("dragover", preventDefaultFunc);
        dropRefStored.removeEventListener("drop", dropFunc);
      }
    };
  }, [dropRef]);

  if (!activeShop) return <></>;

  return (
    <div>
      <div className="mx-auto max-w-sm">
        <div
          ref={dropRef}
          className={cc([
            "border-grey relative flex cursor-pointer flex-col items-center rounded border px-4 py-8 text-center hover:border-transparent hover:ring-2 hover:ring-indigo-200",
            { "border-dashed border-slate-400": !dragActive },
            { "border-transparent ring-2 ring-indigo-200": dragActive },
          ])}
        >
          <div
            onClick={openFileDialog}
            id="dropDiv"
            className="absolute inset-0 z-10"
          />
          <CloudUploadIcon className="mb-5 h-8 w-8 text-gray-700" />
          <span>Drop your images or click here to add them</span>
          <input
            multiple
            accept=".png,.gif,.jpeg,.jpg"
            className="hidden"
            ref={inputRef}
            type="file"
            onChange={onFileChange}
          />
        </div>
      </div>
      <div className="mt-4">
        {mediaList.map((m) => (
          <div key={m} className="relative h-12 w-12">
            <Image
              className="rounded"
              alt={"Product preview"}
              src={`/api/assets/${activeShop.assetsToken}/products/${m}`}
              layout="fill"
              placeholder="blur"
              blurDataURL={`data:image/svg+xml;base64,${shimmerBase64()}`}
              objectFit="contain"
              objectPosition="center"
            />
          </div>
        ))}
      </div>
    </div>
  );
}
