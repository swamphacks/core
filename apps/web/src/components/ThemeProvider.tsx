// https://ui.shadcn.com/docs/dark-mode/vite

import { createContext, useContext, useEffect, useRef, useState } from "react";
import { cn } from "@/utils/cn";
import TablerSun from "~icons/tabler/sun";
import TablerMoon from "~icons/tabler/moon";

type Theme = "dark" | "light" | "system";

type ThemeProviderProps = {
  children: React.ReactNode;
  defaultTheme?: Theme;
  storageKey?: string;
};

type ThemeProviderState = {
  theme: Theme;
  setTheme: (theme: Theme) => void;
};

const initialState: ThemeProviderState = {
  theme: "system",
  setTheme: () => null,
};

const ThemeProviderContext = createContext<ThemeProviderState>(initialState);

export function ThemeProvider({
  children,
  defaultTheme = "system",
  storageKey = "ui-theme",
  ...props
}: ThemeProviderProps) {
  const [theme, setTheme] = useState<Theme>(
    () => (localStorage.getItem(storageKey) as Theme) || defaultTheme,
  );
  const isFirstRender = useRef(true);

  useEffect(() => {
    // Skip effect on initial load, since we are already doing the same thing in the <script> tag in index.html
    if (isFirstRender.current) {
      isFirstRender.current = false;
      return;
    }

    const root = window.document.documentElement;

    root.classList.remove("light", "dark");

    if (theme === "system") {
      const systemTheme = window.matchMedia("(prefers-color-scheme: dark)")
        .matches
        ? "dark"
        : "light";

      root.classList.add(systemTheme);
      return;
    }

    root.classList.add(theme);
  }, [theme]);

  const value = {
    theme,
    setTheme: (theme: Theme) => {
      localStorage.setItem(storageKey, theme);
      setTheme(theme);
    },
  };

  return (
    <ThemeProviderContext.Provider {...props} value={value}>
      {children}
    </ThemeProviderContext.Provider>
  );
}

export const ThemeSwitch = () => {
  const { theme, setTheme } = useTheme();

  return (
    <div className="flex gap-2 text-text-main mt-2 border border-input-border rounded-md w-fit">
      <button
        className={cn(
          theme === "light" &&
            "bg-button-secondary border-r-1 border-input-border",
          "cursor-pointer py-1 px-2 rounded-md flex items-center gap-1",
        )}
        onClick={() => setTheme("light")}
      >
        Light
        <TablerSun />
      </button>
      <button
        className={cn(
          theme === "dark" &&
            "bg-button-secondary border-l-1 border-input-border",
          "cursor-pointer py-1 px-2 rounded-md flex items-center gap-1",
        )}
        onClick={() => setTheme("dark")}
      >
        Dark
        <TablerMoon />
      </button>
    </div>
  );
};

export const useTheme = (): ThemeProviderState => {
  const context = useContext(ThemeProviderContext);

  if (context === undefined)
    throw new Error("useTheme must be used within a ThemeProvider");

  let theme = context.theme;

  if (theme === "system") {
    theme = window.matchMedia("(prefers-color-scheme: dark)").matches
      ? "dark"
      : "light";
  }

  return {
    ...context,
    theme,
  };
};
