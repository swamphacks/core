import { createContext, useContext } from "react";

export interface AppShellContextValue {
  isMobileNavOpen: boolean;
  toggleMobileNav: () => void;
  setMobileNavOpen: (open: boolean) => void;
}

export const AppShellContext = createContext<AppShellContextValue | undefined>(
  undefined,
);

export const useAppShell = () => {
  const ctx = useContext(AppShellContext);
  if (!ctx) throw new Error("useAppShell must be used within <AppShell>");
  return ctx;
};
