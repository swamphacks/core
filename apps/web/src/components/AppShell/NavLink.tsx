import type { ReactNode } from "react";
import { Link } from "react-aria-components";
import { tv } from "tailwind-variants";

const navLink = tv({
  base: "px-3 py-3 transition-colors duration-200 rounded-sm text-lg flex flex-row items-center justify-between w-auto",
  variants: {
    active: {
      true: "bg-navlink-bg-active text-navlink-text-active font-medium",
      false: "bg-navlink-bg-inactive text-navlink-text-inactive font-normal",
    },
  },
});

interface NavLinkProps {
  href?: string;
  label?: string;
  leftSection?: ReactNode;
  rightSection?: ReactNode;
  active?: boolean;
}

const NavLink = ({
  href,
  label,
  leftSection,
  rightSection,
  active = false,
}: NavLinkProps) => {
  return (
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
        <span className="flex items-center justify-center">{rightSection}</span>
      )}
    </Link>
  );
};

NavLink.displayName = "NavLink";

export { NavLink };
