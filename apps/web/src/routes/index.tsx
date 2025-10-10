import { createFileRoute, redirect } from "@tanstack/react-router";
import { Login } from "@/features/Auth/components/Login";
import { PageLoading } from "@/components/PageLoading";
import { z } from "zod";

export const Route = createFileRoute("/")({
  validateSearch: z.object({
    redirect: z.string().optional().catch(""),
  }),
  beforeLoad: async ({ context }) => {
    const { user } = await context.userQuery.promise;

    console.log("Loaded user in beforeLoad:", user);

    if (user) {
      console.log("User is already authenticated, redirecting to portal.");
      throw redirect({
        to: "/portal",
      });
    }
  },
  pendingMs: 50,
  // TODO: Use skeleton components to display a loading state
  pendingComponent: () => PageLoading(),
  component: Index,
});

function Index() {
  return (
    <div className="h-full flex items-center justify-center pb-15">
      <Login />
    </div>
  );
}
