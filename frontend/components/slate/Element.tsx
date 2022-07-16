import { RenderElementProps } from "slate-react";
import cc from "classcat";

import { useDashboard } from "@/providers/DashboardProvider";
import ImageElement from "@/components/slate/ImageElement";

export default function Element({
  attributes,
  children,
  element,
}: RenderElementProps) {
  const { activeShop } = useDashboard();
  const token = activeShop ? activeShop.assetsToken : "";

  switch (element.type) {
    case "bulleted-list":
      return (
        <ul className="list-inside list-disc" {...attributes}>
          {children}
        </ul>
      );
    case "heading-one":
      return (
        <h1
          className={cc([
            "text-xl font-semibold",
            { "text-left": element.align === "left" },
            { "text-right": element.align === "right" },
            { "text-center": element.align === "center" },
            { "whitespace-pre-line text-justify": element.align === "justify" },
          ])}
          {...attributes}
        >
          {children}
        </h1>
      );
    case "heading-two":
      return (
        <h2
          className={cc([
            "text-lg font-medium",
            { "text-left": element.align === "left" },
            { "text-right": element.align === "right" },
            { "text-center": element.align === "center" },
            { "whitespace-pre-line text-justify": element.align === "justify" },
          ])}
          {...attributes}
        >
          {children}
        </h2>
      );
    case "list-item":
      return (
        <li
          className={cc([
            { "text-left": element.align === "left" },
            { "text-right": element.align === "right" },
            { "text-center": element.align === "center" },
            { "whitespace-pre-line text-justify": element.align === "justify" },
          ])}
          {...attributes}
        >
          {children}
        </li>
      );
    case "numbered-list":
      return (
        <ol className="list-inside list-decimal" {...attributes}>
          {children}
        </ol>
      );
    case "image":
      return (
        <ImageElement token={token} attributes={attributes} element={element}>
          {children}
        </ImageElement>
      );
    default:
      return (
        <p
          className={cc([
            { "text-left": element.align === "left" },
            { "text-right": element.align === "right" },
            { "text-center": element.align === "center" },
            { "whitespace-pre-line text-justify": element.align === "justify" },
          ])}
          {...attributes}
        >
          {children}
        </p>
      );
  }
}
