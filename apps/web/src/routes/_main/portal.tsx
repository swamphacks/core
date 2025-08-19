import { Button } from "@/components/ui/Button";
import { EventCard } from "@/features/Event/components/EventCard";
import { auth } from "@/lib/authClient";
import { createFileRoute } from "@tanstack/react-router";
import { Heading } from "react-aria-components";

export const Route = createFileRoute("/_main/portal")({
  component: RouteComponent,
});

function RouteComponent() {
  const { user } = Route.useRouteContext();
  const logout = async () => {
    try {
      await auth.logOut();
      window.location.href = "/";
    } catch (error) {
      console.error("Error during logout:", error);
    }
  };

  return (
    <div className="flex flex-col gap-8">
      <div className="flex flex-col gap-2">
        <Heading className="text-3xl text-text-main">
          Welcome, {user?.name ?? "hacker"}!
        </Heading>
        <h2 className="text-text-secondary text-xl">Ready to start hacking?</h2>
      </div>

      <div className="flex flex-row gap-6">
        <EventCard
          eventId="sfjslfjds"
          status="notApplied"
          title="SwampHacks XI"
          description="Swamphacks's 11th hackathon"
          date="Jan 20-21st"
          location="Newell Hall"
        />

        <EventCard
          eventId="sfjslfjds"
          status="notApplied"
          title="SwampHacks XI"
          description="Swamphacks's 11th hackathon"
          date="Jan 20-21st"
          location="Newell Hall"
        />

        <EventCard
          eventId="sfjslfjds"
          status="notApplied"
          title="SwampHacks XI"
          description="Swamphacks's 11th hackathon"
          date="Jan 20-21st"
          location="Newell Hall"
        />
      </div>
    </div>
  );
}
