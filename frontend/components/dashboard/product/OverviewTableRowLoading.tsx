export default function OverviewTableRowLoading() {
  return (
    <tr className="border-none">
      <td className="h-16 px-6 py-4">
        <div className="relative flex h-full w-full items-center">
          <div className="h-full w-8 animate-pulse rounded bg-slate-100" />
          <div className="ml-4 h-full grow animate-pulse rounded bg-slate-100" />
        </div>
      </td>
      <td className="h-16 px-6 py-4">
        <div className="h-full w-full animate-pulse rounded-lg bg-slate-100" />
      </td>
      <td className="h-16 px-6 py-4">
        <div className="h-full w-full animate-pulse rounded-lg bg-slate-100" />
      </td>
      <td className="h-16 px-6 py-4">
        <div className="h-full w-full animate-pulse rounded-lg bg-slate-100" />
      </td>
      <td className="h-16 px-6 py-4">
        <div className="h-full w-full animate-pulse rounded-lg bg-slate-100" />
      </td>
    </tr>
  );
}
