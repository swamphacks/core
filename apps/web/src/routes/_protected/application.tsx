import { createFileRoute, redirect } from "@tanstack/react-router";
import { ErrorBoundary } from "react-error-boundary";
import { ApplicationForm } from "@/modules/Application/ApplicationForm";
import TablerAlertCircle from "~icons/tabler/alert-circle";
import { useEffect } from "react";
import { hackathonQueryOptions } from "@/modules/Hackathon/hooks/useHackathon";
import { api } from "@/lib/ky";
import { useQueryClient, useSuspenseQuery } from "@tanstack/react-query";
import { useUserQueryKey } from "@/lib/auth/hooks/useUser";
import type { AuthUserResponse } from "@/lib/auth/types";
import { useMyApplication } from "@/modules/Application/hooks/useMyApplication";
import { useApplicationActions } from "@/modules/Application/hooks/useApplicationActions";
import { Button } from "@/components/ui/Button";
import { PageLoading } from "@/components/PageLoading";

export const Route = createFileRoute("/_protected/application")({
  component: RouteComponent,
  beforeLoad: ({ context }) => {
    if (context.user.role === "attendee") {
      throw redirect({
        to: "/hacker-portal",
      });
    }
  },
  pendingComponent: PageLoading,
  loader: ({ context }) => {
    return Promise.all([
      context.queryClient.ensureQueryData(hackathonQueryOptions()),
    ]);
  },
});

