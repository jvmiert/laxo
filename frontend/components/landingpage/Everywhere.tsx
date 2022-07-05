import React, { useState } from "react";

import {
  ChevronLeftIcon,
  ChevronRightIcon,
  ChevronDownIcon,
  DotsHorizontalIcon,
  ChatAlt2Icon,
} from "@heroicons/react/outline";
import Image from "next/image";

import product from "@/assets/product1.jpg";

import lazSearch from "@/assets/lazSearch.png";
import lazShare from "@/assets/lazShare.png";
import lazCart from "@/assets/lazCart.png";
import lazShop from "@/assets/lazShop.svg";
import lazLike from "@/assets/lazLike.png";
import lazStar from "@/assets/lazStar.png";
import gradientArrow from "@/assets/gradientArrow.svg";

import tiktokBack from "@/assets/tiktokBack.png";
import tiktokCart from "@/assets/tiktokCart.png";
import tiktokArrow from "@/assets/tiktokArrow.png";
import tiktokDots from "@/assets/tiktokDots.png";
import tiktokShop from "@/assets/tiktokShop.png";
import tiktokChat from "@/assets/tiktokChat.png";
import tiktokShip from "@/assets/tiktokShip.webp";

import fbBack from "@/assets/fbBack.png";
import fbCart from "@/assets/fbCart.png";
import fbDots from "@/assets/fbDots.png";
import fbShare from "@/assets/fbShare.png";
import fbMark from "@/assets/fbMark.png";

function Facebook({ name }) {
  return (
    <>
      <div className="flex w-full items-center justify-between border-b border-slate-100 py-2 px-3">
        <div className="flex items-center space-x-3">
          <div className="h-4 w-4">
            <Image src={fbBack} alt="back" layout="responsive" />
          </div>
          <span className="text-xs font-semibold">Laxo Shop</span>
        </div>
        <div className="flex space-x-2">
          <div className="h-4 w-4">
            <Image src={fbCart} alt="cart" layout="responsive" />
          </div>
          <div className="h-4 w-4">
            <Image src={fbDots} alt="dots" layout="responsive" />
          </div>
        </div>
      </div>
      <div>
        <Image src={product} alt="product" layout="responsive" />
      </div>
      <div className="flex w-full items-start justify-between px-3 pt-2">
        <div className="text-xs font-bold">{name}</div>
        <div className="flex items-center space-x-3">
          <div className="h-4 w-4">
            <Image
              src={fbMark}
              className="opacity-70"
              alt="bookmark"
              layout="responsive"
            />
          </div>
          <div className="h-4 w-4">
            <Image
              src={fbShare}
              className="opacity-70"
              alt="share"
              layout="responsive"
            />
          </div>
        </div>
      </div>
      <div className="px-3 pb-3 text-[0.6rem]">₫137.000</div>
      <div className="mx-3 rounded bg-[#1a74e5] py-1 text-center text-white">
        Message
      </div>
      <div className="w-full pt-1 text-center text-[0.5rem] text-slate-500">
        Message this seller to ask about the product
      </div>
      <div className="h-2 bg-[#f0f1f5]" />
      <div className="flex items-center justify-between">
        <div className="py-2 pl-3 text-xs font-bold">Description</div>
        <ChevronDownIcon className="mr-3 h-4 w-4 stroke-slate-500 stroke-[3px]" />
      </div>
      <div className="h-2 bg-[#f0f1f5]" />
      <div className="flex w-full items-center justify-between py-2 px-3">
        <div className="flex flex-col">
          <div className="text-xs font-semibold">Laxo Shop</div>
          <div className="text-[0.65rem] text-slate-600">30K followers</div>
        </div>
        <div>
          <ChevronRightIcon className="h-5 w-5" />
        </div>
      </div>
      <div className="mx-3 rounded bg-[#e8f2fe] py-1 text-center text-[#1065d0]">
        Follow
      </div>
    </>
  );
}

