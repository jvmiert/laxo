import React, { useState } from "react";

import { Transition } from "@headlessui/react";

function Tab({ children, isShowing }) {
  return (
    <Transition
      appear={true}
      show={isShowing}
      enter="transition duration-75"
      enterFrom="translate-x-full opacity-30"
      enterTo="translate-x-0 opacity-100"
      leave="transition duration-150"
      leaveFrom="translate-x-0"
      leaveTo="-translate-x-full"
    >
      <div className="absolute">{children}</div>
    </Transition>
  );
}

export default function Everywhere() {
  const [tab, setTab] = useState(1);
  return (
    <div className="relative z-20">
      <div className="flex space-x-1 rounded-xl bg-blue-900/20 p-1">
        <button
          onClick={() => setTab(1)}
          className="w-full rounded-lg py-2.5 text-sm font-medium leading-5 text-blue-700 ring-white ring-opacity-60 ring-offset-2 ring-offset-blue-400 focus:outline-none focus:ring-2"
        >
          Test 1
        </button>
        <button
          onClick={() => setTab(2)}
          className="w-full rounded-lg py-2.5 text-sm font-medium leading-5 text-blue-700 ring-white ring-opacity-60 ring-offset-2 ring-offset-blue-400 focus:outline-none focus:ring-2"
        >
          Test 2
        </button>
        <button className="w-full rounded-lg py-2.5 text-sm font-medium leading-5 text-blue-700 ring-white ring-opacity-60 ring-offset-2 ring-offset-blue-400 focus:outline-none focus:ring-2">
          Test 3
        </button>
      </div>
      <div
        style={{
          boxShadow:
            "0 50px 100px -20px rgba(50,50,93,.25),0 30px 60px -30px rgba(0,0,0,.3),inset 0 -2px 6px 0 rgba(10,37,64,.35)",
        }}
        className="absolute z-30 flex h-[500px] w-[250px] overflow-hidden rounded-3xl bg-white p-2 shadow-xl"
      >
        <div className="h-full w-full overflow-hidden rounded-3xl bg-slate-100 p-3">
          <Tab isShowing={tab === 1}>Test test 1</Tab>
          <Tab isShowing={tab === 2}>Test test 2</Tab>
        </div>
      </div>
      <div className="absolute flex h-[512px] w-[712px] flex-col overflow-hidden rounded-md bg-white shadow-xl">
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
