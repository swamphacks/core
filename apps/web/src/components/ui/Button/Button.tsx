import {
  composeRenderProps,
  Button as RACButton,
  type ButtonProps as RACButtonProps,
} from "react-aria-components";
import { tv } from "tailwind-variants";

export interface ButtonProps extends RACButtonProps {
  variant?: "primary" | "secondary" | "danger" | "icon" | "skeleton";
  size?: "sm" | "md" | "lg" | "auto";
  className?: string;
}

export const button = tv({
  base: "inline-flex cursor-pointer items-center justify-center rounded-md font-medium focus:outline-none",
  variants: {
    variant: {
      skeleton: "",
      primary:
        "bg-button-primary hover:bg-button-primary-hover pressed:bg-button-primary-pressed text-white",
      secondary:
        "bg-button-secondary hover:bg-button-secondary-hover pressed:bg-button-secondary-pressed text-gray-800 dark:bg-zinc-600 dark:hover:bg-zinc-500 dark:pressed:bg-zinc-400 dark:text-zinc-100",
      danger:
        "bg-button-danger hover:bg-button-danger-hover pressed:bg-button-danger-pressed text-white",
      icon: "border-0 p-1 flex items-center justify-center text-gray-600 hover:bg-black/[5%] pressed:bg-black/10 dark:text-zinc-400 dark:hover:bg-white/10 dark:pressed:bg-white/20 disabled:bg-transparent",
    },
    isDisabled: {
      true: "cursor-not-allowed bg-gray-100 dark:bg-zinc-800 text-gray-300 dark:text-zinc-600 forced-colors:text-[GrayText] border-black/5 dark:border-white/5",
    },
    size: {
      auto: "",
      sm: "p-2 text-sm",
      md: "p-2 text-base",
      lg: "p-2 text-base",
    },
  },

  defaultVariants: {
    variant: "primary",
    size: "md",
  },
});

export function Button(props: ButtonProps) {
  return (
    <RACButton
      {...props}
      className={composeRenderProps(props.className, (className, renderProps) =>
        button({
          ...renderProps,
          variant: props.variant,
          size: props.size,
          className,
        }),
      )}
    />
  );
}
