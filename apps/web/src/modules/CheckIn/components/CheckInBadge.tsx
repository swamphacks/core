import { Badge } from "@/components/ui/Badge";

export function CheckInBadge({ isCheckedIn }: { isCheckedIn: boolean }) {
  if (isCheckedIn) {
    return (
      <Badge
        className="
        bg-blue-500/20 text-blue-700 border border-blue-400/30 
        dark:bg-blue-400/20 dark:text-blue-200 dark:border-blue-300/30
      "
      >
        Checked In
      </Badge>
    );
  }

  return (
    <Badge
      className="
      bg-red-500/20 text-red-700 border border-red-400/30 
      dark:bg-red-400/20 dark:text-red-200 dark:border-red-300/30
    "
    >
      Not Checked In
    </Badge>
  );
}
