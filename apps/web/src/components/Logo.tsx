import { cn } from "@/utils/cn";
import { useTheme } from "./ThemeProvider";

export function Logo({ className }: { className?: string }) {
  const { theme } = useTheme();

  return (
    <img
      className={cn(className)}
      src={
        theme === "dark"
          ? "/assets/SwampHacks_Logo_Light.png"
          : "/assets/SwampHacks_Logo_Dark.png"
      }
      alt="SwampHacks Logo"
    />
  );
}
