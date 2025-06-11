import { createFileRoute, redirect } from "@tanstack/react-router";
import { Login } from "@/features/Auth/components/Login";
import { z } from "zod";
import { authClient } from "@/lib/authClient";

export const Route = createFileRoute("/")({
  validateSearch: z.object({
    redirectTo: z.string().optional().catch(""),
  }),
  beforeLoad: async () => {
    const { user, error } = await authClient.getUser();
    if (error) throw error;

    console.log("Loaded user in beforeLoad:", user);

    if (user) {
      console.log("User is already authenticated, redirecting to dashboard.");
      throw redirect({
        to: "/dashboard",
      });
    }

    return { user };
  },
  component: Index,
});

function Index() {
  return (
    <div className="h-full flex items-center justify-center pb-15">
      <Login />
    </div>
  );
}
