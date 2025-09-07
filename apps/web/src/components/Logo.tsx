import { useTheme } from "./ThemeProvider";

export function Logo() {
  const { theme } = useTheme();

  return (
    <img
      className="py-5"
      src={
        theme === "dark"
          ? "/assets/SwampHacks_Logo_Light.png"
          : "/assets/SwampHacks_Logo_Dark.png"
      }
      alt="SwampHacks Logo"
    />
  );
}
