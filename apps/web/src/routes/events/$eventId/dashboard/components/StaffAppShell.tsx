import { Link as TanstackLink } from "@tanstack/react-router";
import { AppShell } from "@/components/AppShell/AppShell";
import { NavLink } from "@/components/AppShell/NavLink";
import TablerLayoutCollage from "~icons/tabler/layout-collage";
import TablerSocial from "~icons/tabler/social";
import React from "react"; // Import React to use ReactNode type

// Accept a 'children' prop
export default function StaffAppShell({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <AppShell>
      <AppShell.Header>
        <div className="w-full px-4 flex flex-row justify-between h-full items-center">
          <h1 className="text-2xl font-bold">Staff Dashboard</h1>
        </div>
      </AppShell.Header>

      <AppShell.Navbar>
        <TanstackLink to="/portal">
          {({ isActive }) => (
            <NavLink
              label="Emails"
              leftSection={
                <TablerLayoutCollage className="w-5 aspect-square" />
              }
              active={isActive}
            />
          )}
        </TanstackLink>

        <TanstackLink to="/community">
          {({ isActive }) => (
            <NavLink
              label="Attendees"
              leftSection={<TablerSocial className="w-5 aspect-square" />}
              active={isActive}
            />
          )}
        </TanstackLink>
      </AppShell.Navbar>

      <AppShell.Main>
        {/* Render the children prop here instead of a hardcoded Outlet */}
        {children}
      </AppShell.Main>
    </AppShell>
  );
}
