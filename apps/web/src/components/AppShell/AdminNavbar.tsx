import { NavLink } from "@/components/AppShell/NavLink";
import TablerInfoCircle from "~icons/tabler/info-circle";
import TablerClipboardCheck from "~icons/tabler/clipboard-check";
import TablerFileText from "~icons/tabler/file-text";
import TablerChartBarPopular from "~icons/tabler/chart-bar-popular";

interface AdminNavbaPropsProps {
  pathname: string;
}

export default function AdminNavbaProps({ pathname }: AdminNavbaPropsProps) {
  const applicationReviewActive = /^\/application-review\/?$/.test(pathname);
  const waitlistActive = /^\/events\/[^/]+\/dashboard\/waitlist\/?$/.test(
    pathname,
  );
  const applicationStatisticsActive = /^\/application-statistics\/?$/.test(
    pathname,
  );

  const commonNavLinks = (
    <>
      <NavLink
        label="Information"
        href="/information"
        leftSection={<TablerInfoCircle className="w-5 aspect-square" />}
        active={pathname.startsWith("/information")}
      />
      <NavLink
        label="Applications"
        leftSection={<TablerFileText className="w-5 aspect-square" />}
        initialExpanded={
          waitlistActive ||
          applicationReviewActive ||
          applicationStatisticsActive
        }
      >
        <NavLink
          label="Application Review"
          href={`/application-review`}
          leftSection={<TablerClipboardCheck className="w-5 aspect-square" />}
          active={applicationReviewActive}
        />
        <NavLink
          label="Statistics"
          href={`/application-statistics`}
          leftSection={<TablerChartBarPopular className="w-5 aspect-square" />}
          active={applicationStatisticsActive}
        />
      </NavLink>
    </>
  );
  return commonNavLinks;
}
