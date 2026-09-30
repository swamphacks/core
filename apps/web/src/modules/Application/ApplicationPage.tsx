import { PageLoading } from "@/components/PageLoading";
import { Button } from "@/components/ui/Button";
import { useUserQueryKey } from "@/lib/auth/hooks/useUser";
import type { AuthUserResponse, UserContext } from "@/lib/auth/types";
import { api } from "@/lib/ky";
import { ApplicationForm } from "@/modules/Application/ApplicationForm";
import { useApplicationActions } from "@/modules/Application/hooks/useApplicationActions";
import { useMyApplication } from "@/modules/Application/hooks/useMyApplication";
import type { Hackathon } from "@/modules/Hackathon/hooks/useHackathon";
import { useQueryClient } from "@tanstack/react-query";
import { useEffect } from "react";
import { ErrorBoundary } from "react-error-boundary";
import TablerAlertCircle from "~icons/tabler/alert-circle";

interface ApplicationPageProps {
  hackathon: Hackathon;
  user: UserContext;
}

export default function ApplicationPage({
  hackathon,
  user,
}: ApplicationPageProps) {
  const queryClient = useQueryClient();
  const application = useMyApplication();

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
  }, [user, queryClient]);

  if (application.isLoading) {
    return <PageLoading />;
  }

  if (!application.data) {
    return <div>Something went wrong while loading application...</div>;
  }

  const applicationResponses = JSON.parse(atob(application.data.application));
  const name = applicationResponses["firstName"];

  if (application.data.status === "accepted") {
    return (
      <Accepted
        name={name}
        rspvDeadline={application.data.rsvpDeadline ?? null}
      />
    );
  }

  if (application.data.status === "confirmed") {
    return <Confirmed name={name} />;
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

  const now = new Date();
  const applicationOpen = new Date(hackathon.applicationOpen);
  const applicationClose = new Date(hackathon.applicationClose);

  let isApplicationOpen;
  if (hackathon.acceptEarlyApplications) {
    const earlyApplicationOpen = new Date(hackathon.earlyApplicationOpen!);
    const earlyApplicationClose = new Date(hackathon.earlyApplicationClose!);
    isApplicationOpen =
      (now >= earlyApplicationOpen && now <= earlyApplicationClose) ||
      (now >= applicationOpen && now <= applicationClose);
  } else {
    isApplicationOpen = now >= applicationOpen && now <= applicationClose;
  }

  if (!isApplicationOpen) {
    return (
      <div className="max-w-xs mx-auto h-full flex flex-col justify-center items-center gap-8 text-text-secondary">
        <div className="flex flex-row items-center justify-center gap-2">
          <TablerAlertCircle />
          <p>Applications are currently closed.</p>
        </div>
      </div>
    );
  }

  return (
    <ErrorBoundary FallbackComponent={Fallback}>
      <ApplicationForm
        hackathon={hackathon}
        application={application.data}
        applicationResponses={applicationResponses}
        user={user}
      />
    </ErrorBoundary>
  );
}

interface AcceptedProps {
  name: string;
  rspvDeadline: string | null;
}

function Accepted({ name, rspvDeadline }: AcceptedProps) {
  const { confirmAttendance, withdrawApplication } = useApplicationActions();

  const handleConfirmAttendance = async () => {
    confirmAttendance.mutate();
  };

  const handleWithdrawApplication = async () => {
    const isConfirmed = window.confirm(
      "Withdraw your attendance? This releases your spot at SwampHacks XII.",
    );

    if (isConfirmed) {
      withdrawApplication.mutate();
    }
  };

  return (
    <div className="w-full sm:max-w-200 mx-auto font-figtree p-2 relative">
      <h1 className="text-2xl">Congrats, {name}! 🎉</h1>
      <div className="my-3 flex flex-col gap-2">
        <p>You've been accepted to hack in SwampHacks XII!</p>
        <p>
          {rspvDeadline ? (
            <>
              Please confirm your attendance by{" "}
              {new Date(rspvDeadline).toLocaleString("en-US", {
                timeZone: "America/New_York",
                month: "long",
                day: "numeric",
                year: "numeric",
                hour: "numeric",
                minute: "2-digit",
                hour12: true,
              }) + " ET"}
              . Failure to do so means you are giving up your spot, and we will
              admit someone from the waitlist.
            </>
          ) : (
            "Please confirm your attendance to reserve your spot."
          )}
        </p>
        <p>
          If you're no longer able to attend, please withdraw your application
          so we can offer your spot to another applicant and maintain an
          accurate attendee count.
        </p>
      </div>
      <div className="flex flex-col w-fit items-start gap-2">
        <Button
          onClick={handleConfirmAttendance}
          isDisabled={
            confirmAttendance.isPending || withdrawApplication.isPending
          }
          size="md"
          className="w-[200px] max-w-full min-h-10"
        >
          {confirmAttendance.isPending ? "Confirming..." : "Confirm Attendance"}
        </Button>
        <Button
          onClick={handleWithdrawApplication}
          isDisabled={
            confirmAttendance.isPending || withdrawApplication.isPending
          }
          size="md"
          className="w-[200px] max-w-full min-h-10"
          variant="danger"
        >
          {withdrawApplication.isPending
            ? "Withdrawing..."
            : "Withdraw Attendance"}
        </Button>
      </div>
    </div>
  );
}

