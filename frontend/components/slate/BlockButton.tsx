import cc from "classcat";
import { useSlate } from "slate-react";

import { isBlockActive, TEXT_ALIGN_TYPES, toggleBlock } from "@/lib/laxoSlate";

export default function BlockButton({
  format,
  className,
  icon,
}: {
  format: any;
  className?: string;
  icon: JSX.Element | string;
}) {
  const editor = useSlate();
  const active = isBlockActive(
    editor,
    format,
    TEXT_ALIGN_TYPES.includes(format) ? "align" : "type",
  );
  return (
    <button
      type="button"
      className={cc([
        className,
        { "bg-white text-gray-700 hover:bg-gray-50": !active },
        { "bg-indigo-50 text-indigo-500": active },
      ])}
      onMouseDown={(event) => {
        event.preventDefault();
        toggleBlock(editor, format);
      }}
    >
      {icon}
    </button>
  );
}
