import { PageUnderConstruction } from "@/components/PageUnderConstruction";
import { createFileRoute, notFound } from "@tanstack/react-router";

export const Route = createFileRoute(
  "/(protected)/events/$eventId/dashboard/(applicant)/application-status",
)({
  beforeLoad: async ({ context }) => {
    if (
      !context.eventRole ||
      !["applicant", "admin"].includes(context.eventRole)
    ) {
      return notFound();
    }
    return {};
  },
  component: RouteComponent,
});

function RouteComponent() {
  return <PageUnderConstruction />;
}