function Confirmed({ name }: { name: string }) {
  const { withdrawApplication } = useApplicationActions();

  const handleWithdraw = () => {
    if (
      window.confirm(
        "Withdraw your attendance? This releases your spot at SwampHacks XII.",
      )
    ) {
      withdrawApplication.mutate();
    }
  };

  return (
    <div className="w-full sm:max-w-200 mx-auto font-figtree p-2 relative">
      <h1 className="text-2xl">You're confirmed, {name}! 🎉</h1>
      <div className="my-3 flex flex-col gap-3">
        <p>Your attendance at SwampHacks XII is confirmed.</p>
        <p>
          If you can no longer attend, withdraw your attendance so we can offer
          your spot to another applicant.
        </p>
        <Button
          onClick={handleWithdraw}
          isDisabled={withdrawApplication.isPending}
          size="md"
          className="w-[200px] max-w-full min-h-10"
          variant="danger"
        >
          {withdrawApplication.isPending
            ? "Withdrawing..."
            : "Withdraw Attendance"}
        </Button>
      </div>
    </div>
  );
}

interface RejectedProps {
  name: string;
}

function Rejected({ name }: RejectedProps) {
  const { joinWaitlist } = useApplicationActions();
  const waitlistClosed = Date.now() >= Date.parse("2026-10-16T00:00:00-04:00");

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

        <Button
          onClick={() => joinWaitlist.mutate()}
          isDisabled={joinWaitlist.isPending || waitlistClosed}
          className="w-fit"
        >
          {waitlistClosed
            ? "Waitlist Closed"
            : joinWaitlist.isPending
              ? "Joining..."
              : "Join Waitlist"}
        </Button>

        <p>The deadline to join is October 15, 2026 at 11:59 PM ET.</p>

        <ol className="flex flex-col gap-2">
          <li>
            1. <strong>Join the Waitlist</strong>: We may have openings
            available closer to the event. Use the Join Waitlist button to join
            the waitlist. Joining does not guarantee a spot; we will contact you
            if you are offered admission.
          </li>
          <li>
            2. <strong>Mentor</strong>: Share your knowledge and guide hackers
            through their projects.{" "}
            <a
              className="underline"
              href="https://swamphack.notion.site/3973b41de22f80b788ced816145e0a2d"
            >
              Sign up to be a mentor here
            </a>
            .
          </li>
          <li>
            3. <strong>Volunteer</strong>: Help us run the event smoothly.{" "}
            <a
              className="underline"
              href="https://swamphack.notion.site/3ae3b41de22f806ba9d5cb3f8ec65bed"
            >
              Sign up to volunteer here
            </a>
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
  const { leaveWaitlist } = useApplicationActions();

  const handleLeave = () => {
    if (
      window.confirm(
        "Leave the waitlist? You will lose your current place. If you rejoin before the deadline, you will join at the end of the queue.",
      )
    ) {
      leaveWaitlist.mutate();
    }
  };

  return (
    <div className="w-full sm:max-w-200 mx-auto font-figtree p-2 relative">
      <h1 className="text-2xl">Hi, {name}!</h1>
      <div className="my-3 flex flex-col gap-3">
        <p>
          You are on the <strong>SwampHacks XII waitlist</strong>. Joining does
          not guarantee admission.
        </p>
        <p>
          Your place is based on when you joined. If you receive an invitation,
          you will have 48 hours to confirm your attendance. Keep an eye on your
          inbox and spam folder.
        </p>
        <Button
          onClick={handleLeave}
          isDisabled={leaveWaitlist.isPending}
          variant="danger"
          size="md"
          className="w-[200px] max-w-full min-h-10"
        >
          {leaveWaitlist.isPending ? "Leaving..." : "Leave Waitlist"}
        </Button>
        <p>
          If you have questions, reach out on our{" "}
          <a className="underline" href="https://discord.com/invite/NfRPv9JtAG">
            Discord server
          </a>{" "}
          or email{" "}
          <a className="underline" href="mailto:contact@swamphacks.com">
            contact@swamphacks.com
          </a>
          .
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

function Fallback() {
  return (
    <div className="h-full flex justify-center items-center gap-2 text-red-400">
      <TablerAlertCircle />
      <p>Something went wrong while loading application form :(</p>
    </div>
  );
}
