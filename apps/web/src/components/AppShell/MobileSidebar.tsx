import type { PropsWithChildren } from "react";
import { useAppShell } from "./AppShellContext";

interface MobileSidebarProps extends PropsWithChildren {
  isOpen: boolean;
}

const SlideoutNavbar = ({ children, isOpen }: MobileSidebarProps) => {
  const { setMobileNavOpen } = useAppShell();

  return (
    <>
      {/* Overlay */}
      <div
        onClick={() => setMobileNavOpen(false)}
        className={`fixed top-16 left-0 right-0 bottom-0 transition-colors duration-300 ease-in-out z-40 ${
          isOpen ? "bg-black/20" : "bg-black/0 pointer-events-none"
        }`}
      />

      {/* Side Navigation */}
      <aside
        className={`fixed top-16 left-0 h-full w-3/4 bg-background shadow-xl transform transition-transform duration-300 ease-in-out z-50 ${
          isOpen ? "translate-x-0" : "-translate-x-full"
        }`}
      >
        <nav className="flex flex-col gap-2 px-2 py-3 h-[92%]">{children}</nav>
      </aside>
    </>
  );
};

export { SlideoutNavbar };