function Tiktok({ name }) {
  return (
    <>
      <div className="relative w-full">
        <div className="absolute inset-x-0 top-0 z-10 flex items-center justify-between bg-gradient-to-b from-black/40 to-black/0 px-3 pt-4 pb-4">
          <div className="h-4 w-4">
            <Image
              src={tiktokBack}
              className="drop-shadow-md"
              alt="back"
              layout="responsive"
            />
          </div>
          <div className="flex space-x-2">
            <div className="h-4 w-4">
              <Image
                src={tiktokArrow}
                className="drop-shadow-md"
                alt="back"
                layout="responsive"
              />
            </div>
            <div className="h-4 w-4">
              <Image
                src={tiktokCart}
                className="drop-shadow-md"
                alt="back"
                layout="responsive"
              />
            </div>
            <div className="h-4 w-4">
              <Image
                src={tiktokDots}
                className="drop-shadow-md"
                alt="back"
                layout="responsive"
              />
            </div>
          </div>
        </div>
        <Image src={product} alt="product" layout="responsive" />
      </div>
      <div className="px-3 pt-2 font-semibold text-[#ff2758]">137.000₫</div>
      <div className="px-3 pt-1 text-xs font-medium text-slate-800">{name}</div>
      <div className="flex items-center justify-between px-3 pb-2 pt-1.5 text-[10px] font-medium">
        <div className="text-slate-400">137 sold</div>
        <div className="flex items-center text-[#fe2c55]">
          <div className="mr-1 flex h-3 w-3">
            <Image src={tiktokShip} alt="ship" />
          </div>
          <div>Free shipping</div>
        </div>
      </div>
      <div className="h-2 bg-slate-50" />
      <div className="flex items-center justify-between">
        <div className="py-2 pl-3 text-xs font-medium">Coupons</div>
        <ChevronRightIcon className="mr-3 h-4 w-4 stroke-slate-400" />
      </div>
      <div className="h-2 bg-slate-50" />
      <div className="flex items-center justify-between">
        <div className="py-2 pl-3 text-xs font-medium">Select options</div>
        <ChevronRightIcon className="mr-3 h-4 w-4 stroke-slate-400" />
      </div>
      <div className="mx-4 border-b border-slate-100 pt-0.5" />
      <div className="flex items-center justify-between pt-0.5">
        <div className="py-2 pl-3 text-xs font-medium">
          Shipping{" "}
          <span className="text-[9px] text-slate-400">(5 - 7 days)</span>
        </div>
        <div className="flex items-center">
          <div className="text-[9px] text-[#fe2c55]">Free</div>
          <ChevronRightIcon className="mr-3 h-4 w-4 stroke-slate-400" />
        </div>
      </div>
      <div className="flex items-center px-2 py-2">
        <div className="relative mr-2 flex flex-col items-center">
          <div className="h-4 w-4">
            <Image src={tiktokShop} alt="shop" />
          </div>
          <div className="text-[10px]">Store</div>
        </div>
        <div className="flex flex-col items-center">
          <div className="h-4 w-4">
            <Image src={tiktokChat} alt="chat" />
          </div>
          <div className="text-[10px]">Chat</div>
        </div>
        <div className="ml-2 flex grow space-x-2 text-center text-xs">
          <div className="flex w-1/2 items-center justify-center rounded border-2 border-[#fe2c55] bg-white py-1 px-2 leading-3 text-[#fe2c55]">
            Add to cart
          </div>
          <div className="flex w-1/2 items-center justify-center rounded bg-[#fe2c55] py-1 px-2 leading-3 text-white">
            Buy with coupon
          </div>
        </div>
      </div>
    </>
  );
}

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
      <div>
        <Image src={product} alt="product" layout="responsive" />
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
      <div className="h-2 bg-slate-100" />
      <div className="flex items-center justify-between">
        <div className="py-2 pl-3 text-xs font-bold">Vouchers</div>
        <ChevronRightIcon className="mr-3 h-3 w-3 stroke-slate-500 stroke-[3px]" />
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
          className="w-full overflow-hidden rounded-3xl bg-[#f6f9fc] p-2"
        >
          <div className="h-full w-full overflow-hidden rounded-3xl bg-white pb-3">
            <Tiktok name={name} />
          </div>
        </div>
      </div>

      <div className="z-20 flex h-[480px] w-[250px] cursor-default select-none overflow-hidden rounded-3xl bg-[#f6f9fc] shadow-xl">
        <div
          style={{
            boxShadow:
              "rgba(0, 0, 0, 0.1) 0px 8px 10px -6px, inset 0 -2px 6px 0 rgba(10,37,64,.35)",
          }}
          className="w-full overflow-hidden rounded-3xl bg-[#f6f9fc] p-2"
        >
          <div className="h-full w-full overflow-hidden rounded-3xl bg-white py-3">
            <Facebook name={name} />
          </div>
        </div>
      </div>
    </div>
  );
}
