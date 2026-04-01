import { Button } from "@/components/ui/Button";
import TablerCheck from "~icons/tabler/check";
import TablerCircleX from "~icons/tabler/circle-x";

interface Props {
  requesterName?: string;
  onAccept: () => void;
  onReject: () => void;
}

export default function TeamJoinRequestCard({
  requesterName,
  onAccept,
  onReject,
}: Props) {
  return (
    <div className="border border-input-border rounded-md px-6 py-4 w-full md:w-fit md:min-w-120 bg-surface flex flex-row justify-between items-center">
      <p>
        <strong>{requesterName ?? "Anonymous"}</strong> has requested to join
        your team.
      </p>
      <div className="flex flex-row">
        <Button
          onClick={onReject}
          variant="icon"
          className="aspect-square p-2.5"
        >
          <TablerCircleX className="h-5 w-5" />
        </Button>
        <Button
          onClick={onAccept}
          variant="icon"
          className="aspect-square p-2.5"
        >
          <TablerCheck className="h-5 w-5" />
        </Button>
      </div>
    </div>
  );
}
