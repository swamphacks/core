import ApplicationReviewPage from "@/features/Application/components/ApplicationReview/ApplicationReviewPage";
import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute(
  "/_protected/events/$eventId/dashboard/_staff/application-review",
)({
  component: RouteComponent,
});

function RouteComponent() {
  const { eventId } = Route.useParams();

  return <ApplicationReviewPage eventId={eventId} />;
}
