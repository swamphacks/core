import type { JSX } from "react";
import { IconWrapper, type IconWrapperProps } from "./IconWrapper";

// All SVG icons must call this function.
// It allows us to do something like <Icon className="..."/> without requiring us to do define and inject
// props inside the Icon directly.
export const createIcon = (iconSvg: () => JSX.Element) => {
  return (props?: Omit<IconWrapperProps, "children">) => (
    <IconWrapper {...props}>{iconSvg()}</IconWrapper>
  );
};
