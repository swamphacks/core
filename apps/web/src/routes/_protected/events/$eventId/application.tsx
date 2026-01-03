import { createFileRoute, useRouter } from "@tanstack/react-router";
import { ErrorBoundary } from "react-error-boundary";
import { ApplicationForm } from "@/features/Application/components/ApplicationForm";
import TablerAlertCircle from "~icons/tabler/alert-circle";
import { useEffect } from "react";
import { useEvent } from "@/features/Event/hooks/useEvent";
import { Button } from "@/components/ui/Button";

export const Route = createFileRoute("/_protected/events/$eventId/application")(
  {
    component: RouteComponent,
  },
);

function RouteComponent() {
  const router = useRouter();
  const { eventId } = Route.useParams();
  const event = useEvent(eventId);

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

  if (event.isLoading) {
    return null;
  }

  if (event.isError || event.error || !event.data) {
    return (
      <div className="w-full h-full bg-surface flex justify-center items-center gap-2 text-red-400">
        <TablerAlertCircle />
        <p>Something went wrong while loading event info :(</p>
      </div>
    );
  }

  if (new Date() > event.data.application_close) {
    return (
      <div className="max-w-xs mx-auto h-full flex flex-col justify-center items-center gap-8 text-text-secondary">
        <div className="flex flex-row items-center justify-center gap-2">
          <TablerAlertCircle />
          <p>Applications have closed for this event...</p>
        </div>

        <Button
          variant="secondary"
          className="w-full"
          onPress={() => router.navigate({ to: "/portal" })}
        >
          Go Back
        </Button>
      </div>
    );
  }

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
