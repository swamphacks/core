import { Button } from "@/components/ui/Button";
import { OverlayArrow, Tooltip, TooltipTrigger } from "react-aria-components";
import TablerCheck from "~icons/tabler/check";
import TablerCircleX from "~icons/tabler/circle-x";

export default function TeamInvitationCard() {
  return (
    <div className="border border-input-border rounded-md px-6 py-4 w-full md:w-fit md:min-w-120 bg-surface flex flex-row justify-between items-center">
      <p>
        <strong>Alex</strong> has invited you to join{" "}
        <strong>The Pirates</strong>.
      </p>
      <div className="flex flex-row">
        <TooltipTrigger delay={150}>
          <Button variant="icon" className="aspect-square p-2.5">
            <TablerCircleX className="h-5 w-5" />
          </Button>
          <Tooltip>
            <OverlayArrow>
              <svg width={8} height={8} viewBox="0 0 8 8">
                <path d="M0 0 L4 4 L8 0" />
              </svg>
            </OverlayArrow>
            Reject
          </Tooltip>
        </TooltipTrigger>
        <Button variant="icon" className="aspect-square p-2.5">
          <TablerCheck className="h-5 w-5" />
        </Button>
      </div>
    </div>
  );
}
