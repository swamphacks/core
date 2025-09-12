import { useRouter } from "@tanstack/react-router";
import { Button, MenuTrigger } from "react-aria-components";
import TablerSettings from "~icons/tabler/settings";
import { Menu, MenuItem } from "@/components/ui/Menu";
import TablerLogout from "~icons/tabler/logout";

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

  return (
    <div className="flex items-center gap-2">
      <MenuTrigger>
        <Button
          aria-label="Account"
          className="inline-flex items-center justify-center rounded-md p-1.5 text-white bg-transparent border-none hover:bg-gray-200 pressed:bg-gray-300 dark:hover:bg-zinc-800 dark:pressed:bg-zinc-700 transition-colors cursor-default outline-hidden focus-visible:ring-2 focus-visible:ring-blue-600"
        >
          <img
            alt=""
            src="https://i.pinimg.com/736x/8b/d2/f6/8bd2f653f38322972e404925ab67294a.jpg"
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
