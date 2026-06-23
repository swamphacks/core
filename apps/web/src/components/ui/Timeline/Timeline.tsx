import { TimelineItem, type TimelineEvent } from "./TimelineItem";

interface TimelineProps {
  events?: TimelineEvent[];
}

const EVENTS: TimelineEvent[] = [
  {
    id: "application-open",
    date: "July 13",
    title: "Applications Open",
    description:
      "Apply for SwampHacks! Applications deadline is August 10, 2026 at 11:59 PM ET.",
    icon: "application-opens",
    completed: true,
  },
  {
    date: "August 10",
    title: "Application Deadline",
    description:
      "You must submit your application by August 10, 2026 at 11:59 PM ET to be considered for SwampHacks XII.",
    icon: "calendar",
    completed: true,
  },
  {
    date: "TBD",
    title: "Decisions Released",
    description:
      "Decision letters will be emailed to you in late September. You will either be accepted, waitlisted, or rejected.",
    icon: "calendar",
    completed: false,
  },
  {
    date: "TBD",
    title: "Confirmation Due",
    description:
      "If accepted, you must confirm your attendance for SwampHacks XII through the hacker portal.",
    icon: "calendar",
    completed: false,
  },
  {
    date: "Oct 16",
    title: "Event Day",
    description:
      "SwampHacks XII begins! More information will be sent to your email.",
    icon: "circle-check-big",
    completed: false,
  },
];

export function Timeline({ events = EVENTS }: TimelineProps) {
  return (
    <div className="max-w-3xl">
      {events.map((event, index) => (
        <TimelineItem
          key={`${event.title}-${index}`}
          event={event}
          isLast={index === events.length - 1}
        />
      ))}
    </div>
  );
}
