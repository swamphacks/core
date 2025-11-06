import { useEffect, type PropsWithChildren, type ReactNode } from "react";
import TablerChevronRight from "~icons/tabler/chevron-right";
import { tv } from "tailwind-variants";
import { useToggleState } from "react-stately";
import { cn } from "@/utils/cn";
import { Button as RAC_Button } from "react-aria-components";
import { Link } from "@tanstack/react-router";
import { useAppShell } from "./AppShellContext";

const navLink = tv({
  base: "px-3 py-2.5 rounded-sm text-md flex flex-row items-center justify-between w-full cursor-pointer transition-none select-none text-navlink-text",
  variants: {
    active: {
      true: "bg-navlink-bg-active font-medium",
      false: "bg-navlink-bg-inactive font-normal hover:scale-101",
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
  closeNavbarOnClick?: boolean;
}

const NavLink = ({
  href,
  label,
  description,
  leftSection,
  rightSection,
  active = false,
  initialExpanded = false,
  closeNavbarOnClick = true,
  children,
}: PropsWithChildren<NavLinkProps>) => {
  const isExpandable = !!children;

  const toggleState = useToggleState({ defaultSelected: initialExpanded });
  const toggle = toggleState.toggle;

  // Handle mobile navigation state in tangent with the AppShell context
  const { setMobileNavOpen } = useAppShell();

  useEffect(() => {
    // if the navlink was previously expanded, don't automatically close it whenever `initialExpanded` change to false
    // only the user can manually close the navlink

    const prevInitialExpanded = !initialExpanded;

    // if it was previously closed (false), then open it (set toggle state to true)
    if (!prevInitialExpanded) {
      toggleState.setSelected(true);
    }
  }, [initialExpanded]);

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
                <span className="text-xs text-navlink-secondary-text">
                  {description}
                </span>
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
        <Link
          className={navLink({ active })}
          to={href}
          onClick={() => (closeNavbarOnClick ? setMobileNavOpen(false) : null)}
        >
          <div className="flex flex-row gap-2 items-center">
            {leftSection && (
              <span className="flex items-center justify-center">
                {leftSection}
              </span>
            )}
            <div className="flex flex-col gap-0.5">
              <span>{label}</span>
              {description && (
                <span className="text-xs text-navlink-secondary-text">
                  {description}
                </span>
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
