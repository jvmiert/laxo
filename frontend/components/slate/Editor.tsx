import cc from "classcat";
import React, { useCallback, useMemo, useState } from "react";
import {
  Editable,
  withReact,
  useSlate,
  Slate,
  ReactEditor,
  RenderElementProps,
  RenderLeafProps,
} from "slate-react";
import {
  BaseEditor,
  Editor as SlateEditor,
  Transforms,
  createEditor,
  Descendant,
  Element as SlateElement,
} from "slate";
import { withHistory, HistoryEditor } from "slate-history";

type LaxoElement = {
  type: string;
  children: LaxoText[];
};

type FormatParameter = "bold" | "italic" | "code" | "underline";

type LaxoText = {
  text: string;
  bold?: boolean;
  italic?: boolean;
  code?: boolean;
  underline?: boolean;
};

declare module "slate" {
  interface CustomTypes {
    Editor: BaseEditor & ReactEditor & HistoryEditor;
    Element: LaxoElement;
    Text: LaxoText;
  }
}

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

const Element = ({ attributes, children, element }: RenderElementProps) => {
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

  //@TODO: We might have to move state up because of our Disclosure usage
  const [editor] = useState(() => withHistory(withReact(createEditor())));
  const [value, setValue] = useState(initialSchema ? JSON.parse(initialSchema) : initialValue);

  return (
    <Slate editor={editor} value={value}>
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
        <MarkButton
          text="<>"
          format="code"
          className="relative -ml-px inline-flex items-center rounded-tr-md border-t border-l border-r border-gray-300 px-4 py-2 text-sm font-medium focus:z-10 focus:border-indigo-500 focus:outline-none focus:ring-1 focus:ring-indigo-500"
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
