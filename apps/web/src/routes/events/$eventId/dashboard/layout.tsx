import { Outlet, useParams, createFileRoute } from "@tanstack/react-router";
import { useQuery } from "@tanstack/react-query";
import StaffAppShell from "./components/StaffAppShell";
import AttendeeAppShell from "./components/AttendeeAppShell";
import ApplicantAppShell from "./components/ApplicantAppShell";
import { api } from "@/lib/ky";
export const Route = createFileRoute("/events/$eventId/dashboard")({
  component: RouteComponent,
});
function useUserRole() {
  const { eventId } = useParams({ from: "/events/$eventId/dashboard" });
  const { data, isLoading, error } = useQuery({
    queryKey: ["event-role", eventId],
    queryFn: async () => {
      const res = await api.get(`events/${eventId}/role`).json<any>();
      return res.role;
    },
    enabled: !!eventId,
  });

  if (isLoading) return undefined;
  if (error) return undefined;
  return data as "staff" | "attendee" | "applicant" | undefined;
}
function RouteComponent() {
  const role = useUserRole();
  if (!role) return <div>Loading...</div>;
  if (role === "staff") {
    return (
      <StaffAppShell>
        <Outlet />
      </StaffAppShell>
    );
  }
  if (role === "attendee") {
    return (
      <AttendeeAppShell>
        <Outlet />
      </AttendeeAppShell>
    );
  }
  if (role === "applicant") {
    return (
      <ApplicantAppShell>
        <Outlet />
      </ApplicantAppShell>
    );
  }
  return <div>Unknown role, Possibly Admin</div>;
}
