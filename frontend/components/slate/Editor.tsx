import isHotkey from "is-hotkey";
import diff from "microdiff";
import React, { useCallback, useMemo, useEffect, useState } from "react";
import {
  Editable,
  withReact,
  Slate,
  RenderElementProps,
  RenderLeafProps,
} from "slate-react";
import {
  Editor as SlateEditor,
  createEditor,
  Descendant,
  Transforms,
} from "slate";
import { withHistory } from "slate-history";
import {
  AlignJustify,
  AlignLeft,
  AlignRight,
  AlignCenter,
  ListOrdered,
  List,
} from "lucide-react";
import { useIntl } from "react-intl";

import { FormatParameter, withImages, toggleMark } from "@/lib/laxoSlate";
import { useDashboard } from "@/providers/DashboardProvider";
import Leaf from "@/components/slate/Leaf";
import Element from "@/components/slate/Element";
import AssetButton from "@/components/slate/AssetButton";
import BlockButton from "@/components/slate/BlockButton";
import MarkButton from "@/components/slate/MarkButton";

const HOTKEYS: {
  [key: string]: FormatParameter;
} = {
  "mod+b": "bold",
  "mod+i": "italic",
  "mod+u": "underline",
};

const initialValue: Descendant[] = [
  {
    type: "paragraph",
    children: [{ text: "" }],
  },
];

export type EditorProps = {
  initialSchema?: string;
  trackDirty?: boolean;
};

