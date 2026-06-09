import { SettingsPage } from "@/modules/Settings/SettingsPage";
import { auth } from "@/lib/authClient";
import { createFileRoute, useRouter } from "@tanstack/react-router";
import { userQueryOptions } from "@/lib/auth/hooks/useUser";
import { useSuspenseQuery } from "@tanstack/react-query";

export const Route = createFileRoute("/_protected/settings")({
  component: RouteComponent,
});

function RouteComponent() {
  const router = useRouter();
  // We call useSuspenseQuery to get user data because they can change their information and
  // this hook lets us subscribe to it and rerender whenever a change happens
  const user = useSuspenseQuery(userQueryOptions());

  const logout = async () => {
    await auth.logOut();
    await router.navigate({ to: "/", replace: true, reloadDocument: true }); // Force reload to reset auth state (Not ideal, but works for now)
  };

  return (
    <div className="flex items-center justify-center w-full h-full overflow-y-auto">
      <SettingsPage logout={logout} user={user.data.user} />
    </div>
  );
}
