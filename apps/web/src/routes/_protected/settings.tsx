import { SettingsPage } from "@/features/Settings/components/SettingsPage";
import { auth } from "@/lib/authClient";
import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/_protected/settings")({
  component: RouteComponent,
});

function RouteComponent() {
  const logout = async () => {
    await auth.logOut();
    location.reload();
  };

  return (
    <div className="flex items-center justify-center w-full h-full overflow-y-auto">
      <SettingsPage logout={logout} />
    </div>
  );
}
