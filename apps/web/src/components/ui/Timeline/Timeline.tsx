import { TimelineItem, type TimelineEvent } from "./TimelineItem";

interface TimelineProps {
  timeline: TimelineEvent[];
}

export function Timeline({ timeline }: TimelineProps) {
  return (
    <div className="max-w-3xl">
      {timeline.map((event, index) => (
        <TimelineItem
          key={`${event.title}-${index}`}
          event={event}
          isLast={index === timeline.length - 1}
        />
      ))}
    </div>
  );
}
