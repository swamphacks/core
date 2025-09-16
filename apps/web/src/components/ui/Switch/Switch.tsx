import React from "react";
import {
  Switch as AriaSwitch,
  type SwitchProps as AriaSwitchProps,
} from "react-aria-components";
import { tv } from "tailwind-variants";
import { composeTailwindRenderProps } from "../utils";

export interface SwitchProps extends Omit<AriaSwitchProps, "children"> {
  children: React.ReactNode;
}

const track = tv({
  base: "flex h-4 w-7 px-px items-center shrink-0 cursor-default rounded-full transition duration-200 ease-in-out shadow-inner border border-transparent",
  variants: {
    isSelected: {
      false:
        "bg-gray-400 dark:bg-zinc-600 group-pressed:bg-gray-500 dark:group-pressed:bg-zinc-700",
      true: "bg-button-primary forced-colors:bg-[Highlight]!",
    },
    isDisabled: {
      true: "bg-gray-200 dark:bg-zinc-700 forced-colors:group-selected:bg-[GrayText]! forced-colors:border-[GrayText]",
    },
  },
});

const handle = tv({
  base: "h-3 w-3 transform rounded-full bg-white outline outline-1 -outline-offset-1 outline-transparent shadow-xs transition duration-200 ease-in-out",
  variants: {
    isSelected: {
      false: "translate-x-0",
      true: "translate-x-[100%]",
    },
    isDisabled: {
      true: "forced-colors:outline-[GrayText]",
    },
  },
});

export function Switch({ children, ...props }: SwitchProps) {
  return (
    <AriaSwitch
      {...props}
      className={composeTailwindRenderProps(
        props.className,
        "group relative flex gap-2 items-center text-text-main/90 disabled:text-gray-300  dark:disabled:text-zinc-600 forced-colors:disabled:text-[GrayText] transition",
      )}
    >
      {(renderProps) => (
        <>
          <div className={track(renderProps)}>
            <span className={handle(renderProps)} />
          </div>
          {children}
        </>
      )}
    </AriaSwitch>
  );
}
