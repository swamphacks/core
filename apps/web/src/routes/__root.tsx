import { ThemeSwitch } from "@/components/ThemeProvider";
import { createRootRoute, Outlet } from "@tanstack/react-router";

const IS_DEV = import.meta.env.DEV;

export const Route = createRootRoute({
  component: () => (
    <>
      <Outlet />
      {IS_DEV && (
        <div className="fixed inline-flex w-fit z-[999] bottom-3 left-3 text-white">
          <ThemeSwitch />
        </div>
      )}
    </>
  ),
});
