import TablerInfoCircle from "~icons/tabler/info-circle";
import { Heading, Text } from "react-aria-components";
import TablerCalendarEventFilled from "~icons/tabler/calendar-event-filled";
import TablerFileExport from "~icons/tabler/file-export";
import TablerUsers from "~icons/tabler/users";
import TablerMapPin from "~icons/tabler/map-pin";
import { useEffect, useState } from "react";
import { differenceInDays, type Duration, intervalToDuration } from "date-fns";
import type { EventOverview } from "@/features/EventOverview/hooks/useEventOverview";

interface EventDetailsProps {
  data: EventOverview;
}

const prettyPrintDate = (input: string | Date): string => {
  const date = new Date(input);

  return new Intl.DateTimeFormat("en-US", {
    month: "short",
    day: "numeric",
    hour: "numeric",
    minute: "2-digit",
    hour12: true,
  }).format(date);
};

export default function EventDetails({ data }: EventDetailsProps) {
  return (
    <div className="border border-input-border rounded-md px-4 py-3 w-full">
      <div className="flex justify-between items-start relative">
        <p className="flex items-center gap-1 bg-badge-bg-admin rounded-md px-2 py-1 w-fit">
          <TablerInfoCircle /> Event Details
        </p>

        <div className="flex items-center gap-2 text-text-secondary absolute right-0">
          <div className="bg-surface rounded-md px-2 py-1">
            <Countdown end_time={new Date(data.event_details.end_time)} />
          </div>
        </div>
      </div>

      <div className="mt-3 space-y-3">
        <div className="flex flex-col gap-2">
          <Heading className="text-xl font-semibold">
            {data.event_details.name}
          </Heading>
          <Text className="text-text-secondary">
            {data.event_details.description || "No description provided"}
          </Text>

          <span className="flex flex-row items-center gap-1.5 text-text-secondary">
            <TablerMapPin className="h-4 w-4" />
            {data.event_details.location
              ? data.event_details.location
              : "Unknown"}
          </span>
        </div>

        <div className="flex flex-col gap-1 text-sm text-text-secondary">
          <Text>Details:</Text>

          <span className="flex flex-row items-center gap-1.5">
            <TablerCalendarEventFilled className="h-4 w-4" /> Event Dates:{" "}
            {prettyPrintDate(data.event_details.start_time)} –{" "}
            {prettyPrintDate(data.event_details.end_time)}
          </span>

          <span className="flex flex-row items-center gap-1.5">
            <TablerFileExport className="h-4 w-4" /> Apply Dates:{" "}
            {prettyPrintDate(data.event_details.application_open)} –{" "}
            {prettyPrintDate(data.event_details.application_close)}
          </span>

          {data.event_details.max_attendees && (
            <span className="flex flex-row items-center gap-1.5">
              <TablerUsers className="h-4 w-4" /> Max Attendees:{" "}
              {data.event_details.max_attendees}
            </span>
          )}
        </div>
      </div>
    </div>
  );
}

interface CountdownProps {
  end_time: Date;
}

function Countdown({ end_time }: CountdownProps) {
  const [duration, setDuration] = useState<Duration | null>(null);

  useEffect(() => {
    const duration = intervalToDuration({
      start: new Date(),
      end: end_time,
    });

    setDuration({
      ...duration,
      days: differenceInDays(end_time, new Date()),
      hours: duration.hours ?? 0,
      minutes: duration.minutes ?? 0,
      seconds: duration.seconds ?? 0,
    });

    const id = setInterval(() => {
      const duration = intervalToDuration({
        start: new Date(),
        end: end_time,
      });

      setDuration({
        ...duration,
        days: differenceInDays(end_time, new Date()),
        hours: duration.hours ?? 0,
        minutes: duration.minutes ?? 0,
        seconds: duration.seconds ?? 0,
      });
    }, 1000);
    return () => clearInterval(id);
  }, []);

  if (!duration) return null;

  const renderSlot = (value: number, label: string) => {
    return (
      <div className="flex flex-col">
        <span className="countdown font-mono text-md">
          <span
            style={{ "--value": value } as React.CSSProperties}
            aria-live="polite"
            aria-label={String(value)}
          >
            {value}
          </span>
        </span>
        <span className="text-xs">{label}</span>
      </div>
    );
  };

  return (
    <div className="grid grid-flow-col gap-5 text-center auto-cols-max">
      {renderSlot(duration.days!, "days")}
      {renderSlot(duration.hours!, "hours")}
      {renderSlot(duration.minutes!, "min")}
      {renderSlot(duration.seconds!, "sec")}
    </div>
  );
}
