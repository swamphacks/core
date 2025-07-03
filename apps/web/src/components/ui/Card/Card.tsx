/* eslint-disable react-refresh/only-export-components */
import { forwardRef } from "react";
import { tv, type VariantProps } from "tailwind-variants";

// Very basic card component, maybe in the future we could do something like https://chakra-ui.com/docs/components/card

export const card = tv({
  base: "inline-block bg-surface rounded-md min-w-96 shadow-xs",
  variants: {
    size: {
      sm: "px-2 py-1",
      md: "px-4 py-3",
    },
  },

  defaultVariants: {
    size: "sm",
  },
});

type CardVariants = VariantProps<typeof card>;

export interface CardProps
  extends CardVariants,
    React.HTMLAttributes<HTMLDivElement> {}

const Card = forwardRef<HTMLDivElement, CardProps>(
  ({ size, className, ...props }, ref) => {
    const cardClassName = card({ size, className });

    return <div {...props} ref={ref} className={cardClassName} />;
  },
);

Card.displayName = "Card";

export { Card };
