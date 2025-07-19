import { cn } from "@/utils/cn";
import type { ComponentType } from "react";
import { composeRenderProps } from "react-aria-components";

export type Icon = ComponentType<React.SVGProps<SVGSVGElement>>;

export function composeTailwindRenderProps<T>(
  className: string | ((v: T) => string) | undefined,
  tw: string,
): string | ((v: T) => string) {
  return composeRenderProps(className, (className) => cn(tw, className));
}
