import React, { useState } from "react";

import {
  ChevronLeftIcon,
  DotsHorizontalIcon,
  ChatAlt2Icon,
} from "@heroicons/react/outline";
import Image from "next/image";

import lazSearch from "@/assets/lazSearch.png";
import lazShare from "@/assets/lazShare.png";
import lazCart from "@/assets/lazCart.png";
import lazShop from "@/assets/lazShop.svg";
import lazLike from "@/assets/lazLike.png";
import lazStar from "@/assets/lazStar.png";
import gradientArrow from "@/assets/gradientArrow.svg";

function Lazada({ name }) {
  return (
    <>
      <div className="flex items-center justify-between px-3 py-2">
        <ChevronLeftIcon className="h-4 w-4" />
        <div className="flex h-6 w-2/5 items-center rounded-full border">
          <div className="ml-1 h-4 w-4">
            <Image src={lazSearch} alt="search" layout="responsive" />
          </div>
        </div>
        <div className="h-4 w-4">
          <Image src={lazShare} alt="share" layout="responsive" />
        </div>
        <div className="h-4 w-4">
          <Image src={lazCart} alt="share" layout="responsive" />
        </div>
        <DotsHorizontalIcon className="h-4 w-4" />
      </div>
      <div className="flex items-center px-3 py-2">
        <div className="grow">
          <p className="font-bold text-[#ea5d63]">
            <span className="mr-1.5 text-xs font-medium">₫</span>
            133,700
          </p>
        </div>
        <div className="flex items-center space-x-2">
          <div className="h-4 w-4">
            <Image src={lazShare} alt="share" />
          </div>
          <div className="h-4 w-4">
            <Image src={lazLike} alt="like" />
          </div>
        </div>
      </div>
      <div className="px-3 py-2 font-semibold">{name}</div>
      <div className="flex items-center px-3 py-2">
        <div className="mr-5 flex items-baseline text-xs">
          <span>5</span>
          <span className="mr-1 text-[8px]">/5</span>
          <div className="mr-0.5 h-2 w-2">
            <Image src={lazStar} alt="star" />
          </div>
          <div className="mr-0.5 h-2 w-2">
            <Image src={lazStar} alt="star" />
          </div>
          <div className="mr-0.5 h-2 w-2">
            <Image src={lazStar} alt="star" />
          </div>
          <div className="mr-0.5 h-2 w-2">
            <Image src={lazStar} alt="star" />
          </div>
          <div className="h-2 w-2">
            <Image src={lazStar} alt="star" />
          </div>
        </div>
        <div className="text-xs">
          <span className="font-medium">89</span> sold
        </div>
      </div>
      <div className="flex items-center border-t border-slate-100 px-3 py-2">
        <div className="relative mr-3 flex flex-col items-center">
          <div className="h-5 w-5">
            <Image src={lazShop} alt="store" />
          </div>
          <div className="text-xs">Store</div>
          <div className="absolute top-1 bottom-1 -right-1.5 w-[1px] bg-slate-200" />
        </div>
        <div className="flex flex-col items-center">
          <ChatAlt2Icon className="h-5 w-5 stroke-[1.5px]" />
          <div className="text-xs">Chat</div>
        </div>
        <div className="ml-2 flex grow space-x-0.5 text-center text-xs text-white">
          <div className="flex w-1/2 items-center justify-center rounded-full bg-[#f3ae3c] py-1 px-2 leading-3">
            Buy Now
          </div>
          <div className="flex w-1/2 items-center justify-center rounded-full bg-[#b73220] py-1 px-2 leading-3">
            Add to Cart
          </div>
        </div>
      </div>
    </>
  );
}

