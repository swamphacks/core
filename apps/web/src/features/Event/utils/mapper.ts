// Mapping event api reponse to usable structure for event card

import type { EventWithUserInfo } from "@/lib/openapi/types";
import type applicationStatus from "../applicationStatus";
import type { EventCardProps } from "../components/EventCard";
import { format } from "date-fns";

function formatDateRange(start: Date, end: Date): string {
  const startDay = format(start, "d");
  const endDay = format(end, "do"); // adds "st", "nd", "rd", "th"
  const month = format(start, "MMM");

  // If start and end are in the same month
  if (format(start, "MMM") === format(end, "MMM")) {
    return `${month} ${startDay}-${endDay}`;
  }

  // Different months
  const startMonth = format(start, "MMM");
  const endMonth = format(end, "MMM");
  return `${startMonth} ${startDay} - ${endMonth} ${endDay}`;
}

export function mapEventsAPIResponseToEventCardProps(
  data: EventWithUserInfo,
): EventCardProps {
  let status: keyof typeof applicationStatus = "notApplied"; // Default

  //TODO: Remove this guard once OpenAPI is updated
  if (
    !data.id ||
    !data.name ||
    !data.start_time ||
    !data.end_time ||
    !data.application_open ||
    !data.application_close
  ) {
    throw new Error("Missing required event fields");
  }

  if (!data.application_status) {
    return {
      eventId: data.id,
      status,
      title: data.name,
      description: data.description ?? "No description",
      date: formatDateRange(new Date(data.start_time), new Date(data.end_time)),
      location: data.location ?? "Unknown",
      banner: data.banner,
    };
  }

  // Handle roles cases
  if (data.event_role?.event_role_type === "staff") {
    status = "staff";
  } else if (data.event_role?.event_role_type === "admin") {
    status = "admin";
  } else if (data.event_role?.event_role_type === "attendee") {
    status = "attending";
  } else if (data.event_role?.event_role_type === "applicant") {
    // Handle application cases
    switch (data.application_status.application_status) {
      case "accepted":
        status = "accepted";
        break;
      case "rejected":
        status = "rejected";
        break;
      case "waitlisted":
        status = "waitlisted";
        break;
      case "submitted":
      case "under_review":
        status = "underReview";
        break;
      case "started":
        status = "notApplied";
        break;
      case "withdrawn":
        status = "notGoing";
        break;
      default:
        status = "notApplied"; // Default case
        break;
    }
  } else {
    // Check if applications are even open
    const now = new Date();
    const appStart = new Date(data.application_open);
    const appEnd = new Date(data.application_close);

    if (now >= appStart && now <= appEnd) {
      // Applications are open
      status = "notApplied";
    } else if (now < appStart) {
      // Applications not yet open
      status = "upcoming";
    } else if (now > appEnd) {
      // Applications closed
      status = "notGoing";
    }
  }

  // Check if the event is completed
  if (new Date(data.end_time) < new Date()) {
    status = "completed";
  }

  return {
    eventId: data.id,
    status,
    title: data.name,
    description: data.description ?? "No description",
    date: formatDateRange(new Date(data.start_time), new Date(data.end_time)),
    location: data.location ?? "Unknown",
    banner: data.banner,
  };
}
