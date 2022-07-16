import cc from "classcat";
import {
  RenderElementProps,
  useSlateStatic,
  useSelected,
  ReactEditor,
} from "slate-react";
import { Transforms } from "slate";
import Image from "next/image";
import { TrashIcon } from "@heroicons/react/solid";

import { LaxoImageElement } from "@/lib/laxoSlate";

export default function ImageElement({
  attributes,
  children,
  element,
  token,
}: RenderElementProps & { token: string; element: LaxoImageElement }) {
  const editor = useSlateStatic();
  const path = ReactEditor.findPath(editor, element);

  const selected = useSelected();

  //@TODO: implement
  // const focused = useFocused();

  return (
    <div {...attributes}>
      {children}
      <div
        draggable
        contentEditable={false}
        className="relative cursor-pointer"
      >
        <div
          className={cc([
            "absolute inset-0 z-10 opacity-40",
            { "bg-black": selected },
            { "bg-transparent": !selected },
          ])}
        />
        <div
          className={cc([
            "absolute inset-0 z-20 flex justify-center",
            { invisible: !selected },
          ])}
        >
          <button
            onClick={() => Transforms.removeNodes(editor, { at: path })}
            className=""
          >
            <TrashIcon className="z-20 h-8 w-8 fill-white" />
          </button>
        </div>
        <Image
          alt={"Product description!"}
          src={`/api/assets/${token}/${element.src}`}
          width={element.width}
          height={element.height}
          layout="responsive"
        />
      </div>
    </div>
  );
}
