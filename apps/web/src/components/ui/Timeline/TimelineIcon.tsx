import TablerFileDescription from "~icons/tabler/file-description";
import TablerCalendarClock from "~icons/tabler/calendar-clock";
import TablerTerminal2 from "~icons/tabler/terminal-2";

export type IconName = "application-opens" | "calendar" | "circle-check-big";

interface TimelineIconProps {
  icon: IconName;
  completed: boolean;
}

const ICON_MAP: Record<IconName, React.ReactNode> = {
  "application-opens": <TablerFileDescription />,
  calendar: <TablerCalendarClock />,
  "circle-check-big": <TablerTerminal2 />,
};

export function TimelineIcon({ icon, completed }: TimelineIconProps) {
  return (
    <div
      className={
        "relative z-10 flex items-center justify-center w-8 h-8 rounded-lg border-2 " +
        (completed
          ? "bg-green-200 border-green-500 text-green-600"
          : "bg-white border-green-200 text-green-400")
      }
    >
      {ICON_MAP[icon]}
    </div>
  );
}
