import { DiscordIcon } from "@/components/icons/Discord";
import { Button } from "@/components/ui/Button";
import { useTheme } from "@/components/ThemeProvider";

// import Auth from "@/lib/auth";
// import { Discord } from "@/lib/auth/providers";

// const authClient = Auth({
//   providers: [Discord],
//   redirect_uri: "test",
// });

const Logo = ({ src }: { src: string }) => {
  return <img className="py-5" src={src} alt="SwampHacks Logo" />;
};

const Login = () => {
  const { theme } = useTheme();

  return (
    <div className="flex flex-col items-center bg-surface rounded-md px-3 pt-6 shadow-md text-text-main">
      <Logo
        src={
          theme === "dark"
            ? "/assets/SwampHacks_Logo_Light.png"
            : "/assets/SwampHacks_Logo_Dark.png"
        }
      />
      <h1 className="text-center font-bold text-2xl">SwampHacks Portal</h1>
      <p className="text-sm text-text-secondary mt-2">
        Log in to access the SwampHacks Portal.
      </p>
      <Button
        className="items-center gap-2 w-[80%] my-4"
        onClick={() => {
          // authClient.signIn("discord", "previous url");
          console.log("hi");
        }}
      >
        <span>
          <DiscordIcon />
        </span>
        Log in with Discord
      </Button>
      <p className="w-[75%] text-xs text-center mb-5 text-text-secondary">
        By logging in, you agree to our Terms of Service and Privacy Policy.
      </p>
    </div>
  );
};

export { Login };
