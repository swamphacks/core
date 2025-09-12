import TablerSettings from "~icons/tabler/settings";
import { useAppShell } from "@/components/AppShell/AppShellContext";
import { cn } from "@/utils/cn";
import { useRouter } from "@tanstack/react-router";

export function Profile({ name, role }: { name: string; role: string }) {
  const router = useRouter();
  const { setMobileNavOpen } = useAppShell();

  return (
    <div className="flex items-center gap-2">
      <img
        src="https://i.pinimg.com/736x/8b/d2/f6/8bd2f653f38322972e404925ab67294a.jpg"
        className="w-8 aspect-square rounded-full"
      />
      <div className="flex justify-between w-full items-center">
        <div>
          <p className="text-sm">{name}</p>
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
