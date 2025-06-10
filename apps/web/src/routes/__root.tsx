import { ThemeSwitch } from "@/components/ThemeProvider";
import { createRootRouteWithContext, Outlet } from "@tanstack/react-router";

interface RouterContext {
  user: undefined; // Fill this in with your user type
}

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
