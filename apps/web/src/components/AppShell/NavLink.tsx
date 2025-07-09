import type { PropsWithChildren, ReactNode } from "react";
import TablerChevronRight from "~icons/tabler/chevron-right";
import { Link } from "react-aria-components";
import { tv } from "tailwind-variants";
import { useToggleState } from "react-stately";
import { cn } from "@/utils/cn";

const navLink = tv({
  base: "px-3 py-2.5 rounded-sm text-sm flex flex-row items-center justify-between w-auto cursor-pointer transition-none",
  variants: {
    active: {
      true: "bg-navlink-bg-active text-navlink-text-active font-medium",
      false:
        "bg-navlink-bg-inactive text-navlink-text-inactive font-normal hover:bg-neutral-200 dark:hover:bg-neutral-800",
    },
  },
});

interface NavLinkProps {
  href?: string;
  label?: string;
  leftSection?: ReactNode;
  rightSection?: ReactNode;
  active?: boolean;
  initialExpanded?: boolean;
}

const NavLink = ({
  href,
  label,
  leftSection,
  rightSection,
  active = false,
  initialExpanded = false,
  children,
}: PropsWithChildren<NavLinkProps>) => {
  const isExpandable = !!children;

  const toggleState = useToggleState({ defaultSelected: initialExpanded });

  const toggle = () => {
    console.log("Toggling nav link");
    console.log("Current state:", toggleState.isSelected);
    toggleState.toggle();
    console.log("New state:", toggleState.isSelected);
  };

  return (
    <div>
      {isExpandable ? (
        <Link className={navLink({ active })} onPress={toggle}>
          <div className="flex flex-row gap-2 items-center select-none">
            {leftSection && (
              <span className="flex items-center justify-center">
                {leftSection}
              </span>
            )}
            <span>{label}</span>
          </div>

          <span
            className={cn(
              "flex items-center justify-center transition-transform duration-150",
              toggleState.isSelected ? "rotate-90" : "rotate-0",
            )}
          >
            <TablerChevronRight className="w-4 aspect-square" />
          </span>
        </Link>
      ) : (
        <Link className={navLink({ active })} href={href}>
          <div className="flex flex-row gap-2 items-center">
            {leftSection && (
              <span className="flex items-center justify-center">
                {leftSection}
              </span>
            )}
            <span>{label}</span>
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
