import { createFileRoute } from "@tanstack/react-router";
import { Login } from "@/features/common/Auth/components/Login";

export const Route = createFileRoute("/")({
  component: Index,
});

function Index() {
  return (
    <div className="h-full flex items-center justify-center pb-20">
      <Login />
    </div>
  );
}
