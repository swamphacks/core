import { useLocation, useRouter } from "@tanstack/react-router";
import { Button, MenuTrigger } from "react-aria-components";
import TablerSettings from "~icons/tabler/settings";
import { Menu, MenuItem } from "@/components/ui/Menu";
import TablerLogout from "~icons/tabler/logout";
import { auth } from "@/lib/authClient";
import TablerUserCog from "~icons/tabler/user-cog";
import TablerUser from "~icons/tabler/user";

export function MobileProfile({
  name,
  role,
  logout,
}: {
  name: string;
  role: string;
  logout: () => void;
}) {
  const router = useRouter();
  const { data } = auth.useUser();
  const { user } = data!;
  const pathname = useLocation({ select: (loc) => loc.pathname });

  const isAdminPortal = pathname.startsWith("/admin");

  return (
    <div className="flex items-center gap-2">
      <MenuTrigger>
        <Button
          aria-label="Account"
          className="inline-flex items-center justify-center rounded-md p-1.5 text-white bg-transparent border-none hover:bg-gray-200 pressed:bg-gray-300 dark:hover:bg-zinc-800 dark:pressed:bg-zinc-700 transition-colors cursor-default outline-hidden focus-visible:ring-2 focus-visible:ring-blue-600"
        >
          <img
            alt=""
            src={
              user?.image ||
              "https://static.vecteezy.com/system/resources/thumbnails/009/292/244/small_2x/default-avatar-icon-of-social-media-user-vector.jpg"
            }
            className="w-9 h-9 rounded-full"
          />
        </Button>

        <Menu
          className="outline-hidden"
          header={
            <div className="pt-2 pl-4">
              <p className="text-sm truncate max-w-35">{name}</p>
              <p className="text-sm text-text-secondary opacity-85">{role}</p>
            </div>
          }
        >
          <MenuItem
            id="portal"
            onAction={() =>
              router.navigate({
                to: isAdminPortal ? "/portal" : "/admin/overview",
              })
            }
          >
            {isAdminPortal ? <TablerUser /> : <TablerUserCog />}
            {isAdminPortal ? "User" : "Admin"} Portal
          </MenuItem>
          <MenuItem
            id="settings"
            onAction={() => router.navigate({ to: "/settings" })}
          >
            <TablerSettings />
            Account Settings
          </MenuItem>
          <MenuItem
            className="text-badge-text-rejected"
            id="logout"
            onAction={logout}
          >
            <TablerLogout />
            Log Out
          </MenuItem>
        </Menu>
      </MenuTrigger>
    </div>
  );
}
