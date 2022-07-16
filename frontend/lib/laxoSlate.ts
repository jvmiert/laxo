import { ReactEditor } from "slate-react";
import { BaseEditor, Editor, Element, Transforms } from "slate";
import { HistoryEditor } from "slate-history";

export type LaxoParagraphElement = {
  type: "paragraph";
  children: LaxoText[];
  align?: "left" | "center" | "right" | "justify";
};

export type LaxoHeadingElement = {
  type: "heading-one" | "heading-two" | "heading-three";
  children: LaxoText[];
  align?: "left" | "center" | "right" | "justify";
};

export type LaxoListItemElement = {
  type: "list-item";
  children: LaxoText[];
  align?: "left" | "center" | "right" | "justify";
};

//@TODO: align shouldn't be here
export type LaxoListElement = {
  type: "bulleted-list" | "numbered-list";
  children: LaxoListItemElement[];
  align?: "left" | "center" | "right" | "justify";
};

export type LaxoBlockElements = Exclude<LaxoElement, LaxoImageElement>;

//@TODO: align and children shouldn't be here
export type LaxoImageElement = {
  type: "image";
  src: string;
  width: number;
  height: number;
  align?: "left" | "center" | "right" | "justify";
  children: LaxoText[];
};

type LaxoElement =
  | LaxoImageElement
  | LaxoListElement
  | LaxoListItemElement
  | LaxoHeadingElement
  | LaxoParagraphElement;

export type FormatParameter = "bold" | "italic" | "code" | "underline";
export type BlockFormatParameter =
  | "left"
  | "center"
  | "right"
  | "justify"
  | "heading-one"
  | "heading-two"
  | "heading-three"
  | "bulleted-list"
  | "numbered-list";

type LaxoText = {
  text: string;
  bold?: boolean;
  italic?: boolean;
  code?: boolean;
  underline?: boolean;
};

export const TEXT_ALIGN_TYPES = ["left", "center", "right", "justify"];
export const LIST_TYPES = ["numbered-list", "bulleted-list"];

declare module "slate" {
  interface CustomTypes {
    Editor: BaseEditor & ReactEditor & HistoryEditor;
    Element: LaxoElement;
    Text: LaxoText;
  }
}

export const withImages = (editor: Editor) => {
  const { isVoid } = editor;

  editor.isVoid = (element) => {
    return element.type === "image" ? true : isVoid(element);
  };

  return editor;
};

export function isBlockActive(
  editor: Editor,
  format: BlockFormatParameter,
  blockType: "type" | "align" = "type",
) {
  const { selection } = editor;
  if (!selection) return false;

  const [match] = Array.from(
    Editor.nodes(editor, {
      at: Editor.unhangRange(editor, selection),
      match: (n) =>
        !Editor.isEditor(n) && Element.isElement(n) && n[blockType] === format,
    }),
  );

  return !!match;
}

export const toggleBlock = (editor: Editor, format: BlockFormatParameter) => {
  const isActive = isBlockActive(
    editor,
    format,
    TEXT_ALIGN_TYPES.includes(format) ? "align" : "type",
  );
  const isList = LIST_TYPES.includes(format);

  Transforms.unwrapNodes(editor, {
    match: (n) =>
      !Editor.isEditor(n) &&
      Element.isElement(n) &&
      LIST_TYPES.includes(n.type) &&
      !TEXT_ALIGN_TYPES.includes(format),
    split: true,
  });
  let newProperties: Partial<Element>;
  if (TEXT_ALIGN_TYPES.includes(format)) {
    newProperties = {
      align: isActive
        ? undefined
        : (format as "left" | "center" | "right" | "justify"),
    };
  } else {
    newProperties = {
      type: isActive
        ? "paragraph"
        : isList
        ? "list-item"
        : (format as Exclude<
            BlockFormatParameter,
            "left" | "center" | "right" | "justify"
          >),
    };
  }

  Transforms.setNodes<Element>(editor, newProperties);

  if (!isActive && isList) {
    const block = {
      type: format as "numbered-list" | "bulleted-list",
      children: [],
    };
    Transforms.wrapNodes(editor, block);
  }
};

export const isMarkActive = (editor: Editor, format: FormatParameter) => {
  const marks = Editor.marks(editor);
  return marks ? marks[format] === true : false;
};

export const toggleMark = (editor: Editor, format: FormatParameter) => {
  const isActive = isMarkActive(editor, format);

  if (isActive) {
    Editor.removeMark(editor, format);
  } else {
    Editor.addMark(editor, format, true);
  }
};
