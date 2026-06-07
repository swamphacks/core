import { NavLink } from "@/components/AppShell/NavLink";
import TablerLayoutDashboard from "~icons/tabler/layout-dashboard";
import TablerInfoCircle from "~icons/tabler/info-circle";

interface AttendeeNavbarProps {
  pathname: string;
}

export default function AttendeeNavbar({ pathname }: AttendeeNavbarProps) {
  const commonNavLinks = (
    <>
      <NavLink
        label="Information"
        href="/information"
        leftSection={<TablerInfoCircle className="w-5 aspect-square" />}
        active={pathname.startsWith("/information")}
      />
      <NavLink
        label="Hacker Portal"
        href="/hacker-portal"
        leftSection={<TablerLayoutDashboard className="w-5 aspect-square" />}
        active={pathname.startsWith("/hacker-portal")}
      />
    </>
  );
  return commonNavLinks;
}
