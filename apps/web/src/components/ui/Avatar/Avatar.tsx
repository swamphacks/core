import { cn } from "@/utils/cn";
import * as React from "react";
import { tv } from "tailwind-variants";

export interface AvatarProps extends React.HTMLAttributes<HTMLDivElement> {
  src?: string;
  alt?: string;
  size?: "xs" | "sm" | "md" | "lg" | "xl";
  shape?: "circle" | "square" | "rounded";
  bordered?: boolean;
  fallback?: string;
  status?: "online" | "offline" | "busy" | "away";
  className?: string;
}

export const avatar = tv({
  base: "flex shrink-0 items-center justify-center overflow-hidden bg-surface text-text-secondary select-none",
  variants: {
    size: {
      xs: "w-6 h-6 text-[10px]",
      sm: "w-8 h-8 text-xs",
      md: "w-10 h-10 text-sm",
      lg: "w-12 h-12 text-base",
      xl: "w-16 h-16 text-lg",
    },
    shape: {
      circle: "rounded-full",
      rounded: "rounded-xl",
      square: "rounded-none",
    },
    bordered: {
      true: "ring-1 ring-border",
    },
  },
  defaultVariants: {
    size: "md",
    shape: "circle",
  },
});

export function Avatar({
  src,
  alt,
  size,
  shape,
  bordered,
  fallback,
  status,
  className,
  ...props
}: AvatarProps) {
  function getInitials(input: string): string {
    const words = input.trim().split(/\s+/);

    // Get first letters of first two words
    if (words.length >= 2) {
      const first = words[0][0];
      const second = words[1][0];
      return (first + second).toUpperCase();
    }

    // Get first word's letters
    const single = words[0];
    const initials = single.slice(0, 2).toUpperCase();

    return initials;
  }

  const initials = fallback ? getInitials(fallback) : null;

  return (
    <div {...props} className={cn("relative inline-block", className)}>
      <div className={avatar({ size, shape, bordered })}>
        {src ? (
          <img
            src={src}
            alt={alt || "avatar"}
            className="w-full h-full object-cover"
            loading="lazy"
          />
        ) : fallback ? (
          <span className="uppercase font-medium">{initials}</span>
        ) : (
          <img
            src="https://ui-avatars.com/api/?background=0D8ABC&color=fff&name=John+Apple"
            alt="default avatar"
            className="w-full h-full object-cover"
            loading="lazy"
          />
        )}
      </div>

      {status && (
        <span
          className={tv({
            base: "absolute bottom-[0.5px] right-[2px] block rounded-full ring-2 ring-white dark:ring-neutral-900",
            variants: {
              size: {
                xs: "w-1.5 h-1.5",
                sm: "w-2 h-2",
                md: "w-2.5 h-2.5",
                lg: "w-3 h-3",
                xl: "w-4 h-4",
              },
              status: {
                online: "bg-green-500",
                offline: "bg-gray-400",
                busy: "bg-red-500",
                away: "bg-yellow-400",
              },
            },
            defaultVariants: {
              size: "md",
            },
          })({ size, status })}
        />
      )}
    </div>
  );
}
