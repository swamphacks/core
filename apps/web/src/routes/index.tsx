import { z } from "zod";
import { createFileRoute } from "@tanstack/react-router";
import { Login } from "@/features/Auth/components/Login";

export const Route = createFileRoute("/")({
  validateSearch: z.object({
    redirectTo: z.string().optional().catch(""),
  }),
  component: Index,
});

function Index() {
  return (
    <div className="h-full flex items-center justify-center pb-15">
      <Login />
    </div>
  );
}
