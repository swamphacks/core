import { createFileRoute } from "@tanstack/react-router";
import { ApplicationForm } from "@/features/Application/components/ApplicationForm";

export const Route = createFileRoute("/application")({
  component: RouteComponent,
});

function RouteComponent() {
  return (
    <div className="w-full h-screen bg-surface">
      <div className="w-full bg-surface transition-[background]">
        <ApplicationForm
          title="SwampHacks XI Application"
          description="SwampHacks is the University of Florida’s largest annual hackathon. A 36-hour tech event where students from across the country come together to build projects, learn new skills, and connect with fellow innovators."
        />
      </div>
    </div>
  );
}
