import { Button } from "@/components/ui/Button";
import TablerCheck from "~icons/tabler/check";

export default function TeamInvitationCard() {
  return (
    <div className="border border-input-border rounded-md px-6 py-4 w-full md:w-fit md:min-w-120 bg-surface flex flex-row justify-between items-center">
      <p>
        <strong>Alex</strong> has invited you to join{" "}
        <strong>The Pirates</strong>.
      </p>
      <div className="flex flex-row">
        <Button variant="icon">X</Button>
        <Button variant="icon" className="aspect-square p-2.5">
          <TablerCheck className="h-5 w-5" />
        </Button>
      </div>
    </div>
  );
}
