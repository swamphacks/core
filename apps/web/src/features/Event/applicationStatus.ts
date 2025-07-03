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
        "bg-eventBtn-bg-rejected text-eventBtn-text-rejected hover:bg-eventBtn-bg-rejected-hover",
      text: "Learn more",
    },
  },
  attending: {
    className: "bg-badge-bg-attending text-badge-text-attending",
    text: "Attending",
    icon: TablerUserCheck,
    button: {
      className:
        "bg-eventBtn-bg-attending text-eventBtn-text-attending hover:bg-eventBtn-bg-attending-hover",
      text: "Dashboard",
    },
  },
  accepted: {
    className: "bg-badge-bg-accepted text-badge-text-accepted",
    text: "Accepted",
    icon: TablerConfetti,
    button: {
      className:
        "bg-eventBtn-bg-accepted-attending text-eventBtn-text-accepted-attending hover:bg-eventBtn-bg-accepted-attending-hover",
      text: "I'm Attending",
    },
  },
  waitlisted: {
    className: "bg-badge-bg-waitlisted text-badge-text-waitlisted",
    text: "Waitlisted",
    icon: TablerClockPause,
    button: {
      className:
        "bg-eventBtn-bg-waitlisted text-eventBtn-text-waitlisted hover:bg-eventBtn-bg-waitlisted-hover",
      text: "What's Next?",
    },
  },
  underReview: {
    className: "bg-badge-bg-underReview text-badge-text-underReview",
    text: "Under Review",
    icon: TablerHourglassFilled,
    button: {
      className:
        "bg-eventBtn-bg-underReview text-eventBtn-text-underReview hover:bg-eventBtn-bg-underReview-hover",
      text: "Dashboard",
    },
  },
  notApplied: {
    className: "bg-badge-bg-notApplied text-badge-text-notApplied",
    text: "Not Applied",
    icon: TablerPointFilled,
    button: {
      className:
        "bg-eventBtn-bg-notApplied text-eventBtn-text-notApplied hover:bg-eventBtn-bg-notApplied-hover",
      text: "Apply Now",
    },
  },
  staff: {
    className: "bg-badge-bg-staff text-badge-text-staff",
    text: "Staff",
    icon: TablerId,
    button: {
      className:
        "bg-eventBtn-bg-staff text-eventBtn-text-staff hover:bg-eventBtn-bg-staff-hover",
      text: "Dashboard",
    },
  },
  admin: {
    className: "bg-badge-bg-admin text-badge-text-admin",
    text: "Admin",
    icon: TablerSettings2,
    button: {
      className:
        "bg-eventBtn-bg-admin text-eventBtn-text-admin hover:bg-eventBtn-bg-admin-hover",
      text: "Dashboard",
    },
  },
  notGoing: {
    className: "bg-badge-bg-notGoing text-badge-text-notGoing",
    text: "Not Going",
    icon: TablerBan,
    button: {
      className:
        "bg-eventBtn-bg-notGoing text-eventBtn-text-notGoing hover:bg-eventBtn-bg-notGoing-hover",
      text: "Help Us Improve",
    },
  },
  completed: {
    className: "bg-badge-bg-completed text-badge-text-completed",
    text: "Completed",
    icon: TablerCalendarCheck,
    button: {
      className:
        "bg-eventBtn-bg-completed text-eventBtn-text-completed hover:bg-eventBtn-bg-completed-hover",
      text: "Event Summary",
    },
  },
});

export default applicationStatus;
