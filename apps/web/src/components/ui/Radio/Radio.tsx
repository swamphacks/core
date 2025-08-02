import { Radio as RACRadio, type RadioProps } from "react-aria-components";
import { tv } from "tailwind-variants";
import { composeTailwindRenderProps } from "@/components/ui/utils";

const styles = tv({
  base: "w-5 h-5 rounded-full border-2 bg-white dark:bg-zinc-900 transition-all",
  variants: {
    isSelected: {
      false:
        "border-border dark:border-gray-400 group-pressed:border-gray-400 dark:group-pressed:border-zinc-300",
      true: "border-[7px] border-blue-600 dark:border-blue-300 forced-colors:border-[Highlight]! group-pressed:border-blue-800 dark:group-pressed:border-blue-200",
    },
    isInvalid: {
      true: "border-red-700 dark:border-red-600 group-pressed:border-red-800 dark:group-pressed:border-red-700 forced-colors:border-[Mark]!",
    },
    isDisabled: {
      true: "border-gray-200 dark:border-zinc-700 forced-colors:border-[GrayText]!",
    },
  },
});

export function Radio(props: RadioProps) {
  return (
    <RACRadio
      {...props}
      className={composeTailwindRenderProps(
        props.className,
        "flex relative gap-2 items-center group text-text-main text-sm disabled:text-gray-300  dark:disabled:text-zinc-600 forced-colors:disabled:text-[GrayText] transition",
      )}
    >
      {(renderProps) => (
        <>
          <div className={styles(renderProps)} />
          {props.children}
        </>
      )}
    </RACRadio>
  );
}
