import {
  composeRenderProps,
  Button as RACButton,
  type ButtonProps as RACButtonProps,
} from "react-aria-components";
import { tv } from "tailwind-variants";

export interface ButtonProps extends RACButtonProps {
  variant?:
    | "primary"
    | "secondary"
    | "danger"
    | "icon"
    | "skeleton"
    | "unstyled";
  size?: "sm" | "md" | "lg" | "auto";
  className?: string;
}

export const button = tv({
  base: "inline-flex cursor-pointer items-center justify-center rounded-md font-medium focus:outline-none gap-2",
  variants: {
    variant: {
      skeleton: "",
      primary:
        "bg-button-primary hover:bg-button-primary-hover pressed:bg-button-primary-pressed text-white",
      secondary:
        "bg-button-secondary hover:bg-button-secondary-hover pressed:bg-button-secondary-pressed",
      danger:
        "bg-button-danger hover:bg-button-danger-hover pressed:bg-button-danger-pressed text-white",
      icon: "border-0 p-1 flex items-center justify-center text-gray-600 hover:bg-black/[5%] pressed:bg-black/10 dark:text-zinc-400 dark:hover:bg-white/10 dark:pressed:bg-white/20 disabled:bg-transparent",
      unstyled: "",
    },
    isDisabled: {
      true: "cursor-not-allowed bg-gray-200 dark:bg-neutral-700 text-text-main/30 border-black/5 dark:border-white/5",
    },
    size: {
      auto: "",
      sm: "py-2 px-4 text-sm",
      md: "py-2 px-4 text-base",
      lg: "py-2 px-4 text-lg",
    },
  },
  compoundVariants: [
    {
      variant: "unstyled",
      class: "", // override everything
    },
  ],
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
          size: props.variant === "unstyled" ? undefined : props.size,
          className,
        }),
      )}
    />
  );
}