function RouteComponent() {
  const { user } = Route.useRouteContext();
  const queryClient = useQueryClient();
  const hackathon = useSuspenseQuery(hackathonQueryOptions());
  const application = useMyApplication();

  const now = new Date();
  const applicationOpen = new Date(hackathon.data.applicationOpen);
  const applicationClose = new Date(hackathon.data.applicationClose);
  const isApplicationOpen = now >= applicationOpen && now <= applicationClose;

  // Show a confirmation dialog when the user closes the tab
  // useEffect(() => {
  //   function beforeUnload(e: BeforeUnloadEvent) {
  //     e.preventDefault();
  //   }

  //   window.addEventListener("beforeunload", beforeUnload);

  //   return () => {
  //     window.removeEventListener("beforeunload", beforeUnload);
  //   };
  // }, []);

  useEffect(() => {
    if (user.hasSeenNewApplicationStatus === false) {
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
  }, [user]);

  if (application.isLoading) {
    return <PageLoading />;
  }

  if (!application.data) {
    return <div>Something went wrong while loading application...</div>;
  }

  const applicationResponses = JSON.parse(atob(application.data.application));
  const name = applicationResponses["firstName"];

  if (application.data.status === "accepted") {
    return <Accepted name={name} />;
  }

  if (application.data.status === "rejected") {
    return <Rejected name={name} />;
  }

  if (application.data.status === "waitlisted") {
    return <Waitlisted name={name} />;
  }

  if (application.data.status === "withdrawn") {
    return <Withdrawn name={name} />;
  }

  if (!isApplicationOpen) {
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
      <ApplicationForm
        hackathon={hackathon.data}
        application={application.data}
        applicationResponses={applicationResponses}
      />
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

interface AcceptedProps {
  name: string;
}

function Accepted({ name }: AcceptedProps) {
  const { confirmAttendance, withdrawApplication } = useApplicationActions();

  const handleConfirmAttendance = async () => {
    await confirmAttendance.mutateAsync();
    window.location.reload();
  };

  const handleWithdrawApplication = async () => {
    const isConfirmed = window.confirm(
      "Are you sure you want to withdraw your application?",
    );

    if (isConfirmed) {
      await withdrawApplication.mutateAsync();
    }
  };

  return (
    <div className="w-full sm:max-w-200 mx-auto font-figtree p-2 relative">
      <h1 className="text-2xl">Congrats, {name}! 🎉</h1>
      <div className="my-3 flex flex-col gap-2">
        <p>You've been accepted to hack in SwampHacks XII!</p>
        <p>
          Please confirm your attendance by January 2nd. Failure to do so means
          you are giving up your spot, and we will admit someone from a
          waitlist.
        </p>
        <p>
          If you're no longer able to attend, please withdraw your application
          so we can offer your spot to another applicant and maintain an
          accurate attendee count.
        </p>
      </div>
      <div className="flex flex-col w-fit items-start gap-2">
        <Button onClick={handleConfirmAttendance}>Confirm Attendance</Button>
        <Button
          onClick={handleWithdrawApplication}
          className="max-w-45 py-2 mt-2"
          variant="secondary"
          size="sm"
        >
          Withdraw Application
        </Button>
      </div>
    </div>
  );
}

interface RejectedProps {
  name: string;
}

function Rejected({ name }: RejectedProps) {
  return (
    <div className="w-full sm:max-w-200 mx-auto font-figtree p-2 relative">
      <h1 className="text-2xl">Hi, {name}!</h1>
      <div className="my-3 flex flex-col gap-3">
        <p>
          We sincerely appreciate your interest in SwampHacks XII and the time
          you took to apply. After careful consideration, we are unable to
          accept you as a hacker at this time.
        </p>

        <p>
          However, we’d love to stay connected and invite you to get involved in
          other ways:
        </p>

        <ol className="flex flex-col gap-2">
          <li>
            1. <strong>Join the Waitlist</strong>: We may have openings
            available closer to the event. You can join the waitlist in the{" "}
            <a className="underline" href="https://app.swamphacks.com/">
              SwampHacks Portal
            </a>{" "}
            or by signing up in person on the day of check-in if space allows.
            The waitlist operates on a first-come, first-served basis.
          </li>
          <li>
            2. <strong>Mentor</strong>: Share your knowledge and guide hackers
            through their projects.{" "}
            <a className="underline" href="https://swamphacks.com/mentor">
              Sign up to be a mentor here
            </a>
            .
          </li>
          <li>
            3. <strong>Volunteer</strong>: Help us run the event smoothly.{" "}
            <a className="underline" href="https://swamphacks.com/volunteer">
              Sign up to volunteer here
            </a>
            .
          </li>
        </ol>

        <p>
          If you have any questions, reach out in our{" "}
          <a href="https://discord.com/invite/NfRPv9JtAG">Discord server</a> or
          email us at{" "}
          <a href="mailto:contact@swamphacks.com">contact@swamphacks.com</a>
        </p>
      </div>
    </div>
  );
}

interface WaitlistedProps {
  name: string;
}

function Waitlisted({ name }: WaitlistedProps) {
  return (
    <div className="w-full sm:max-w-200 mx-auto font-figtree p-2 relative">
      <h1 className="text-2xl">Hi, {name}!</h1>

      <div className="my-3 flex flex-col gap-3">
        <p>
          Thank you for applying to SwampHacks XII! We were very impressed by
          your application. At this time, we’re placing you on our{" "}
          <strong>waitlist</strong> due to limited capacity.
        </p>

        <p>
          If spots open up, we’ll be sending out invitations on a rolling basis
          leading up to the event. Waitlist decisions are made as space becomes
          available. Please keep an eye out on your email or the hacker portal
          for updates!
        </p>

        <p>
          If you have questions, feel free to reach out on our{" "}
          <a href="https://discord.com/invite/NfRPv9JtAG">Discord server</a> or
          email us at{" "}
          <a href="mailto:contact@swamphacks.com">contact@swamphacks.com</a>
        </p>
      </div>
    </div>
  );
}

interface WithdrawnProps {
  name: string;
}

function Withdrawn({ name }: WithdrawnProps) {
  return (
    <div className="w-full sm:max-w-200 mx-auto font-figtree p-2 relative">
      <h1 className="text-2xl">Hi, {name}!</h1>

      <div className="my-3 flex flex-col gap-3">
        <p>Your application for SwampHacks XII has been withdrawn.</p>

        <p>
          We appreciate your interest in SwampHacks and the time you took to
          apply. While we're sorry you won't be able to join us this year, we
          hope to see you at a future event.
        </p>

        <p>
          If your plans change and registration is still open, please reach out
          to our team and we'll do our best to help.
        </p>

        <p>
          You can also stay connected with the SwampHacks community through our{" "}
          <a className="underline" href="https://discord.com/invite/NfRPv9JtAG">
            Discord server
          </a>{" "}
          and follow future announcements for upcoming events and opportunities.
        </p>

        <p>
          If you have any questions, feel free to contact us at{" "}
          <a className="underline" href="mailto:contact@swamphacks.com">
            contact@swamphacks.com
          </a>
          .
        </p>
      </div>
    </div>
  );
}
