import TablerSettings from "~icons/tabler/settings";
import { useAppShell } from "@/components/AppShell/AppShellContext";
import { cn } from "@/utils/cn";
import { useRouter } from "@tanstack/react-router";
import { auth } from "@/lib/authClient";

export function Profile({ name, role }: { name: string; role: string }) {
  const router = useRouter();
  const { setMobileNavOpen } = useAppShell();
  const { data } = auth.useUser();
  const { user } = data!;

  return (
    <div className="flex items-center gap-2">
      <img
        src={
          user?.image ||
          "https://static.vecteezy.com/system/resources/thumbnails/009/292/244/small_2x/default-avatar-icon-of-social-media-user-vector.jpg"
        }
        className="w-8 aspect-square rounded-full"
      />
      <div className="flex justify-between w-full items-center">
        <div className="w-full">
          <p className="text-sm truncate max-w-40">{name}</p>
          <p className="text-sm text-text-secondary opacity-85">{role}</p>
        </div>
        <div className="flex items-center gap-1">
          <ProfileIcon
            onClick={() => {
              router.navigate({ to: "/settings" });
              setMobileNavOpen(false);
            }}
          >
            <TablerSettings />
          </ProfileIcon>
        </div>
      </div>
    </div>
  );
}

function ProfileIcon({
  children,
  className,
  ...props
}: React.PropsWithChildren<React.HTMLAttributes<HTMLDivElement>>) {
  return (
    <div
      {...props}
      className={cn(
        "hover:bg-navlink-bg-active p-1 rounded-md hover:cursor-pointer duration-100 opacity-80",
        className,
      )}
    >
      {children}
    </div>
  );
}
