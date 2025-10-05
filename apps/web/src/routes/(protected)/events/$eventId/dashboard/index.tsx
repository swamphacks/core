import { PageUnderConstruction } from "@/components/PageUnderConstruction";
import { createFileRoute, redirect } from "@tanstack/react-router";

export const Route = createFileRoute("/(protected)/events/$eventId/dashboard/")(
  {
    component: RouteComponent,
    beforeLoad: ({ context, params }) => {
      if (context.eventRole === "applicant") {
        throw redirect({
          to: `/events/$eventId/dashboard/application-status`,
          params: {
            eventId: params.eventId,
          },
        });
      }
    },
  },
);

function RouteComponent() {
  return <PageUnderConstruction />;
}
