import { DEVThemeSwitch } from "@/components/ThemeProvider";
import type { auth } from "@/lib/authClient";
import { createRootRouteWithContext, Outlet } from "@tanstack/react-router";

interface RouterContext {
  userQuery: ReturnType<typeof auth.useUser>;
}

const IS_DEV = import.meta.env.DEV;

export const Route = createRootRouteWithContext<RouterContext>()({
  component: () => (
    <>
      <Outlet />
      {IS_DEV && (
        <div className="fixed inline-flex w-fit z-[999] bottom-3 left-3 text-white">
          <DEVThemeSwitch />
        </div>
      )}
    </>
  ),
});
