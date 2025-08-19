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

function mapEventsAPIResponseToEventCardProps(data: EventWithUserInfo): EventCardProps {
    let status: keyof typeof applicationStatus  = "notApplied" // Default

    if (!data.application_status) {
        return {
            eventId: data.id,
            status,
            title: data.name,
            description: data.description ?? "No description",
            date: 
        }
    }

    switch (data.application_status?.application_status) {
        case "under_review" && data.application_status?.valid == true:

    }
}

