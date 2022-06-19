import cc from "classcat";
import React, { useCallback, useMemo, useState } from "react";
import Image from "next/image";
import {
  Editable,
  withReact,
  useSlate,
  Slate,
  RenderElementProps,
  RenderLeafProps,
  useSlateStatic,
  useSelected,
  useFocused,
  ReactEditor,
} from "slate-react";
import {
  Editor as SlateEditor,
  createEditor,
  Descendant,
  Transforms,
} from "slate";
import { withHistory } from "slate-history";
import { TrashIcon } from "@heroicons/react/solid";

import { FormatParameter, withImages } from "@/lib/laxoSlate";
import { useDashboard } from "@/providers/DashboardProvider";

const isMarkActive = (editor: SlateEditor, format: FormatParameter) => {
  const marks = SlateEditor.marks(editor);
  return marks ? marks[format] === true : false;
};

const toggleMark = (editor: SlateEditor, format: FormatParameter) => {
  const isActive = isMarkActive(editor, format);

  if (isActive) {
    SlateEditor.removeMark(editor, format);
  } else {
    SlateEditor.addMark(editor, format, true);
  }
};

const MarkButton = ({
  format,
  className,
  text,
}: {
  format: FormatParameter;
  className?: string;
  text: string;
}) => {
  const editor = useSlate();
  const active = isMarkActive(editor, format);

  return (
    <button
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
};

const ImageElement = ({
  attributes,
  children,
  element,
  token,
}: RenderElementProps & { token: string }) => {
  const editor = useSlateStatic();
  const path = ReactEditor.findPath(editor, element);

  const selected = useSelected();
  const focused = useFocused();

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
};

const Element = ({ attributes, children, element }: RenderElementProps) => {
  const { activeShop } = useDashboard();
  const token = activeShop ? activeShop.assetsToken : "";

  switch (element.type) {
    case "bulleted-list":
      return <ul {...attributes}>{children}</ul>;
    case "heading-one":
      return <h1 {...attributes}>{children}</h1>;
    case "heading-two":
      return <h2 {...attributes}>{children}</h2>;
    case "list-item":
      return <li {...attributes}>{children}</li>;
    case "numbered-list":
      return <ol {...attributes}>{children}</ol>;
    case "image":
      return (
        <ImageElement token={token} attributes={attributes} element={element}>
          {children}
        </ImageElement>
      );
    default:
      return <p {...attributes}>{children}</p>;
  }
};

const Leaf = ({ attributes, children, leaf }: RenderLeafProps) => {
  if (leaf.bold) {
    children = <strong>{children}</strong>;
  }

  if (leaf.code) {
    children = <code>{children}</code>;
  }

  if (leaf.italic) {
    children = <em>{children}</em>;
  }

  if (leaf.underline) {
    children = <u>{children}</u>;
  }

  return <span {...attributes}>{children}</span>;
};

const initialValue: Descendant[] = [
  {
    type: "paragraph",
    children: [{ text: "" }],
  },
];

export type EditorProps = {
  initialSchema: string;
};

export default function Editor({ initialSchema }: EditorProps) {
  const renderElement = useCallback(
    (props: RenderElementProps) => <Element {...props} />,
    [],
  );
  const renderLeaf = useCallback(
    (props: RenderLeafProps) => <Leaf {...props} />,
    [],
  );

  const [editor] = useState(() =>
    withImages(withHistory(withReact(createEditor()))),
  );
  //https://github.com/ianstormtaylor/slate/issues/3477
  //@TODO: figure out if we need above
  //const editor = useMemo(
  //  () => withImages(withHistory(withReact(createEditor()))),
  //  [],
  //);

  const slateValue = useMemo(() => {
    const parsedSchema = JSON.parse(initialSchema);
    // Slate throws an error if the value on the initial render is invalid
    // so we directly set the value on the editor in order
    // to be able to trigger normalization on the initial value before rendering
    editor.children = parsedSchema;
    try {
      SlateEditor.normalize(editor, { force: true });
    } catch {
      editor.children = initialValue;
      SlateEditor.normalize(editor, { force: true });
    }
    Transforms.deselect(editor);
    // We return the normalized internal value so that the rendering can take over from here
    return editor.children;
  }, [editor, initialSchema]);

  return (
    <Slate editor={editor} value={slateValue}>
      <div className="relative z-0 inline-flex rounded-md shadow-sm">
        <MarkButton
          text="B"
          format="bold"
          className="relative inline-flex items-center rounded-tl-md border-t border-l border-r border-gray-300 px-4 py-2 text-sm font-bold focus:z-10 focus:border-indigo-500 focus:outline-none focus:ring-1 focus:ring-indigo-500"
        />
        <MarkButton
          text="I"
          format="italic"
          className="relative -ml-px inline-flex items-center border-t border-l border-r border-gray-300 px-4 py-2 font-serif text-sm font-bold italic focus:z-10 focus:border-indigo-500 focus:outline-none focus:ring-1 focus:ring-indigo-500"
        />
        <MarkButton
          text="U"
          format="underline"
          className="relative -ml-px inline-flex items-center border-t border-l border-r border-gray-300 px-4 py-2 text-sm font-medium underline focus:z-10 focus:border-indigo-500 focus:outline-none focus:ring-1 focus:ring-indigo-500"
        />
      </div>
      <Editable
        renderElement={renderElement}
        renderLeaf={renderLeaf}
        spellCheck
        className="focus:shadow-outline block min-h-[200px] w-full appearance-none rounded-b-md rounded-tr-md border border-gray-300 py-2 px-3 leading-tight text-gray-700 shadow focus:outline-none"
      />
    </Slate>
  );
}
