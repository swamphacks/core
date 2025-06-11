import { ThemeSwitch } from "@/components/ThemeProvider";
import type { authClient } from "@/lib/authClient";
import { createRootRouteWithContext, Outlet } from "@tanstack/react-router";

interface RouterContext {
  user: ReturnType<typeof authClient.useUser>;
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
