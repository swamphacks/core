import { createFileRoute } from "@tanstack/react-router";
import { ErrorBoundary } from "react-error-boundary";
import { ApplicationForm } from "@/modules/Application/ApplicationForm";
import TablerAlertCircle from "~icons/tabler/alert-circle";
import { useEffect } from "react";
import { useHackathon } from "@/modules/Hackathon/hooks/useHackathon";
import { Spinner } from "@/components/ui/Spinner";
import { api } from "@/lib/ky";
import { auth } from "@/lib/authClient";
import { useQueryClient } from "@tanstack/react-query";
import { useUserQueryKey } from "@/lib/auth/hooks/useUser";
import type { AuthUserResponse } from "@/lib/auth/types";

export const Route = createFileRoute("/_protected/_user/application")({
  component: RouteComponent,
});

function RouteComponent() {
  const hackathon = useHackathon();
  const { data } = auth.useUser();
  const { user } = data!;
  const queryClient = useQueryClient();

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

  useEffect(() => {
    if (
      user?.hasSeenNewApplicationStatus === false &&
      hackathon.data &&
      !hackathon.isLoading
    ) {
      api.post(`users/me/acknowledge-new-application-status`);

      queryClient.setQueryData(
        useUserQueryKey,
        (oldData: AuthUserResponse) => ({
          ...oldData,
          user: {
            ...oldData.user,
            hasSeenNewApplicationStatus: true,
          },
        }),
      );
    }
  }, [user, hackathon]);

  if (hackathon.isLoading) {
    return (
      <div className="flex w-full justify-center pt-10 gap-2 text-text-secondary">
        <Spinner />
        <p>Loading hackathon...</p>
      </div>
    );
  }

  if (!hackathon.data) {
    return (
      <div className="h-full flex justify-center items-center gap-2 text-red-400">
        <TablerAlertCircle />
        <p>Something went wrong while loading hackathon information :(</p>
      </div>
    );
  }

  if (new Date() > new Date(hackathon.data.applicationClose)) {
    return (
      <div className="max-w-xs mx-auto h-full flex flex-col justify-center items-center gap-8 text-text-secondary">
        <div className="flex flex-row items-center justify-center gap-2">
          <TablerAlertCircle />
          <p>Applications have closed!</p>
        </div>
      </div>
    );
  }

  return (
    <ErrorBoundary FallbackComponent={Fallback}>
      <ApplicationForm hackathon={hackathon.data} />
    </ErrorBoundary>
  );
}

function Fallback() {
  return (
    <div className="h-full flex justify-center items-center gap-2 text-red-400">
      <TablerAlertCircle />
      <p>Something went wrong while loading application form :(</p>
    </div>
  );
}
