// TODO: Remove the eslint-disable line and fix the fast refresh issue
import { tv } from "tailwind-variants";
import applicationStatus from "../applicationStatus";
import { Button, button, type ButtonProps } from "@/components/ui/Button";
import { cn } from "@/utils/cn";
import { useRouter } from "@tanstack/react-router";
import { toast } from "react-toastify";
import { DialogTrigger } from "react-aria-components";
import { EventWithdrawalModal } from "./EventWithdrawalModal";
import { EventWaitlistModal } from "./EventWaitlistModal";

type ApplicationStatusTypes = keyof typeof applicationStatus;

const applicationStatusVariants = Object.fromEntries(
  Object.entries(applicationStatus).map(([key, value]) => [
    key,
    value.button.className,
  ]),
) as {
  [K in ApplicationStatusTypes]: (typeof applicationStatus)[K]["button"]["className"];
};

export const eventButton = tv({
  extend: button,
  variants: {
    variant: {},
    status: applicationStatusVariants,
  },
});

interface EventButtonProps extends ButtonProps {
  status: ApplicationStatusTypes;
  eventId: string;
  text?: string;
}

const EventButton = ({
  status: statusProp,
  eventId,
  className,
  text,
}: EventButtonProps) => {
  const eventButtonClassName = eventButton({
    status: statusProp,
    variant: "skeleton",
    className,
  });
  const router = useRouter();

  const onClick = () => {
    switch (statusProp) {
      case "accepted":
        // Make API call and then navigate to dashboard with accepted query param
        console.log("Accepted button clicked for event");
        router.navigate({
          to: `/events/${eventId}/dashboard`,
        });
        break;
      case "attending":
      case "staff":
      case "admin":
      case "underReview":
        // Navigate to the dashboard
        router.navigate({
          to: `/events/${eventId}/dashboard`,
        });
        break;
      case "waitlisted":
        // Navigate to the waitlist page
        router.navigate({
          to: `/events/${eventId}/waitlist/info`,
        });
        break;
      case "rejected":
        break;
      case "notApplied":
        // Navigate to the application page
        router.navigate({
          to: `/events/${eventId}/application`,
        });
        break;
      case "notGoing":
        break;
      case "completed":
        // Navigate to the event summary
        router.navigate({
          to: `/events/${eventId}/summary`,
        });
        break;
      case "upcoming":
        toast.info("Applications for this event are not currently open.", {
          position: "bottom-right",
        });
        break;
    }
  };

  if (statusProp === "notGoing") {
    return (
      <DialogTrigger>
        <Button
          variant="skeleton"
          className={cn(eventButtonClassName, "font-semibold")}
        >
          {text || applicationStatus[statusProp].button.text}
        </Button>
        <EventWithdrawalModal eventId={eventId}></EventWithdrawalModal>
      </DialogTrigger>
    );
  } else if (statusProp === "rejected") {
    return (
      <DialogTrigger>
        <Button
          variant="skeleton"
          className={cn(eventButtonClassName, "font-semibold")}
        >
          {text || applicationStatus[statusProp].button.text}
        </Button>
        <EventWaitlistModal eventId={eventId}></EventWaitlistModal>
      </DialogTrigger>
    );
  } else {
    return (
      <Button
        variant="skeleton"
        className={cn(eventButtonClassName, "font-semibold")}
        onClick={onClick}
      >
        {text || applicationStatus[statusProp].button.text}
      </Button>
    );
  }
};

export { EventButton };
