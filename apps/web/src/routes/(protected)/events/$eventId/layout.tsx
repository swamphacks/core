import { getUserEventRole } from "@/features/Event/api/getUserEventRole";
import { createFileRoute, Outlet } from "@tanstack/react-router";

export const Route = createFileRoute("/(protected)/events/$eventId")({
  component: RouteComponent,
  beforeLoad: async ({ context, params }) => {
    const data = await context.queryClient.fetchQuery({
      queryKey: ["eventRole", params.eventId],
      queryFn: () => getUserEventRole(params.eventId),
    });

    // Parse assigned at string to Date object
    let parsedAssignedAt: Date | null = null;
    if (data?.assigned_at) {
      parsedAssignedAt = new Date(data.assigned_at);
    }

    return {
      eventRole: data?.role,
      roleAssignedAt: parsedAssignedAt ?? null,
    };
  },
});

function RouteComponent() {
  return <Outlet />;
}
