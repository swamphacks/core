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

  if (!data.application_status) {
    return {
      eventId: data.id,
      status,
      title: data.name,
      description: data.description ?? "No description",
      date: formatDateRange(new Date(data.start_time), new Date(data.end_time)),
      location: data.location ?? "Unknown",
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
    // If no specific role or application status, default to notApplied
    status = "notApplied";
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
  };
}
