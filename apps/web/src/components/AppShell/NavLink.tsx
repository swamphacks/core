import type { PropsWithChildren, ReactNode } from "react";
import TablerChevronRight from "~icons/tabler/chevron-right";
import { tv } from "tailwind-variants";
import { useToggleState } from "react-stately";
import { cn } from "@/utils/cn";
import { Button as RAC_Button } from "react-aria-components";
import { Link } from "@tanstack/react-router";

const navLink = tv({
  base: "px-3 py-2.5 rounded-sm text-sm flex flex-row items-center justify-between w-full cursor-pointer transition-none select-none",
  variants: {
    active: {
      true: "bg-navlink-bg-active text-navlink-text-active font-medium",
      false:
        "bg-navlink-bg-inactive text-navlink-text-inactive font-normal hover:scale-101",
    },
  },
});

interface NavLinkProps {
  href?: string;
  label: string;
  description?: string;
  leftSection?: ReactNode;
  rightSection?: ReactNode;
  active?: boolean;
  initialExpanded?: boolean;
}

const NavLink = ({
  href,
  label,
  description,
  leftSection,
  rightSection,
  active = false,
  initialExpanded = false,
  children,
}: PropsWithChildren<NavLinkProps>) => {
  const isExpandable = !!children;

  const toggleState = useToggleState({ defaultSelected: initialExpanded });
  const toggle = toggleState.toggle;

  return (
    <div>
      {isExpandable ? (
        <RAC_Button className={navLink({ active })} onPress={toggle}>
          <div className="flex flex-row gap-2 items-center">
            {leftSection && (
              <span className="flex items-center justify-center">
                {leftSection}
              </span>
            )}
            <div className="flex flex-col gap-0.5">
              <span>{label}</span>
              {description && (
                <span className="text-xs text-neutral-500">{description}</span>
              )}
            </div>
          </div>

          <span
            className={cn(
              "flex items-center justify-center transition-transform duration-150",
              toggleState.isSelected ? "rotate-90" : "rotate-0",
            )}
          >
            <TablerChevronRight className="w-4 aspect-square" />
          </span>
        </RAC_Button>
      ) : (
        <Link className={navLink({ active })} to={href}>
          <div className="flex flex-row gap-2 items-center">
            {leftSection && (
              <span className="flex items-center justify-center">
                {leftSection}
              </span>
            )}
            <div className="flex flex-col gap-0.5">
              <span>{label}</span>
              {description && (
                <span className="text-xs text-neutral-500">{description}</span>
              )}
            </div>
          </div>
          {rightSection && (
            <span className="flex items-center justify-center">
              {rightSection}
            </span>
          )}
        </Link>
      )}

      {/* Expanded children links */}
      <div>
        {isExpandable && toggleState.isSelected && (
          <div className="pl-6 mt-1">{children}</div>
        )}
      </div>
    </div>
  );
};

NavLink.displayName = "NavLink";

export { NavLink };
