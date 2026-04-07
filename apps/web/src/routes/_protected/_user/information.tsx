import { createFileRoute } from "@tanstack/react-router";
import { useState } from "react";
import { OnboardingModal } from "@/modules/Onboarding/OnboardingModal";
import Cookies from "js-cookie";

export const Route = createFileRoute("/_protected/_user/information")({
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
  const { showOnboardingModal } = Route.useRouteContext();
  const [isModalOpen, setIsModalOpen] = useState(showOnboardingModal);

  return (
    <div>
      <p>TODO</p>
      <OnboardingModal
        isOpen={isModalOpen}
        onOpenChange={setIsModalOpen}
      ></OnboardingModal>
    </div>
  );
}
