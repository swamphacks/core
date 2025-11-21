import type { ComponentType } from "react";
import TablerBan from "~icons/tabler/ban";
import TablerUserCheck from "~icons/tabler/user-check";
import TablerConfetti from "~icons/tabler/confetti";
import TablerClockPause from "~icons/tabler/clock-pause";
import TablerHourglassFilled from "~icons/tabler/hourglass-filled";
import TablerPointFilled from "~icons/tabler/point-filled";
import TablerId from "~icons/tabler/id";
import TablerSettings2 from "~icons/tabler/settings-2";
import TablerCalendarCheck from "~icons/tabler/calendar-check";

type ApplicationStatus = {
  [k: string]: {
    className: string;
    text: string;
    icon?: ComponentType<React.SVGProps<SVGSVGElement>>;

    button: {
      className: string;
      text: string;
    };
  };
};

const defineStatus = <const T extends ApplicationStatus>(status: T) => {
  return status;
};

const applicationStatus = defineStatus({
  rejected: {
    className: "bg-badge-bg-rejected text-badge-text-rejected",
    text: "Rejected",
    icon: TablerBan,
    button: {
      className:
        "bg-event-button-bg-rejected text-event-button-text-rejected hover:bg-event-button-bg-rejected-hover",
      text: "Learn more",
    },
  },
  attending: {
    className: "bg-badge-bg-attending text-badge-text-attending",
    text: "Attending",
    icon: TablerUserCheck,
    button: {
      className:
        "bg-event-button-bg-attending text-event-button-text-attending hover:bg-event-button-bg-attending-hover",
      text: "Dashboard",
    },
  },
  accepted: {
    className: "bg-badge-bg-accepted text-badge-text-accepted",
    text: "Accepted",
    icon: TablerConfetti,
    button: {
      className:
        "bg-event-button-bg-accepted-attending text-event-button-text-accepted-attending hover:bg-event-button-bg-accepted-attending-hover",
      text: "I'm Attending",
    },
  },
  waitlisted: {
    className: "bg-badge-bg-waitlisted text-badge-text-waitlisted",
    text: "Waitlisted",
    icon: TablerClockPause,
    button: {
      className:
        "bg-event-button-bg-waitlisted text-event-button-text-waitlisted hover:bg-event-button-bg-waitlisted-hover",
      text: "What's Next?",
    },
  },
  underReview: {
    className: "bg-badge-bg-under-review text-badge-text-under-review",
    text: "Under Review",
    icon: TablerHourglassFilled,
    button: {
      className:
        "bg-event-button-bg-under-review text-event-button-text-under-review hover:bg-event-button-bg-under-review-hover",
      text: "Dashboard",
    },
  },
  notApplied: {
    className: "bg-badge-bg-not-applied text-badge-text-not-applied",
    text: "Not Applied",
    icon: TablerPointFilled,
    button: {
      className:
        "bg-event-button-bg-not-applied text-event-button-text-not-applied hover:bg-event-button-bg-not-applied-hover",
      text: "Apply Now",
    },
  },
  staff: {
    className: "bg-badge-bg-staff text-badge-text-staff",
    text: "Staff",
    icon: TablerId,
    button: {
      className:
        "bg-event-button-bg-staff text-event-button-text-staff hover:bg-event-button-bg-staff-hover",
      text: "Dashboard",
    },
  },
  admin: {
    className: "bg-badge-bg-admin text-badge-text-admin",
    text: "Admin",
    icon: TablerSettings2,
    button: {
      className:
        "bg-event-button-bg-admin text-event-button-text-admin hover:bg-event-button-bg-admin-hover",
      text: "Dashboard",
    },
  },
  notGoing: {
    className: "bg-badge-bg-not-going text-badge-text-not-going",
    text: "Not Going",
    icon: TablerBan,
    button: {
      className:
        "bg-event-button-bg-not-going text-event-button-text-not-going hover:bg-event-button-bg-not-going-hover",
      text: "Help Us Improve",
    },
  },
  completed: {
    className: "bg-badge-bg-completed text-badge-text-completed",
    text: "Completed",
    icon: TablerCalendarCheck,
    button: {
      className:
        "bg-event-button-bg-completed text-event-button-text-completed hover:bg-event-button-bg-completed-hover",
      text: "Event Summary",
    },
  },
  upcoming: {
    className: "bg-badge-bg-not-applied text-badge-text-not-applied",
    text: "Not Applied",
    icon: TablerPointFilled,
    button: {
      className:
        "bg-event-button-bg-completed text-event-button-text-completed hover:bg-event-button-bg-completed-hover",
      text: "Coming Soon",
    },
  },
});

export default applicationStatus;
