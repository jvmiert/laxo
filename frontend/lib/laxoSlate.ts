import { ReactEditor } from "slate-react";
import { BaseEditor, Editor } from "slate";
import { HistoryEditor } from "slate-history";

type LaxoElement = {
  type:
    | "paragraph"
    | "image"
    | "heading-one"
    | "heading-two"
    | "heading-three"
    | "bulleted-list"
    | "list-item"
    | "numbered-list";
  children?: LaxoText[];
  src?: string;
  width?: number;
  height?: number;
};

export type FormatParameter = "bold" | "italic" | "code" | "underline";

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
