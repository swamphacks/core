import { cn } from "@/utils/cn";
import { useTheme } from "./ThemeProvider";

interface LogoProps {
  className?: string;
  onClick?: () => void;
  logo?: string;
  label?: string;
  hideOnMobile?: boolean;
  hideLabelOnMobile?: boolean;
  hideLogoOnMobile?: boolean;
}

export function Logo({
  className,
  onClick,
  label,
  logo,
  hideOnMobile,
  hideLabelOnMobile,
  hideLogoOnMobile,
}: LogoProps) {
  const { theme } = useTheme();

  return (
    <div
      onClick={onClick}
      className={cn(
        "flex items-center gap-2",
        className,
        hideOnMobile && "hidden md:flex",
      )}
    >
      <div className={cn("w-13", hideLogoOnMobile && "hidden md:flex")}>
        <img
          src={
            logo
              ? logo
              : theme === "dark" // Default logos if no logo prop is provided
                ? "/assets/SwampHacks_Logo_Light.png"
                : "/assets/SwampHacks_Logo_Dark.png"
          }
          alt="SwampHacks Logo"
        />
      </div>

      {label && (
        <h1
          className={cn(
            "text-xl font-bold",
            hideLabelOnMobile && "hidden md:flex",
          )}
        >
          {label}
        </h1>
      )}
    </div>
  );
}
