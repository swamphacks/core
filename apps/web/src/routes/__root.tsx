/* eslint-disable @typescript-eslint/no-empty-object-type */
import { ThemeSwitch } from "@/components/ThemeProvider";
import { createRootRouteWithContext, Outlet } from "@tanstack/react-router";

interface RouterContext {}

const IS_DEV = import.meta.env.DEV;

export const Route = createRootRouteWithContext<RouterContext>()({
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
