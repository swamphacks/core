import { NavLink } from "@/components/AppShell/NavLink";
import TablerInfoCircle from "~icons/tabler/info-circle";
import TablerClipboardCheck from "~icons/tabler/clipboard-check";
import TablerFileText from "~icons/tabler/file-text";
import TablerChartBarPopular from "~icons/tabler/chart-bar-popular";
import TablerSearch from "~icons/tabler/search";

interface AdminNavbaPropsProps {
  pathname: string;
}

export default function AdminNavbaProps({ pathname }: AdminNavbaPropsProps) {
  const applicationReviewActive = /^\/application-review\/?$/.test(pathname);
  const applicationStatisticsActive = /^\/application-statistics\/?$/.test(
    pathname,
  );
  const applicationSearchActive = /^\/application-search\/?$/.test(pathname);

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
          applicationReviewActive ||
          applicationStatisticsActive ||
          applicationSearchActive
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
        <NavLink
          label="Search"
          href={`/application-search`}
          leftSection={<TablerSearch className="w-5 aspect-square" />}
          active={applicationSearchActive}
        />
      </NavLink>
    </>
  );
  return commonNavLinks;
}