export default function Everywhere() {
  const [name, setName] = useState("Cool name");

  return (
    <div className="relative flex space-x-8">
      <style>{`
        @keyframes connectFlowDiagramArrowsLinear{
          to{
              -webkit-mask-position:11px 0;
              mask-position:11px 0
          }
      `}</style>
      <div
        style={{
          maskImage: `url(${gradientArrow.src})`,
          animation: "connectFlowDiagramArrowsLinear 500ms linear infinite",
        }}
        className="absolute top-1/2 left-0 right-10 block h-2 bg-gradient-to-r from-violet-500 to-fuchsia-500"
      />
      <div className="z-20 flex h-[480px] w-[712px] flex-col overflow-hidden rounded-md bg-white shadow-xl">
        <div className="flex py-1">
          <div className="flex items-center space-x-1.5 px-4">
            <div className="h-3 w-3 rounded-full bg-slate-100" />
            <div className="h-3 w-3 rounded-full bg-slate-100" />
            <div className="h-3 w-3 rounded-full bg-slate-100" />
          </div>
          <div className="grow">
            <div className="my-2 mx-16 rounded bg-slate-100">
              <div className="py-1 text-center text-[12px] font-light text-slate-900">
                dashboard.laxo.vn
              </div>
            </div>
          </div>
        </div>
        <div className="grow bg-slate-100">
          <div className="mx-auto mt-8 w-8/12 rounded-md bg-white py-4 px-3 shadow-sm">
            <div className="flex w-full justify-between rounded-xl bg-gray-50 px-2 py-1.5">
              <h3 className="text-sm font-medium leading-6 text-gray-900">
                General
              </h3>
            </div>
            <div className="grid grid-cols-8 gap-4 p-2 text-[12px]">
              <div className="col-span-5">
                <label className="mb-1 block text-gray-700">Name</label>

                <input
                  type="text"
                  className="focus:shadow-outline w-full appearance-none rounded border py-2 px-3 leading-tight text-gray-700 shadow focus:outline-none focus:ring focus:ring-indigo-200"
                  onChange={(e) => setName(e.target.value)}
                  value={name}
                />
              </div>

              <div className="col-span-3">
                <label className="mb-1 block text-gray-700">SKU</label>

                <input
                  type="text"
                  className="focus:shadow-outline w-full appearance-none rounded border py-2 px-3 leading-tight text-gray-700 shadow focus:outline-none focus:ring focus:ring-indigo-200"
                  defaultValue="Submit"
                />
              </div>
            </div>
          </div>
        </div>
      </div>

      <div className="z-20 flex h-[480px] w-[250px] cursor-default select-none overflow-hidden rounded-3xl bg-[#f6f9fc] shadow-xl">
        <div
          style={{
            boxShadow:
              "rgba(0, 0, 0, 0.1) 0px 8px 10px -6px, inset 0 -2px 6px 0 rgba(10,37,64,.35)",
          }}
          className="overflow-hidden rounded-3xl bg-[#f6f9fc] p-2"
        >
          <div className="h-full w-full overflow-hidden rounded-3xl bg-white py-3">
            <Lazada name={name} />
          </div>
        </div>
      </div>

      <div className="z-20 flex h-[480px] w-[250px] cursor-default select-none overflow-hidden rounded-3xl bg-[#f6f9fc] shadow-xl">
        <div
          style={{
            boxShadow:
              "rgba(0, 0, 0, 0.1) 0px 8px 10px -6px, inset 0 -2px 6px 0 rgba(10,37,64,.35)",
          }}
          className="overflow-hidden rounded-3xl bg-[#f6f9fc] p-2"
        >
          <div className="h-full w-full overflow-hidden rounded-3xl bg-white py-3">
            <Lazada name={name} />
          </div>
        </div>
      </div>

      <div className="z-20 flex h-[480px] w-[250px] cursor-default select-none overflow-hidden rounded-3xl bg-[#f6f9fc] shadow-xl">
        <div
          style={{
            boxShadow:
              "rgba(0, 0, 0, 0.1) 0px 8px 10px -6px, inset 0 -2px 6px 0 rgba(10,37,64,.35)",
          }}
          className="overflow-hidden rounded-3xl bg-[#f6f9fc] p-2"
        >
          <div className="h-full w-full overflow-hidden rounded-3xl bg-white py-3">
            <Lazada name={name} />
          </div>
        </div>
      </div>
    </div>
  );
}
