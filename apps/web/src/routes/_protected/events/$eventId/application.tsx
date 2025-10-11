import { createFileRoute } from "@tanstack/react-router";
import { ErrorBoundary } from "react-error-boundary";
import { ApplicationForm } from "@/features/Application/components/ApplicationForm";
import TablerAlertCircle from "~icons/tabler/alert-circle";
import { useEffect } from "react";

export const Route = createFileRoute("/_protected/events/$eventId/application")(
  {
    component: RouteComponent,
  },
);

function RouteComponent() {
  const { eventId } = Route.useParams();

  // Show a confirmation dialog when the user closes the tab
  useEffect(() => {
    function beforeUnload(e: BeforeUnloadEvent) {
      e.preventDefault();
    }

    window.addEventListener("beforeunload", beforeUnload);

    return () => {
      window.removeEventListener("beforeunload", beforeUnload);
    };
  }, []);

  return (
    <ErrorBoundary FallbackComponent={Fallback}>
      <div className="w-full h-screen bg-white dark:bg-background">
        {/* this padding left prevent the page being shifted when the form fully loads because of the scrollbar */}
        <div className="w-full bg-white dark:bg-background transition-[background] sm:pl-[calc(100vw-100%)]">
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
