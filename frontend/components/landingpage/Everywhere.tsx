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

function Lazada() {
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
      <div className="px-3 py-2 font-semibold">
        This is a very cool product name
      </div>
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
  return (
    <div className="relative z-20">
      <div
        style={{
          boxShadow:
            "0 50px 100px -20px rgba(50,50,93,.25),0 30px 60px -30px rgba(0,0,0,.3),inset 0 -2px 6px 0 rgba(10,37,64,.35)",
        }}
        className="absolute z-30 flex h-[480px] w-[250px] translate-y-1/2 -translate-x-8 cursor-default select-none overflow-hidden rounded-3xl bg-[#f6f9fc] p-2 shadow-xl"
      >
        <div className="h-full w-full overflow-hidden rounded-3xl bg-white py-3">
          <Lazada />
        </div>
      </div>
      <div className="absolute flex h-[512px] w-[712px] -translate-y-1/2 flex-col overflow-hidden rounded-md bg-white shadow-xl">
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
        <div className="grow bg-slate-100"></div>
      </div>
    </div>
  );
}
