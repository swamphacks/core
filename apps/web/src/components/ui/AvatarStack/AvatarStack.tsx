import * as React from "react";
import { tv } from "tailwind-variants";
import { cn } from "@/utils/cn";
import { Avatar, type AvatarProps } from "../Avatar/Avatar";

export interface AvatarStackProps extends React.HTMLAttributes<HTMLDivElement> {
  avatars: (AvatarProps & { id?: string | number | string })[];
  size?: "xs" | "sm" | "md" | "lg" | "xl";
  overlap?: "tight" | "medium" | "loose";
  max?: number; // Maximum avatars to show, show +N for the rest
  bordered?: boolean;
}

export const avatarStack = tv({
  base: "flex items-center",
  variants: {
    overlap: {
      tight: "-space-x-4",
      medium: "-space-x-3",
      loose: "-space-x-2",
    },
  },
  defaultVariants: {
    overlap: "medium",
  },
});

export function AvatarStack({
  avatars,
  size = "md",
  overlap = "medium",
  max,
  bordered = true,
  className,
  ...props
}: AvatarStackProps) {
  const displayAvatars = max ? avatars.slice(0, max) : avatars;
  const remainingCount = max && avatars.length > max ? avatars.length - max : 0;

  return (
    <div {...props} className={cn(avatarStack({ overlap }), className)}>
      {displayAvatars.map((avatar) => (
        <Avatar
          key={avatar.id || avatar.src || avatar.fallback}
          size={size}
          bordered={bordered}
          {...avatar}
        />
      ))}

      {remainingCount > 0 && (
        <div
          className={cn(
            "flex items-center justify-center bg-surface text-text-secondary font-medium rounded-full",
            {
              "w-6 h-6 text-xs": size === "xs",
              "w-8 h-8 text-xs": size === "sm",
              "w-10 h-10 text-sm": size === "md",
              "w-12 h-12 text-base": size === "lg",
              "w-16 h-16 text-lg": size === "xl",
            },
            bordered &&
              "ring-1 ring-border",
          )}
          style={{ zIndex: displayAvatars.length + 1 }}
        >
          +{remainingCount}
        </div>
      )}
    </div>
  );
}
