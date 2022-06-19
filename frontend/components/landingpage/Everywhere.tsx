export default function Everywhere() {
  return (
    <div className="relative z-20">
      <div className="absolute flex h-[512px] w-[712px] flex-col overflow-hidden rounded-md bg-white shadow-xl shadow-indigo-400">
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
