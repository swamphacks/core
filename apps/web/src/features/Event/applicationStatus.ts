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
  },
  attending: {
    className: "bg-badge-bg-attending text-badge-text-attending",
    text: "Attending",
    icon: TablerUserCheck,
  },
  accepted: {
    className: "bg-badge-bg-accepted text-badge-text-accepted",
    text: "Accepted",
    icon: TablerConfetti,
  },
  waitlisted: {
    className: "bg-badge-bg-waitlisted text-badge-text-waitlisted",
    text: "Waitlisted",
    icon: TablerClockPause,
  },
  underReview: {
    className: "bg-badge-bg-underReview text-badge-text-underReview",
    text: "Under Review",
    icon: TablerHourglassFilled,
  },
  notApplied: {
    className: "bg-badge-bg-notApplied text-badge-text-notApplied",
    text: "Not Applied",
    icon: TablerPointFilled,
  },
  staff: {
    className: "bg-badge-bg-staff text-badge-text-staff",
    text: "Staff",
    icon: TablerId,
  },
  admin: {
    className: "bg-badge-bg-admin text-badge-text-admin",
    text: "Admin",
    icon: TablerSettings2,
  },
  notGoing: {
    className: "bg-badge-bg-notGoing text-badge-text-notGoing",
    text: "Not Going",
    icon: TablerBan,
  },
  completed: {
    className: "bg-badge-bg-completed text-badge-text-completed",
    text: "Completed",
    icon: TablerCalendarCheck,
  },
});

export default applicationStatus;
