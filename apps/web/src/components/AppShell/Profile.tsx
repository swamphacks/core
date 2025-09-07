import TablerSettings from "~icons/tabler/settings";
import TablerLogout from "~icons/tabler/logout";

export function Profile({ name, role }: { name: string; role: string }) {
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
          <div className="hover:bg-navlink-bg-active p-1 rounded-md">
            <TablerSettings />
          </div>
          <div className="hover:bg-navlink-bg-active p-1 rounded-md">
            <TablerLogout className="text-badge-text-rejected" />
          </div>
        </div>
      </div>
    </div>
  );
}
