import {
  UserCheck,
  Rejected,
  Confetti,
  ClockPause,
  Hourglass,
  Point,
  ID,
  Settings,
  CalendarCheck,
} from "@/components/icons";
import type { IconWrapperProps } from "@/components/icons/IconWrapper";
import type { ComponentType } from "react";

type ApplicationStatus = {
  [k: string]: {
    className: string;
    text: string;
    icon?: ComponentType<IconWrapperProps>;
  };
};

const defineStatus = <const T extends ApplicationStatus>(status: T) => {
  return status;
};

const applicationStatus = defineStatus({
  rejected: {
    className: "bg-badge-bg-rejected text-badge-text-rejected",
    text: "Rejected",
    icon: Rejected,
  },
  attending: {
    className: "bg-badge-bg-attending text-badge-text-attending",
    text: "Attending",
    icon: UserCheck,
  },
  accepted: {
    className: "bg-badge-bg-accepted text-badge-text-accepted",
    text: "Accepted",
    icon: Confetti,
  },
  waitlisted: {
    className: "bg-badge-bg-waitlisted text-badge-text-waitlisted",
    text: "Waitlisted",
    icon: ClockPause,
  },
  underReview: {
    className: "bg-badge-bg-underReview text-badge-text-underReview",
    text: "Under Review",
    icon: Hourglass,
  },
  notApplied: {
    className: "bg-badge-bg-notApplied text-badge-text-notApplied",
    text: "Not Applied",
    icon: Point,
  },
  staff: {
    className: "bg-badge-bg-staff text-badge-text-staff",
    text: "Staff",
    icon: ID,
  },
  admin: {
    className: "bg-badge-bg-admin text-badge-text-admin",
    text: "Admin",
    icon: Settings,
  },
  notGoing: {
    className: "bg-badge-bg-notGoing text-badge-text-notGoing",
    text: "Not Going",
    icon: Rejected,
  },
  completed: {
    className: "bg-badge-bg-completed text-badge-text-completed",
    text: "Completed",
    icon: CalendarCheck,
  },
});

export default applicationStatus;
