import cc from "classcat";
import { useSlate } from "slate-react";
import { FormatParameter } from "@/lib/laxoSlate";

import { isMarkActive, toggleMark } from "@/lib/laxoSlate";

export default function MarkButton({
  format,
  className,
  text,
}: {
  format: FormatParameter;
  className?: string;
  text: string;
}) {
  const editor = useSlate();
  const active = isMarkActive(editor, format);

  return (
    <button
      type="button"
      onMouseDown={(event) => {
        event.preventDefault();
        toggleMark(editor, format);
      }}
      className={cc([
        className,
        { "bg-white text-gray-700 hover:bg-gray-50": !active },
        { "bg-indigo-50 text-indigo-500": active },
      ])}
    >
      {text}
    </button>
  );
}
