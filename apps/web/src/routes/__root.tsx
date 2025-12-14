import { ThemeSwitch } from "@/components/ThemeProvider";
import NotFoundPage from "@/features/NotFound/NotFoundPage";
import type { auth } from "@/lib/authClient";
import type { QueryClient } from "@tanstack/react-query";
import { createRootRouteWithContext, Outlet } from "@tanstack/react-router";

interface RouterContext {
  userQuery: ReturnType<typeof auth.useUser>;
  queryClient: QueryClient;
}

const IS_DEV = import.meta.env.DEV;

export const Route = createRootRouteWithContext<RouterContext>()({
  component: () => (
    <>
      <Outlet />
      {IS_DEV && (
        <div className="fixed inline-flex w-fit z-[999] bottom-3 right-3 text-white">
          <ThemeSwitch />
        </div>
      )}
    </>
  ),
  notFoundComponent: NotFoundPage,
});
