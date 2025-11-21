import type applicationStatus from "../applicationStatus";
import type { EventCardProps } from "../components/EventCard";
import { format } from "date-fns";
import type { EventWithUserInfo } from "../hooks/useEventsWithUserInfo";

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

const statusMap: Record<string, keyof typeof applicationStatus> = {
  accepted: "accepted",
  rejected: "rejected",
  waitlisted: "waitlisted",
  submitted: "underReview",
  under_review: "underReview",
  started: "notApplied",
  withdrawn: "notGoing",
  staff: "staff",
  admin: "admin",
  attendee: "attending",
};

export function mapEventsAPIResponseToEventCardProps(
  data: EventWithUserInfo,
): EventCardProps {
  let status: keyof typeof applicationStatus = "notApplied"; // Default

  if (!data.application_status) {
    return buildEventCardProps(data, status);
  }

  // Handle roles cases
  if (data.event_role?.event_role_type === "applicant") {
    // Handle applicant cases
    status =
      statusMap[data.application_status.application_status] ?? "notApplied";
  } else if (data.event_role?.event_role_type) {
    // Handle cases where user is staff/admin/attendee
    status = statusMap[data.event_role.event_role_type] ?? "notApplied";
  } else {
    // Handle cases where role type is undefined (applications may not be open)
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

  return buildEventCardProps(data, status);
}

/**
 *
 * @param event - Event data from API
 * @param status - Mapped application status
 * @returns An object conforming to EventCardProps
 */
function buildEventCardProps(
  event: EventWithUserInfo,
  status: keyof typeof applicationStatus,
): EventCardProps {
  return {
    eventId: event.id,
    status,
    title: event.name,
    description: event.description ?? "No description",
    date: formatDateRange(new Date(event.start_time), new Date(event.end_time)),
    location: event.location ?? "Unknown",
    banner: event.banner,

    application_close: new Date(event.application_close),
    application_open: new Date(event.application_open),
    end_time: new Date(event.end_time),
    start_time: new Date(event.start_time),
    location_url: event.location_url ?? undefined,
    decision_release: event.decision_release
      ? new Date(event.decision_release)
      : undefined,
    rsvp_deadline: event.rsvp_deadline
      ? new Date(event.rsvp_deadline)
      : undefined,
    website_url: event.website_url ?? undefined,
  };
}
