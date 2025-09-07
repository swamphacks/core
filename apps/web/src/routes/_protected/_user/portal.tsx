import { EventCard } from "@/features/Event/components/EventCard";
import { useEventsWithUserInfo } from "@/features/Event/hooks/useEventsWithUserInfo";
import { createFileRoute } from "@tanstack/react-router";
import { Heading, Text } from "react-aria-components";
import { useState } from "react";
import { OnboardingModal } from "@/features/Onboarding/components/OnboardingModal";
import Cookies from "js-cookie";

export const Route = createFileRoute("/_protected/_user/portal")({
  beforeLoad: (context) => {
    const { user } = context.context;
    const hasSkippedCookie = Cookies.get("welcome-modal-skipped") === "true";
    const showOnboardingModal = !hasSkippedCookie && !!user && !user.onboarded;
    return {
      showOnboardingModal: showOnboardingModal,
    };
  },
  component: RouteComponent,
});

function RouteComponent() {
  const { user, showOnboardingModal } = Route.useRouteContext();
  const { data, isLoading, isError } = useEventsWithUserInfo();
  const [isModalOpen, setIsModalOpen] = useState(showOnboardingModal);

  if (isLoading) {
    return (
      <div className="flex flex-col gap-8">
        <div className="flex flex-col gap-2">
          <Heading className="text-3xl text-text-main">
            Welcome, {user?.name ?? "hacker"}!
          </Heading>
          <h2 className="text-text-secondary text-xl">
            Ready to start hacking?
          </h2>
        </div>

        <div className="flex flex-row gap-6">
          {[...Array(3)].map((_, i) => (
            <div
              key={i}
              className="w-96 h-64 bg-neutral-200 dark:bg-neutral-800 animate-pulse rounded"
            />
          ))}
        </div>
      </div>
    );
  }

  if (isError || data === undefined) {
    return (
      <div className="flex flex-col gap-4">
        <div className="flex flex-col gap-2">
          <Heading className="text-3xl text-text-main">
            Welcome, {user?.name ?? "hacker"}!
          </Heading>
          <h2 className="text-text-secondary text-xl">
            Ready to start hacking?
          </h2>
        </div>

        <div className="flex flex-row gap-6">
          <Text className="text-red-500">
            Whoops, something went wrong, please refresh and try again!
          </Text>
        </div>
      </div>
    );
  }

  return (
    <div>
      <div className="flex flex-col gap-8">
        <div className="flex flex-col gap-2">
          <Heading className="text-3xl text-text-main">
            Welcome, {user?.name ?? "hacker"}!
          </Heading>
          <h2 className="text-text-secondary text-xl">
            Ready to start hacking?
          </h2>
        </div>

        <div className="flex flex-row gap-6">
          {data.map((event) => (
            <EventCard key={event.eventId} {...event} />
          ))}

          {data.length === 0 && (
            <Text className="text-violet-600">
              Awww such empty. Please check back later for events!
            </Text>
          )}
        </div>
      </div>
      <OnboardingModal
        isOpen={isModalOpen}
        onOpenChange={setIsModalOpen}
      ></OnboardingModal>
    </div>
  );
}