export default function Editor({
  initialSchema = "",
  trackDirty = true,
}: EditorProps) {
  const t = useIntl();

  const renderElement = useCallback(
    (props: RenderElementProps) => <Element {...props} />,
    [],
  );
  const renderLeaf = useCallback(
    (props: RenderLeafProps) => <Leaf {...props} />,
    [],
  );

  const {
    slateResetRef,
    dashboardDispatch,
    toggleSlateDirtyState,
    slateIsDirty,
    slateEditorRef,
  } = useDashboard();

  const openImageInsert = () => {
    dashboardDispatch({
      type: "open_image_insert",
    });
  };

  const [editor] = useState(() =>
    withImages(withHistory(withReact(createEditor()))),
  );
  //https://github.com/ianstormtaylor/slate/issues/3477
  //@TODO: figure out if we need above
  //const editor = useMemo(
  //  () => withImages(withHistory(withReact(createEditor()))),
  //  [],
  //);

  //@TODO: We should probably check if the initialSchema is empty
  // and handle that
  const slateValue = useMemo(() => {
    let parsedSchema;
    try {
      parsedSchema = JSON.parse(initialSchema);
    } catch (e) {
      console.log("json error", e);
      parsedSchema = initialValue;
    }
    // Slate throws an error if the value on the initial render is invalid
    // so we directly set the value on the editor in order
    // to be able to trigger normalization on the initial value before rendering
    editor.children = parsedSchema;
    try {
      SlateEditor.normalize(editor, { force: true });
    } catch (e) {
      console.log("normalize error", e);
      editor.children = initialValue;
      SlateEditor.normalize(editor, { force: true });
    }
    Transforms.deselect(editor);
    // We return the normalized internal value so that the rendering can take over from here
    return editor.children;
  }, [editor, initialSchema]);

  const resetEditor = useCallback(() => {
    editor.children = slateValue;
    Transforms.deselect(editor);
  }, [editor, slateValue]);

  useEffect(() => {
    if (trackDirty) {
      slateResetRef.current = resetEditor;
    }
  }, [resetEditor, slateResetRef, trackDirty]);

  useEffect(() => {
    if (trackDirty) {
      slateEditorRef.current = editor;
    }
  }, [editor, slateEditorRef, trackDirty]);

  const onEditorChange = (d: Descendant[]) => {
    if (trackDirty && !slateIsDirty && diff(slateValue, d).length > 0) {
      toggleSlateDirtyState();
    }
  };

  return (
    <Slate editor={editor} value={slateValue} onChange={onEditorChange}>
      <div className="relative z-0 inline-flex rounded-md shadow-sm">
        <MarkButton
          text="B"
          format="bold"
          className="relative inline-flex items-center rounded-tl-md border-t border-l border-r px-4 py-2 text-sm font-bold focus:z-10 focus:border-indigo-500 focus:outline-none focus:ring-1 focus:ring-indigo-500"
        />
        <MarkButton
          text="I"
          format="italic"
          className="relative -ml-px inline-flex items-center border-t border-l border-r px-4 py-2 font-serif text-sm font-bold italic focus:z-10 focus:border-indigo-500 focus:outline-none focus:ring-1 focus:ring-indigo-500"
        />
        <MarkButton
          text="U"
          format="underline"
          className="relative -ml-px inline-flex items-center border-t border-l border-r px-4 py-2 font-serif text-sm font-bold italic focus:z-10 focus:border-indigo-500 focus:outline-none focus:ring-1 focus:ring-indigo-500"
        />
        <BlockButton
          icon={
            <span>
              H<sub className="text-[9px]">1</sub>
            </span>
          }
          format="heading-one"
          className="relative -ml-px inline-flex items-center border-t border-l border-r px-4 py-2 text-sm font-medium focus:z-10 focus:border-indigo-500 focus:outline-none focus:ring-1 focus:ring-indigo-500"
        />
        <BlockButton
          icon={
            <span>
              H<sub className="text-[9px]">2</sub>
            </span>
          }
          format="heading-two"
          className="relative -ml-px inline-flex items-center border-t border-l border-r px-4 py-2 text-sm font-medium focus:z-10 focus:border-indigo-500 focus:outline-none focus:ring-1 focus:ring-indigo-500"
        />
        <BlockButton
          icon={<ListOrdered className="h-4 w-4" />}
          format="numbered-list"
          className="relative -ml-px inline-flex items-center border-t border-l border-r px-4 py-2 text-sm font-medium focus:z-10 focus:border-indigo-500 focus:outline-none focus:ring-1 focus:ring-indigo-500"
        />
        <BlockButton
          icon={<List className="h-4 w-4" />}
          format="bulleted-list"
          className="relative -ml-px inline-flex items-center border-t border-l border-r px-4 py-2 text-sm font-medium focus:z-10 focus:border-indigo-500 focus:outline-none focus:ring-1 focus:ring-indigo-500"
        />
        <BlockButton
          icon={<AlignLeft className="h-4 w-4" />}
          format="left"
          className="relative -ml-px inline-flex items-center border-t border-l border-r px-4 py-2 text-sm font-medium focus:z-10 focus:border-indigo-500 focus:outline-none focus:ring-1 focus:ring-indigo-500"
        />
        <BlockButton
          icon={<AlignCenter className="h-4 w-4" />}
          format="center"
          className="relative -ml-px inline-flex items-center border-t border-l border-r px-4 py-2 text-sm font-medium focus:z-10 focus:border-indigo-500 focus:outline-none focus:ring-1 focus:ring-indigo-500"
        />
        <BlockButton
          icon={<AlignRight className="h-4 w-4" />}
          format="right"
          className="relative -ml-px inline-flex items-center border-t border-l border-r px-4 py-2 text-sm font-medium focus:z-10 focus:border-indigo-500 focus:outline-none focus:ring-1 focus:ring-indigo-500"
        />
        <BlockButton
          icon={<AlignJustify className="h-4 w-4" />}
          format="justify"
          className="relative -ml-px inline-flex items-center border-t border-l border-r px-4 py-2 text-sm font-medium focus:z-10 focus:border-indigo-500 focus:outline-none focus:ring-1 focus:ring-indigo-500"
        />
        <AssetButton openFunc={openImageInsert} />
      </div>
      <Editable
        renderElement={renderElement}
        renderLeaf={renderLeaf}
        spellCheck
        className="focus:shadow-outline block min-h-[200px] w-full appearance-none rounded-b-md rounded-tr-md border py-2 px-3 leading-tight text-gray-700 shadow focus:outline-none"
        onKeyDown={(event) => {
          for (const hotkey in HOTKEYS) {
            if (isHotkey(hotkey, event as any)) {
              event.preventDefault();
              const mark = HOTKEYS[hotkey];
              toggleMark(editor, mark);
            }
          }
        }}
      />
    </Slate>
  );
}
