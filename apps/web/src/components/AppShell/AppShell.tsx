import {
  Children,
  isValidElement,
  memo,
  useMemo,
  type FC,
  type PropsWithChildren,
  type ReactNode,
} from "react";
import { SlideoutNavbar } from "./MobileSidebar";
import { AppShellContext } from "./AppShellContext";
import { useToggleState } from "react-stately";
import { Button as RACButton } from "react-aria-components";
import TablerMenu2 from "~icons/tabler/menu-2";
import IconX from "~icons/tabler/x";
import { auth } from "@/lib/authClient";
import { Profile } from "./Profile";
import { MobileProfile } from "@/components/AppShell/MobileProfile";

interface AppShellComponent extends FC<PropsWithChildren> {
  Header: FC<PropsWithChildren>;
  Navbar: FC<PropsWithChildren>;
  Main: FC<PropsWithChildren>;
}

function extractAppShellChildren(children: ReactNode) {
  let header: ReactNode = null;
  let navbar: ReactNode = null;
  let main: ReactNode = null;

  Children.forEach(children, (child) => {
    if (!isValidElement(child)) return;

    if (child.type === AppShell.Header) header = child;
    else if (child.type === AppShell.Navbar) navbar = child;
    else if (child.type === AppShell.Main) main = child;
  });

  return { header, navbar, main };
}

const AppShellBase: FC<PropsWithChildren> = ({ children }) => {
  const { toggle, isSelected, setSelected } = useToggleState();

  const { header, navbar, main } = useMemo(
    () => extractAppShellChildren(children),
    [children],
  );
  const { data } = auth.useUser();
  const user = data?.user;

  if (data?.error || !user) {
    return <p>Something went wrong while loading user information.</p>;
  }

  return (
    <AppShellContext.Provider
      value={{
        isMobileNavOpen: isSelected,
        toggleMobileNav: toggle,
        setMobileNavOpen: setSelected,
      }}
    >
      <div className="flex h-screen w-screen flex-col relative overflow-hidden">
        {/* Topbar */}
        <header className="bg-surface md:hidden h-16 w-full flex items-center px-4 md:px-2 border-b-1 border-neutral-300 dark:border-neutral-800">
          {/* Mobile Burger menu */}
          <RACButton onPress={toggle}>
            {isSelected ? (
              <IconX className="w-8 h-8" />
            ) : (
              <TablerMenu2 className="w-8 h-8" />
            )}
          </RACButton>
          <div className="flex justify-between w-full items-center">
            {header}
            <MobileProfile
              name={user.name}
              role={user.role}
              logout={async () => {
                await auth.logOut();
                window.location.href = "/";
              }}
            />
          </div>
        </header>

        <div className="flex flex-1 min-h-0 relative">
          {/* Desktop Sidebar */}
          {navbar && (
            <aside className="w-64 h-full px-2 py-3 border-r bg-surface border-neutral-200 dark:border-neutral-800 hidden md:block">
              <nav className="flex flex-col gap-2 h-full">
                {header}
                <div className="flex flex-col justify-between h-full">
                  <div>{navbar}</div>
                  <div className="flex flex-col gap-5">
                    <Profile name={user.name} role={user.role} />
                  </div>
                </div>
              </nav>
            </aside>
          )}

          {/* Slideout for Mobile */}
          <SlideoutNavbar isOpen={isSelected}>
            <div className="flex flex-col justify-between h-full">
              <div>{navbar}</div>
            </div>
          </SlideoutNavbar>

          {/* Main content */}
          {main && (
            <main className="flex-1 p-6 overflow-y-auto relative bg-background">
              {main}
            </main>
          )}
        </div>
      </div>
    </AppShellContext.Provider>
  );
};

const AppShell = memo(AppShellBase) as unknown as AppShellComponent;

AppShell.Header = memo(({ children }: PropsWithChildren) => <>{children}</>);
AppShell.Navbar = memo(({ children }: PropsWithChildren) => <>{children}</>);
AppShell.Main = memo(({ children }: PropsWithChildren) => <>{children}</>);

AppShell.displayName = "AppShell";
AppShell.Header.displayName = "AppShell.Header";
AppShell.Navbar.displayName = "AppShell.Navbar";
AppShell.Main.displayName = "AppShell.Main";

export { AppShell };
