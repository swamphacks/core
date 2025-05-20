import { createFileRoute } from "@tanstack/react-router";
import { Login } from "@/features/auth/components/Login";

export const Route = createFileRoute("/")({
  component: Index,
});

function Index() {
  return (
    <div className="font-[Figtree] flex items-center justify-center pb-20 bg-background">
      <Login />
    </div>
  );
}
