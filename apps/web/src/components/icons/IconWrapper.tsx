import { cn } from "@/utils/cn";
import type { JSX } from "react";
import React from "react";

export type IconWrapperProps = React.HTMLAttributes<SVGElement> & {
  children: JSX.Element;
};

// Simplified version of https://www.jacobparis.com/content/react-as-child
// This wrapper is used to wrap any SVG icon which allows us to inject additional props to those icons (like className)
// without requiring us to define those props on the icons themselves which reduces code duplication.
const IconWrapper = ({ children, ...props }: IconWrapperProps) => {
  if (!React.isValidElement(children)) {
    throw new Error("Invalid icon component passed to IconWrapper.");
  }

  if (children.type !== "svg") {
    throw new Error(
      `Icon must be an svg element. The erroneous icon has type \`${children.type}\``,
    );
  }

  return React.cloneElement(children, {
    ...props,
    className: cn(children.props.className, props.className), // custom props will override the icon's
  });
};

export { IconWrapper };
