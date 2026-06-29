import { Button } from "@/components/ui/Button";
import { TimelineIcon, type IconName } from "./TimelineIcon";
import TablerCalendar from "~icons/tabler/calendar";
import { Link } from "@tanstack/react-router";

export interface TimelineEvent {
  id?: string;
  date: string;
  title: string;
  description: string;
  icon: IconName;
  completed: boolean;
}

interface TimelineItemProps {
  event: TimelineEvent;
  isLast: boolean;
}

export function TimelineItem({ event, isLast }: TimelineItemProps) {
  const { id, date, title, description, icon, completed } = event;

  return (
    <div className={"relative flex items-start " + (!isLast ? "pb-8" : "pb-0")}>
      {!isLast && (
        <div
          className={
            "absolute left-4 top-8 w-0.5 h-full " +
            (completed ? "bg-green-400" : "bg-green-200")
          }
        ></div>
      )}
      <TimelineIcon icon={icon} completed={completed} />
      <div className="flex-1 ml-4">
        <div className="flex items-center gap-1 text-sm text-text-secondary mb-1">
          <TablerCalendar width={14} height={14} />
          {date}
        </div>
        <h3 className="font-semibold tex-text-main mb-1">{title}</h3>
        <p className="text-sm text-text-secondary">{description}</p>
        {id === "application-open" && (
          <Link to="/application">
            <Button
              className="mt-2 bg-green-600 hover:bg-green-700 pressed:bg-green-700"
              size="sm"
            >
              Application
            </Button>
          </Link>
        )}
      </div>
    </div>
  );
}
