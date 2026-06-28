import { createFileRoute } from "@tanstack/react-router";
import Cookies from "js-cookie";
import { PageLoading } from "@/components/PageLoading";
import { hackathonQueryOptions } from "@/modules/Hackathon/hooks/useHackathon";
import { useSuspenseQuery } from "@tanstack/react-query";
import InformationPage from "@/modules/Information/InformationPage";

export const Route = createFileRoute("/_protected/information")({
  beforeLoad: (context) => {
    const { user } = context.context;
    const hasSkippedCookie = Cookies.get("welcome-modal-skipped") === "true";
    const showOnboardingModal = !hasSkippedCookie && !!user && !user.onboarded;
    return {
      showOnboardingModal: showOnboardingModal,
    };
  },
  pendingComponent: PageLoading,
  component: RouteComponent,
  loader: ({ context }) => {
    return Promise.all([
      context.queryClient.ensureQueryData(hackathonQueryOptions()),
    ]);
  },
});

function RouteComponent() {
  const { user } = Route.useRouteContext();
  // const [isModalOpen, setIsModalOpen] = useState(showOnboardingModal);
  const hackathon = useSuspenseQuery(hackathonQueryOptions());

  return <InformationPage user={user} hackathon={hackathon.data} />;
}
