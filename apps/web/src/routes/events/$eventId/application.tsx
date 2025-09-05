import { createFileRoute } from "@tanstack/react-router";
import { ErrorBoundary } from "react-error-boundary";
import { ApplicationForm } from "@/features/Application/components/ApplicationForm";
import TablerAlertCircle from "~icons/tabler/alert-circle";

export const Route = createFileRoute("/events/$eventId/application")({
  component: RouteComponent,
});

function RouteComponent() {
  const { eventId } = Route.useParams();
  return (
    <ErrorBoundary FallbackComponent={Fallback}>
      <div className="w-full h-screen bg-surface">
        {/* this padding left prevent the page being shifted when the form fully loads because of the scrollbar */}
        <div className="w-full bg-surface transition-[background] sm:pl-[calc(100vw-100%)]">
          <ApplicationForm eventId={eventId} />
        </div>
      </div>
    </ErrorBoundary>
  );
}

function Fallback() {
  return (
    <div className="w-full h-full bg-surface flex justify-center items-center gap-2 text-red-400">
      <TablerAlertCircle />
      <p>Something went wrong while loading form :(</p>
    </div>
  );
}
