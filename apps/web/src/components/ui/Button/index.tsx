import { forwardRef } from "react";
// RAC = React Aria Components
import {
  Button as RAC_Button,
  type ButtonProps as RAC_ButtonProps,
} from "react-aria-components";
import { tv, type VariantProps } from "tailwind-variants";

// Define styles with tailwind-variants
const button = tv({
  base: "inline-flex cursor-pointer px-4 py-2 items-center justify-center rounded-md text-sm font-medium focus:outline-none focus:ring-2 focus:ring-offset-2",
  variants: {
    color: {
      primary:
        "bg-button-primary text-white hover:bg-button-primary-hover focus:ring-blue-500",
      secondary: "bg-gray-600 text-white hover:bg-gray-700 focus:ring-gray-500",
      danger: "bg-red-600 text-white hover:bg-red-700 focus:ring-red-500",
    },
  },
  defaultVariants: {
    color: "primary",
  },
});

// Combine Tailwind variant types with RAC props
type ButtonVariants = VariantProps<typeof button>;

interface ButtonProps extends ButtonVariants, RAC_ButtonProps {}

const Button = forwardRef<HTMLButtonElement, ButtonProps>(
  ({ color, className, ...props }, ref) => {
    // Combine the props to apply styles from tailwind-variants along with React Aria button props
    const btnClassName = button({ color });

    return (
      <RAC_Button
        {...props}
        ref={ref}
        className={`${btnClassName} ${className}`}
      />
    );
  },
);

Button.displayName = "Button";

export { Button };
