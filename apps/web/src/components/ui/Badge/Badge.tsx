/* eslint-disable react-refresh/only-export-components */
import { forwardRef } from "react";
import { tv, type VariantProps } from "tailwind-variants";

export const badge = tv({
  base: "inline-flex items-center gap-1 rounded-xl font-medium select-none",
  variants: {
    type: {
      default: "bg-badge-bg-default text-badge-text-default",
    },
    size: {
      sm: "px-2 py-1 text-xs",
      md: "px-2.5 py-1.5 text-sm",
    },
    // TODO: Border radius?
  },

  defaultVariants: {
    size: "sm",
    type: "default",
  },
});

type BadgeVariants = VariantProps<typeof badge>;

export interface BadgeProps
  extends BadgeVariants,
    React.HTMLAttributes<HTMLSpanElement> {}

const Badge = forwardRef<HTMLSpanElement, BadgeProps>(
  ({ size, type, className, ...props }, ref) => {
    const badgeClassName = badge({ size, type, className });

    return <span {...props} ref={ref} className={badgeClassName} />;
  },
);

Badge.displayName = "Badge";

export { Badge };
