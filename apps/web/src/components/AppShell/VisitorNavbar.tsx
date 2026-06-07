import { NavLink } from "@/components/AppShell/NavLink";
import TablerInfoCircle from "~icons/tabler/info-circle";
import TablerClipboard from "~icons/tabler/clipboard";

interface VisitorNavbarProps {
  pathname: string;
}

export default function VisitorNavbar({ pathname }: VisitorNavbarProps) {
  const commonNavLinks = (
    <>
      <NavLink
        label="Information"
        href="/information"
        leftSection={<TablerInfoCircle className="w-5 aspect-square" />}
        active={pathname.startsWith("/information")}
      />

      <NavLink
        label="Application"
        href="/application"
        leftSection={<TablerClipboard className="w-5 aspect-square" />}
        active={pathname.startsWith("/application")}
      />
    </>
  );

  return <>{commonNavLinks}</>;
}
