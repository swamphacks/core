import { NavLink } from "@/components/AppShell/NavLink";
import TablerCode from "~icons/tabler/code";
import TablerClipboard from "~icons/tabler/clipboard";
import TablerAlertCircleFilled from "~icons/tabler/alert-circle-filled";

interface ApplicantNavbarProps {
  pathname: string;
  hasSeenNewApplicationStatus: boolean | null;
}

export default function ApplicantNavbar({
  pathname,
  hasSeenNewApplicationStatus,
}: ApplicantNavbarProps) {
  const commonNavLinks = (
    <>
      <NavLink
        label="Information"
        href="/information"
        leftSection={<TablerCode className="w-5 aspect-square" />}
        active={pathname.startsWith("/information")}
      />

      <NavLink
        label="Application"
        href="/application"
        leftSection={<TablerClipboard className="w-5 aspect-square" />}
        rightSection={
          hasSeenNewApplicationStatus === false && (
            <TablerAlertCircleFilled className="text-orange-400" />
          )
        }
        active={pathname.startsWith("/application")}
      />
    </>
  );
  return commonNavLinks;
}
