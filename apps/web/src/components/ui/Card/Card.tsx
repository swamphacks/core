import { forwardRef } from "react";
import { tv, type VariantProps } from "tailwind-variants";

// Very basic card component, maybe in the future we could do something like https://chakra-ui.com/docs/components/card

export const card = tv({
  base: "inline-block bg-surface rounded-md w-full sm:max-w-96 shadow-xs",
  variants: {},

  defaultVariants: {},
});

type CardVariants = VariantProps<typeof card>;

export interface CardProps
  extends CardVariants,
    React.HTMLAttributes<HTMLDivElement> {}

const Card = forwardRef<HTMLDivElement, CardProps>(
  ({ className, ...props }, ref) => {
    const cardClassName = card({ className });

    return <div {...props} ref={ref} className={cardClassName} />;
  },
);

Card.displayName = "Card";

export { Card };
