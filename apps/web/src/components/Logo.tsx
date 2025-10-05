import { cn } from "@/utils/cn";
import { useTheme } from "./ThemeProvider";

interface LogoProps {
  className?: string;
  onClick?: () => void;
}

export function Logo({ className, onClick }: LogoProps) {
  const { theme } = useTheme();

  return (
    <div
      onClick={onClick}
      className={cn("flex items-center gap-2 ml-2", className)}
    >
      <div className="w-13">
        <img
          src={
            theme === "dark"
              ? "/assets/SwampHacks_Logo_Light.png"
              : "/assets/SwampHacks_Logo_Dark.png"
          }
          alt="SwampHacks Logo"
        />
      </div>

      <h1 className="text-xl font-bold">SwampHacks</h1>
    </div>
  );
}
