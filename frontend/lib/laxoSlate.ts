import { ReactEditor } from "slate-react";
import { BaseEditor, Editor } from "slate";
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
